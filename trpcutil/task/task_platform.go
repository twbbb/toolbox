// Package task 对应定时任务平台相关调用
package task

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.woa.com/th-go/toolbox/trpcutil/entity"
	"git.woa.com/th-go/toolbox/trpcutil/monitor"
	"net/http"
	"runtime/debug"
	"time"
)

// 使用示例：
const (
	webhookURL = "https://qyapi.weixin.qq.com/cgi-bin/webhook"
)

var distributeLock entity.DistributLock
var taskFuncs = map[string]*entity.PlatformTask{}

// InitDistributeLock 初始化分布式锁
func InitDistributeLock(lock entity.DistributLock) {
	distributeLock = lock
}

// AddTask 添加任务
func AddTask(taskNam string, taskFunc entity.PlatformTaskFunc) {
	if taskFuncs[taskNam] != nil {
		log.Errorf("DoPlatformTask|AddTask| 重复注册任务:%s", taskNam)
		return
	}
	taskFuncs[taskNam] = &entity.PlatformTask{
		TaskName: taskNam,
		Func:     taskFunc,
	}
}

// AddTaskWithReportCorn 添加任务并添加corn,corn("0 0 1 * * *")会上报到智妍
func AddTaskWithReportCorn(taskNam string, taskFunc entity.PlatformTaskFunc, cornTime string) {
	if taskFuncs[taskNam] != nil {
		log.Errorf("DoPlatformTask|AddTask| 重复注册任务:%s", taskNam)
		return
	}
	taskFuncs[taskNam] = &entity.PlatformTask{
		TaskName: taskNam,
		Func:     taskFunc,
		CornTime: cornTime,
	}
}

// DoPlatformTask 处理平台下发的任务
func DoPlatformTask(ctx context.Context, msg string) (err error) {
	var taskDto entity.PlatformTaskDto
	err = json.Unmarshal([]byte(msg), &taskDto)
	if err != nil {
		log.ErrorContextf(ctx, "DoPlatformTask|DoTask|json 解析出现错误:%+v", err)
		return err
	}
	log.InfoContextf(ctx, "DoPlatformTask|DoTask|taskDto:%s", jsonMarshal(taskDto))
	newCtx, cancel := getWithCanCelContext(ctx, taskDto)
	go func() {
		// 确保上下文在函数退出时取消
		defer cancel()
		err = doTask(newCtx, taskDto)
	}()
	return
}

// doTask 处理任务
func doTask(ctx context.Context, taskDto entity.PlatformTaskDto) (err error) {
	start := time.Now()
	// 处理异常情况
	defer func() {
		if p := recover(); p != nil {
			log.ErrorContextf(ctx, "DoPlatformTask|DoTask:%s|panic:%v,%s", taskDto.TaskName, p, debug.Stack())
			debug.PrintStack()
			monitor.ReportPlatformTask(taskDto.TaskName, false, getUseTime(start), taskFuncs[taskDto.TaskName].CornTime)
		}

	}()

	taskFunc := taskFuncs[taskDto.TaskName]
	if taskFunc == nil {
		log.ErrorContextf(ctx, "DoPlatformTask|DoTask|未注册的任务:%s", taskDto.TaskName)
		return
	}
	if taskDto.LockKey != "" && distributeLock != nil {
		lock := distributeLock.Lock(ctx, taskDto.LockKey, time.Duration(taskDto.LockTime)*time.Second)
		if !lock {
			log.ErrorContextf(ctx, "DoPlatformTask|DoTask| 获取锁失败:%s", taskDto.LockKey)
			robotAlarm(ctx, taskDto, fmt.Sprintf(",获取锁失败:%s", taskDto.LockKey))
			monitor.ReportPlatformTask(taskDto.TaskName, false, getUseTime(start), taskFuncs[taskDto.TaskName].CornTime)
			return
		}
		defer distributeLock.UnLock(ctx, taskDto.LockKey)
	}
	err = taskFunc.Func(ctx, taskDto.TransformMsg)
	if err != nil {
		log.ErrorContextf(ctx, "DoPlatformTask|DoTask| 任务执行出现错误:%s|err:%+v", taskDto.TaskName, err)
		return
	}
	monitor.ReportPlatformTask(taskDto.TaskName, true, getUseTime(start), taskFuncs[taskDto.TaskName].CornTime)
	log.InfoContextf(ctx, "DoPlatformTask|DoTask|taskName:%s|costTime:%s|info:%s", taskDto.TaskName, time.Since(start), taskDto.TransformMsg)
	if taskDto.AlarmTimeSec > 0 && time.Since(start).Seconds() > float64(taskDto.AlarmTimeSec) {
		log.ErrorContextf(ctx, "DoPlatformTask|DoTask|taskName:%s|cost:%v|info:%s|超时告警", taskDto.TaskName, time.Since(start), taskDto.TransformMsg)
		robotAlarm(ctx, taskDto, fmt.Sprintf(",耗时:%s,超过了设置的默认告警时间%ds", time.Since(start).String(), taskDto.AlarmTimeSec))
	}
	return err
}

func robotAlarm(ctx context.Context, taskDto entity.PlatformTaskDto, markMsg string) {
	if taskDto.RobotLink == "" {
		return
	}
	trpcMsg := trpc.Message(ctx)
	rpcName := trpcMsg.ServerRPCName()
	ip := trpc.GlobalConfig().Global.LocalIP
	msg := fmt.Sprintf("rpcName:%s,ip:%s,taskName:%s,%s", rpcName, ip, taskDto.TaskName, markMsg)
	values := entity.RobotSendReq{
		Msgtype: "markdown",
		Markdown: entity.Markdown{
			Content: msg,
		},
	}
	jsonValue, _ := json.Marshal(values)
	robotHookURL := fmt.Sprintf("%s/send?key=%s", webhookURL, taskDto.RobotLink)
	_, err := http.Post(
		robotHookURL,
		"application/json",
		bytes.NewBuffer(jsonValue),
	)
	if err != nil {
		log.ErrorContextf(ctx, "DoPlatformTask|robotSendMDAlarm|err:%v", err)
	}
}

func getWithCanCelContext(ctx context.Context, dto entity.PlatformTaskDto) (newCtx context.Context, cancel context.CancelFunc) {
	//ctx, _ = codec.WithCloneContextAndMessage(ctx)
	ctx = trpc.CloneContext(ctx)
	if dto.TimeoutCtx > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(dto.TimeoutCtx)*time.Minute)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}
	return ctx, cancel
}

// JsonMarshal , 变为json
func jsonMarshal(v interface{}) string {
	if v == nil {
		return ""
	}
	marshal, _ := json.Marshal(v)
	return string(marshal)
}
