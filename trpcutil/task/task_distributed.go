// Package task 定时任务
package task

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/codec"
	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.woa.com/th-go/toolbox/trpcutil/entity"
	"git.woa.com/th-go/toolbox/trpcutil/monitor"
	"git.woa.com/th-go/toolbox/util/osutil"
)

// InitTaskTickerEnv 初始定时任务，包含env
func InitTaskTickerEnv(
	runFunc func(ctx context.Context) error, taskName string, secondDuration int,
	lock entity.DistributLock,
) error {
	baseLockKey := ""
	if trpc.GlobalConfig() != nil {
		if trpc.GlobalConfig().Server.Server != "" {
			baseLockKey = trpc.GlobalConfig().Server.Server + ":"
		} else {
			baseLockKey = "task:"
		}
		if trpc.GlobalConfig().Global.EnvName != "" {
			baseLockKey = baseLockKey + trpc.GlobalConfig().Global.EnvName + ":"
		}
	}

	if baseLockKey == "" {
		baseLockKey = "task:"
	}

	return InitTaskTicker(runFunc, taskName, secondDuration, baseLockKey, lock)
}

// InitTaskTicker 初始定时任务
func InitTaskTicker(
	runFunc func(ctx context.Context) error, taskName string, secondDuration int, baseLockKey string,
	lock entity.DistributLock,
) error {
	taskFunc := GetDefaultTaskDistributFunc(runFunc, taskName, secondDuration, baseLockKey, lock)
	return DoDistributTask(taskFunc)
}

// GetDefaultTaskDistributFunc 获取默认的分布式定时任务
func GetDefaultTaskDistributFunc(
	runFunc func(ctx context.Context) error, taskName string, secondDuration int, baseLockKey string,
	lock entity.DistributLock,
) entity.TaskDistributFunc {
	taskFunc := entity.TaskDistributFunc{
		TaskName:      taskName,
		CornDuration:  secondDuration,
		TaskFunc:      runFunc,
		InitStart:     false,
		MaxRunTime:    60 * 60, // 60分钟
		BaseLockKey:   baseLockKey,
		DistributLock: lock,
		Logger:        log.GetDefaultLogger(),
		ReportName:    "task",
	}
	return taskFunc
}

// DoDistributTask 执行分布式定时任务
func DoDistributTask(taskFunc entity.TaskDistributFunc) error {
	initRequestId := fmt.Sprintf("%s_%d", taskFunc.TaskName, 0)
	if taskFunc.InitStart {
		if err := runDistributProcess(taskFunc, initRequestId); err != nil {
			return err
		}
	}
	// 创建一个计时器
	go func() {
		ticks := time.NewTicker(time.Duration(taskFunc.CornDuration) * time.Second)
		tick := ticks.C
		runTimes := 0
		var isRunning bool // 用来标记任务是否正在运行

		for {
			<-tick
			runTimes++
			if runTimes >= math.MaxInt32 {
				runTimes = 1
			}
			requestId := fmt.Sprintf("%s_%d", taskFunc.TaskName, runTimes)
			if isRunning { // 检查是否有任务正在运行
				log.Info("task is running：%s", requestId)
				continue
			}

			isRunning = true
			runDistributProcess(taskFunc, requestId)
			isRunning = false
		}
	}()
	return nil
}

// 分布式执行
func runDistributProcess(
	taskFunc entity.TaskDistributFunc, requestId string,
) (err error) {

	// 创建可以取消的上下文
	ctx, cancel := context.WithCancel(trpc.BackgroundContext())
	defer cancel() // 确保上下文在函数退出时取消

	begin := time.Now()
	loggerInfo := taskFunc.Logger
	ctx, msg := codec.EnsureMessage(ctx)
	trpc.SetMetaData(ctx, "requestId", []byte(requestId))
	loggerInfo = loggerInfo.With(log.Field{Key: "requestId", Value: requestId})
	msg.WithLogger(loggerInfo)

	// 定时任务，防止重复执行，比如多台机器，3分钟只能执行一次
	cornKey := taskFunc.BaseLockKey + taskFunc.TaskName + ":corn"
	lockTime := taskFunc.CornDuration
	if lockTime > 3 {
		lockTime = lockTime - 2
	}
	lockState := taskFunc.DistributLock.Lock(ctx, cornKey, time.Second*time.Duration(lockTime))
	if !lockState {
		log.DebugContextf(ctx, "runDistributProcess corn fail:%s", cornKey)
		return nil
	}

	maxTimeLockKey := taskFunc.BaseLockKey + taskFunc.TaskName + ":maxlock"
	// 任务执行最大时间的锁
	lock := taskFunc.DistributLock.Lock(ctx, maxTimeLockKey, time.Second*time.Duration(taskFunc.MaxRunTime))
	if !lock {
		log.DebugContextf(ctx, "runDistributProcess lock fail:%s", maxTimeLockKey)
		return nil
	}
	// 重启时候，接收到 SIGTERM 信号时，关闭上下文信号，修改锁定时间为2分钟。
	osutil.ListenForSIGTERM(
		ctx, cancel, func() {
			log.InfoContextf(ctx, "runDistributProcess SIGTERM")
			taskFunc.DistributLock.UnLock(ctx, maxTimeLockKey)
			lock := taskFunc.DistributLock.Lock(ctx, maxTimeLockKey, time.Minute*2)
			if lock {
				log.InfoContextf(ctx, "MaxRunTime lock 2 min suc at SIGTERM")
			} else {
				log.InfoContextf(ctx, "MaxRunTime lock fail:%s", maxTimeLockKey)
			}
		},
	)

	defer taskFunc.DistributLock.UnLock(ctx, maxTimeLockKey)

	defer func() {
		// 放回EnsureMessage分配的的msg
		codec.PutBackMessage(msg)
		// panic时候程序不挂
		if r := recover(); r != nil {
			// check exactly what the panic was and create error.
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("unknow panic")
			}
			monitor.ReportCache(requestId, taskFunc.ReportName, taskFunc.TaskName, false, getUseTime(begin))

			log.ErrorContextf(
				ctx, "%s,%s,%s,%v,runProcess panic %+v", requestId, taskFunc.ReportName, taskFunc.TaskName,
				getUseTime(begin), r,
			)
		}

	}()
	log.InfoContextf(ctx, "runDistributProcess start")

	err = taskFunc.TaskFunc(ctx)
	if err != nil {
		log.ErrorContextf(ctx, "runProcess err%+v", err)
	}
	monitor.ReportCache(requestId, taskFunc.ReportName, taskFunc.TaskName, true, getUseTime(begin))
	log.InfoContextf(ctx, "runDistributProcess end")
	return err
}
