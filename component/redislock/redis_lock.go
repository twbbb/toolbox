package redislock

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// TaskRedisLock 定时任务redis锁
type TaskRedisLock struct {
	// redis客户端
	Client *redis.Client
}

// Lock 分布式锁
func (rc *TaskRedisLock) Lock(ctx context.Context, key string, expiration time.Duration) bool {
	bool, err := rc.Client.SetNX(ctx, key, "1", expiration).Result()
	if err != nil {
		return false
	}
	return bool
}

// UnLock 分布式锁解锁
func (rc *TaskRedisLock) UnLock(ctx context.Context, key string) bool {
	_, err := rc.Client.Del(ctx, key).Result()
	if err != nil {
		return false
	}

	return true
}
