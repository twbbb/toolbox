package task

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/codec"
	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.woa.com/th-go/toolbox/trpcutil/entity"
	"git.woa.com/th-go/toolbox/trpcutil/report"
)

// InitTaskTicker 初始定时任务
func InitTaskTicker(
	runFunc func(ctx context.Context) error, taskName string, secondDuration int, lock entity.DistributLock,
) error {
	taskFunc := getDefaultTaskDistributFunc(runFunc, taskName, secondDuration, lock)
	return DoDistributTask(taskFunc)
}

func getDefaultTaskDistributFunc(
	runFunc func(ctx context.Context) error, taskName string, secondDuration int, lock entity.DistributLock,
) entity.TaskDistributFunc {
	taskFunc := entity.TaskDistributFunc{
		TaskName:      taskName,
		CornDuration:  secondDuration,
		TaskFunc:      runFunc,
		InitStart:     false,
		MaxRunTime:    30,
		BaseLockKey:   "task:",
		DistributLock: lock,
		Logger:        log.GetDefaultLogger(),
		ReportName:    "task",
	}
	return taskFunc
}

// 执行分布式定时任务
func DoDistributTask(taskFunc entity.TaskDistributFunc) error {
	initRequestId := fmt.Sprintf("%s_%d", taskFunc.TaskName, 0)
	if taskFunc.InitStart {
		if err := runDistributProcess(taskFunc, initRequestId); err != nil {
			return err
		}
	}
	// 创建一个计时器
	go func() {
		ticks := time.NewTicker(time.Duration(taskFunc.CornDuration) * time.Minute)
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // 确保上下文在函数退出时取消

	// 监听 SIGTERM 信号
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)
	go func() {
		select {
		case <-sigs:
			log.InfoContextf(ctx, "runDistributProcess SIGTERM")
			cancel() // 在接收到 SIGTERM 信号时取消上下文
			return
		case <-ctx.Done():
			return // 如果主上下文被取消，也停止此 goroutine
		}
	}()

	begin := time.Now()
	loggerInfo := taskFunc.Logger
	ctx, msg := codec.EnsureMessage(ctx)
	trpc.SetMetaData(ctx, "requestId", []byte(requestId))
	loggerInfo = loggerInfo.With(log.Field{Key: "requestId", Value: requestId})
	msg.WithLogger(loggerInfo)

	lockKey := taskFunc.BaseLockKey + taskFunc.TaskName + ":maxlock"

	// 任务执行的锁
	lock := taskFunc.DistributLock.Lock(ctx, lockKey, time.Minute*time.Duration(taskFunc.MaxRunTime))
	if !lock {
		return nil
	}

	defer taskFunc.DistributLock.UnLock(ctx, lockKey)

	// 定时任务，防止重复执行，比如多台机器，3分钟只能执行一次
	cornKey := taskFunc.BaseLockKey + taskFunc.TaskName + ":corn"
	lock = taskFunc.DistributLock.Lock(ctx, cornKey, time.Minute*time.Duration(taskFunc.CornDuration))
	if !lock {
		return nil
	}

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
			report.ReportCache(requestId, taskFunc.ReportName, taskFunc.TaskName, false, getUseTimeSec(begin))

			log.ErrorContextf(ctx, "runProcess panic %+v", r)
		}

	}()
	log.InfoContextf(ctx, "runDistributProcess start")

	err = taskFunc.TaskFunc(ctx)
	if err != nil {
		log.ErrorContextf(ctx, "runProcess err%+v", err)
	}
	report.ReportCache(requestId, taskFunc.ReportName, taskFunc.TaskName, true, getUseTimeSec(begin))
	log.InfoContextf(ctx, "runDistributProcess end")
	return err
}
