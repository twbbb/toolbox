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
)

// InitCacheTicker 初始化本地缓存任务
func InitCacheTicker(runFunc func(ctx context.Context) error, taskName string, secondDuration int) error {
	taskFunc := GetDefaultTaskFunc(runFunc, taskName, secondDuration)
	return DoTask(taskFunc)
}

// GetDefaultTaskFunc 获取默认的taskFunc
func GetDefaultTaskFunc(runFunc func(ctx context.Context) error, taskName string, secondDuration int) entity.TaskFunc {
	taskFunc := entity.TaskFunc{
		TaskName:     taskName,
		CornDuration: secondDuration,
		TaskFunc:     runFunc,
		InitStart:    true,
		Logger:       log.GetDefaultLogger(),
		ReportName:   "cache",
	}
	return taskFunc
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

			monitor.ReportCache(requestId, taskFunc.ReportName, taskFunc.TaskName, false, getUseTime(begin))
			log.ErrorContextf(ctx, "runProcess panic %+v", r)
		}

	}()
	log.InfoContextf(ctx, "runProcess start")
	err = taskFunc.TaskFunc(ctx)
	if err != nil {
		log.ErrorContextf(ctx, "runProcess err%+v", err)
	}
	monitor.ReportCache(requestId, taskFunc.ReportName, taskFunc.TaskName, err == nil, getUseTime(begin))
	log.InfoContextf(ctx, "runProcess end")
	return err
}

func getUseTime(begin time.Time) float64 {
	// 计算经过的时间，并转换为秒
	elapsed := time.Since(begin).Seconds()
	// 格式化输出，保留两位小数
	elapsedRounded := math.Round(elapsed*1000) / 1000
	return elapsedRounded
}
