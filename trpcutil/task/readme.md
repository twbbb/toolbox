- **task**: 本地定时任务实现,用于本地缓存，会统一上报对应任务服务名，定时任务名，耗时，错误，次数 等
- **task_distributed**：分布式注册，智研监控会有相关耗时(task框架已统一上报)
- **task_platform**：定时任务平台对接，会统一上报对应任务服务名，定时任务名，耗时，错误，次数 等


## task_platform
> 对于注册的方法，会自动上报智妍自定义监控。支持超时企业微信告警，加锁防止重复执行
使用案例

1. 初始化定时任务
```go
// Init 
func Init() error {
	// 如果定时任务不需要加锁，可不进行如下初始化
	task.InitDistributeLock(redis.Get())
	// 添加refreshAllLegoModules定时任务
	task.AddTask("refreshAllLegoModules", func(ctx context.Context, data string) error {
		return service.RefreshAllLegoModules(ctx)
	})
	// 添加refreshGrayWenzhenDoctor定时任务
	task.AddTask("refreshGrayWenzhenDoctor", func(ctx context.Context, data string) error {
		return service.RefreshGrayWenzhenDoctor(ctx)
	})

    task.AddTaskWithReportCorn("test2", func(ctx context.Context, data string) error {
        return service.test2(ctx)
    }, time.Duration(11)*time.Second)
	return nil
}
```
2. pb定时任务实现层使用如下：
```go
// DoTask 定时任务入口
func (t *ThMedicalCenter) DoTask(ctx context.Context, req *pb.ReceiveMsgRequest) (rsp *pb.ReceiveMsgResponse, err error) {
	task.DoPlatformTask(ctx, req.Msg)
	rsp = &pb.ReceiveMsgResponse{}
	return rsp, err
}
```
3. 定时任务平台配置案例
在定时任务的业务消息中，对于不同的任务配置不同的内容

执行refreshAllLegoModules
```json
{"taskName": "refreshAllLegoModules","lockKey": "SyncPublicPkgList_1000005891",
  "lockTime": 30, "alarmTimeSec": 1, "robotLink": "0d9a679d-10f6-44fe-a210-37f0690bdb6e" }
```

执行refreshGrayWenzhenDoctor，并设置分布式锁lockKey,超过60s进行企业微信机器人告警(robotLink机器人，alarmTimeSec超时时间)
```json
{"taskName": "refreshGrayWenzhenDoctor","lockKey": "refreshGrayWenzhenDoctor_task",
"lockTime": 30, "alarmTimeSec": 60, "robotLink": "0d9a679d-10f6-44fe-a210-37f0690bdb6e" }
```


## task
对应本地缓存相关任务执行