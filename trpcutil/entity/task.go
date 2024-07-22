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
