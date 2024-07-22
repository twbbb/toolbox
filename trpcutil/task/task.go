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
	"git.woa.com/th-go/toolbox/trpcutil/report"
)

// InitCacheTicker 初始化缓存和定时器
func InitCacheTicker(runFunc func(ctx context.Context) error, taskName string, secondDuration int) error {
	taskFunc := entity.TaskFunc{
		TaskName:     taskName,
		CornDuration: secondDuration,
		TaskFunc:     runFunc,
		InitStart:    true,
		Logger:       log.GetDefaultLogger(),
	}
	return DoTask(taskFunc)
}

// DoTask 执行定时任务
func DoTask(taskFunc entity.TaskFunc) error {
	initRequestId := fmt.Sprintf("%s_%d", taskFunc.TaskName, 0)
	if taskFunc.InitStart {
		if err := runProcess(taskFunc, initRequestId); err != nil {
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
			runProcess(taskFunc, requestId)
			isRunning = false
		}
	}()

	return nil
}

// 执行
func runProcess(
	taskFunc entity.TaskFunc, requestId string,
) (err error) {

	begin := time.Now()
	loggerInfo := taskFunc.Logger
	ctx, msg := codec.EnsureMessage(context.Background())
	trpc.SetMetaData(ctx, "requestId", []byte(requestId))
	loggerInfo = loggerInfo.With(log.Field{Key: "requestId", Value: requestId})
	msg.WithLogger(loggerInfo)

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

			report.ReportCache(requestId, "task", taskFunc.TaskName, true, getUseTimeSec(begin))
			log.ErrorContextf(ctx, "runProcess panic %+v", r)
		}

	}()

	err = taskFunc.TaskFunc(ctx)
	if err != nil {
		log.ErrorContextf(ctx, "runProcess err%+v", err)
	}
	report.ReportCache(requestId, "task", taskFunc.TaskName, err == nil, getUseTimeSec(begin))
	return err
}

func getUseTimeSec(begin time.Time) float64 {
	// 计算经过的时间，并转换为秒
	elapsed := time.Since(begin).Seconds()
	// 格式化输出，保留两位小数
	elapsedRounded := math.Round(elapsed*100) / 100
	return elapsedRounded
}
