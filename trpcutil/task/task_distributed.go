// Package task 定时任务
package task

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/twbbb/toolbox/trpcutil/entity"
	"github.com/twbbb/toolbox/trpcutil/monitor"
	"github.com/twbbb/toolbox/util/osutil"
	"trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/codec"
	"trpc.group/trpc-go/trpc-go/log"
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

// InitTaskTickerEnvWithCron 初始定时任务，包含env，支持 cron 表达式
func InitTaskTickerEnvWithCron(
	runFunc func(ctx context.Context) error, taskName string, cronExpr string,
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

	return InitTaskTickerWithCron(runFunc, taskName, cronExpr, baseLockKey, lock)
}

// InitTaskTickerWithCron 初始定时任务，支持 cron 表达式
func InitTaskTickerWithCron(
	runFunc func(ctx context.Context) error, taskName string, cronExpr string, baseLockKey string,
	lock entity.DistributLock,
) error {
	taskFunc := GetDefaultTaskDistributFunc(runFunc, taskName, 0, baseLockKey, lock)
	taskFunc.CronExpr = cronExpr
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
// 若 taskFunc.CronExpr 不为空则使用 cron 表达式调度，否则使用 CornDuration 秒级间隔
func DoDistributTask(taskFunc entity.TaskDistributFunc) error {
	initRequestId := fmt.Sprintf("%s_%d", taskFunc.TaskName, 0)
	if taskFunc.InitStart {
		if err := runDistributProcess(taskFunc, initRequestId); err != nil {
			return err
		}
	}

	if taskFunc.CronExpr != "" {
		return doDistributTaskByCron(taskFunc)
	}
	return doDistributTaskByTicker(taskFunc)
}

// doDistributTaskByCron 使用 cron 表达式调度
func doDistributTaskByCron(taskFunc entity.TaskDistributFunc) error {
	// 支持秒级（6位）cron 表达式
	c := cron.New(cron.WithSeconds())
	runTimes := 0
	var isRunning bool

	_, err := c.AddFunc(taskFunc.CronExpr, func() {
		runTimes++
		if runTimes >= math.MaxInt32 {
			runTimes = 1
		}
		requestId := fmt.Sprintf("%s_%d", taskFunc.TaskName, runTimes)
		if isRunning {
			log.Info("task is running：%s", requestId)
			return
		}
		isRunning = true
		runDistributProcess(taskFunc, requestId)
		isRunning = false
	})
	if err != nil {
		return fmt.Errorf("cron expr invalid [%s]: %w", taskFunc.CronExpr, err)
	}
	c.Start()
	return nil
}

// doDistributTaskByTicker 使用秒级 ticker 调度（原有逻辑）
func doDistributTaskByTicker(taskFunc entity.TaskDistributFunc) error {
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
	lockTime := calcLockTime(taskFunc)
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
	useTime := getUseTime(begin)
	monitor.ReportCache(requestId, taskFunc.ReportName, taskFunc.TaskName, true, useTime)
	log.InfoContextf(ctx, "runDistributProcess end, taskName: %s, cost: %.3fs", taskFunc.TaskName, useTime)
	return err
}

// calcLockTime 计算分布式锁的持有时间（秒）
// - 若使用 CronExpr，通过解析表达式计算相邻两次触发的间隔
// - 否则由 CornDuration 推算（减去2秒缓冲）
func calcLockTime(taskFunc entity.TaskDistributFunc) int {
	if taskFunc.CronExpr != "" {
		parser := cron.NewParser(
			cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
		)
		schedule, err := parser.Parse(taskFunc.CronExpr)
		if err == nil {
			now := time.Now()
			next1 := schedule.Next(now)
			next2 := schedule.Next(next1)
			intervalSec := int(next2.Sub(next1).Seconds())
			if intervalSec > 2 {
				return intervalSec - 2
			}
			return intervalSec
		}
	}
	lockTime := taskFunc.CornDuration
	if lockTime > 3 {
		lockTime = lockTime - 2
	}
	return lockTime
}
