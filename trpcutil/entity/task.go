package entity

import (
	"context"
	"time"

	"git.code.oa.com/trpc-go/trpc-go/log"
)

// TaskFunc 定时任务结构体
type TaskFunc struct {
	// 任务名称
	TaskName string
	// 任务执行间隔
	CornDuration int
	// 任务执行函数
	TaskFunc func(ctx context.Context) error
	// 启动执行
	InitStart bool
	// Logger
	Logger log.Logger
	// 上报名称
	ReportName string
}

// DistributLock 分布式锁
type DistributLock interface {
	// Lock 锁定
	Lock(ctx context.Context, key string, expiration time.Duration) bool

	// UnLock 解锁
	UnLock(ctx context.Context, key string) bool
}

// TaskDistributFunc 分布式定时任务结构体
type TaskDistributFunc struct {
	// 任务名称
	TaskName string
	// 任务执行间隔
	CornDuration int
	// 任务执行函数
	TaskFunc func(ctx context.Context) error
	// 启动执行
	InitStart bool
	// 最大执行时间，分钟，比如30，则30分钟没有执行完，也会释放锁
	MaxRunTime int
	// 分布式锁的basekey
	BaseLockKey string
	// 分布式锁
	DistributLock DistributLock

	// Logger
	Logger log.Logger
	// 上报名称
	ReportName string
}

// PlatformTaskDto 定时任务平台消息传递
type PlatformTaskDto struct {
	TaskName     string `json:"taskName"`
	TransformMsg string `json:"transformMsg"` // 透传到定时任务的信息
	LockKey      string `json:"lockKey"`      // 加锁的key,没有则不加锁
	LockTime     int    `json:"lockTime"`     // 默认加锁的时间,单位为秒
	TimeoutCtx   int    `json:"timeoutCtx"`   // 任务超时时间，单位为min，默认没有超时时间
	AlarmTimeSec int    `json:"alarmTimeSec"` // 运行时间超过这个限制则进行告警，单位为s
	RobotLink    string `json:"robotLink"`    // 运行时间超过这个限制则进行企业微信告警
	Corn         string `json:"corn"`         // 任务的corn表达式，不影响调用，只是会上报到自定义监控
}

// PlatformTaskFunc 定时任务平台消息传递消息处理函数类型
type PlatformTaskFunc func(ctx context.Context, data string) error
type PlatformTaskFuncWithCallback func(ctx context.Context, data string) (interface{}, error)
type PlatformTaskCallbackFunc func(ctx context.Context, data interface{}) error

// PlatformTask 定时任务平台消息传递消息处理结构体
type PlatformTask struct {
	TaskName string
	// 无回调函数的方法
	Func PlatformTaskFunc
	// 有后续回掉函数的方法
	FuncWithCallBack PlatformTaskFuncWithCallback
	// 回调函数
	CallbackFunc PlatformTaskCallbackFunc
}
