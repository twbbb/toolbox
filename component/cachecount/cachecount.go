package cachecount

import (
	"context"
	"fmt"
	"sync"

	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.woa.com/th-go/toolbox/trpcutil/task"
	"github.com/redis/go-redis/v9"
)

// CacheCount 高并发场景下的本地缓存使用次数统计组件
type CacheCount struct {
	storage  *RedisStorage    // Redis存储实例
	countMap map[string]int64 // 本地计数map，key为业务key，value为计数
	mu       sync.RWMutex     // 读写锁，保证并发安全
}

// New 创建一个新的CacheCount实例
// redisKeyPrefix: Redis key前缀，实际存储的key格式为: prefix:YYYYMMDD（使用Hash结构，按天存储）
// redisClient: Redis客户端实例
func New(redisKeyPrefix string, redisClient *redis.Client) (*CacheCount, error) {
	cc := &CacheCount{
		storage:  NewRedisStorage(redisKeyPrefix, redisClient),
		countMap: make(map[string]int64),
	}
	taskFunc := task.GetDefaultTaskFunc(cc.flush, "flushLoop"+redisKeyPrefix, 60*2)
	taskFunc.InitStart = false
	err := task.DoTask(taskFunc)
	if err != nil {
		log.Errorf("[CacheCount.New] init cache ticker failed, prefix=%s, err=%v", redisKeyPrefix, err)
		return nil, err
	}

	log.Infof("[CacheCount.New] success, prefix=%s", redisKeyPrefix)
	return cc, nil
}
func (cc *CacheCount) IncrId(key int) {
	cc.Incr(fmt.Sprintf("%d", key))
}

// Incr 增加指定key的计数
// key: 业务key
func (cc *CacheCount) Incr(key string) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	cc.countMap[key]++
}

// IncrBy 增加指定key的计数（指定增量）
// key: 业务key
// delta: 增量值
func (cc *CacheCount) IncrBy(key string, delta int64) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	cc.countMap[key] += delta
}

// flush 将本地计数刷新到Redis并清理本地map
func (cc *CacheCount) flush(ctx context.Context) error {
	// 如果没有数据，直接返回
	if len(cc.countMap) == 0 {
		log.DebugContextf(ctx, "[CacheCount.flush] no data to flush")
		return nil
	}

	// 先创建新的map，然后交换并立即释放锁
	cc.mu.Lock()
	localData := cc.countMap
	dataSize := len(localData)
	cc.countMap = make(map[string]int64)
	cc.mu.Unlock()

	// log.InfoContextf(ctx, "[CacheCount.flush] start flush, dataSize=%d", dataSize)
	err := cc.storage.BatchIncrBy(ctx, localData)
	if err != nil {
		log.ErrorContextf(ctx, "[CacheCount.flush] flush to redis failed, dataSize=%d, err=%v", dataSize, err)
		return err
	}
	log.InfoContextf(ctx, "[CacheCount.flush] flush success, dataSize=%d", dataSize)
	return nil
}

// Flush 手动触发刷新到Redis
func (cc *CacheCount) Flush(ctx context.Context) error {
	return cc.flush(ctx)
}

// GetLocalCount 获取本地map中指定key的计数
// key: 业务key
func (cc *CacheCount) GetLocalCount(key string) int64 {
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	return cc.countMap[key]
}

// GetRedisCount 从Redis中获取指定天数的所有计数数据
// days: 天数，0表示今天，1表示昨天，2表示前天，以此类推
// 返回: map[string]int64，key为业务key，value为计数
func (cc *CacheCount) GetRedisCount(ctx context.Context, days int) (map[string]int64, error) {
	log.InfoContextf(ctx, "[CacheCount.GetRedisCount] start, days=%d", days)
	counts, err := cc.storage.GetCountsByDay(ctx, days)
	if err != nil {
		log.ErrorContextf(ctx, "[CacheCount.GetRedisCount] failed, days=%d, err=%v", days, err)
		return nil, err
	}
	log.InfoContextf(ctx, "[CacheCount.GetRedisCount] success, days=%d, countsSize=%d", days, len(counts))
	return counts, nil
}
