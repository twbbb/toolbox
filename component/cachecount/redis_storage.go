package cachecount

import (
	"context"
	"fmt"
	"time"

	"git.code.oa.com/trpc-go/trpc-go/log"
	"github.com/redis/go-redis/v9"
)

// RedisStorage Redis存储操作结构体
type RedisStorage struct {
	keyPrefix   string        // Redis key前缀，实际key格式为: prefix:YYYYMMDD
	redisClient *redis.Client // Redis客户端实例
}

// NewRedisStorage 创建Redis存储实例
func NewRedisStorage(keyPrefix string, redisClient *redis.Client) *RedisStorage {
	return &RedisStorage{
		keyPrefix:   keyPrefix,
		redisClient: redisClient,
	}
}

// getRedisKey 根据日期生成Redis key
func (rs *RedisStorage) getRedisKey(date time.Time) string {
	if date.IsZero() {
		date = time.Now()
	}
	return rs.keyPrefix + ":" + date.Format("20060102")
}

// getTodayRedisKey 获取今天的Redis key
func (rs *RedisStorage) getTodayRedisKey() string {
	return rs.getRedisKey(time.Now())
}

// BatchIncrBy 批量增加计数到Redis
func (rs *RedisStorage) BatchIncrBy(ctx context.Context, data map[string]int64) error {
	if len(data) == 0 {
		return nil
	}

	todayKey := rs.getTodayRedisKey()

	pipe := rs.redisClient.Pipeline()
	for key, count := range data {
		pipe.HIncrBy(ctx, todayKey, key, count)
	}

	// 设置过期时间为7天，避免历史数据一直占用内存
	pipe.Expire(ctx, todayKey, 7*24*time.Hour)

	_, err := pipe.Exec(ctx)
	if err != nil {
		log.ErrorContextf(
			ctx, "[RedisStorage.BatchIncrBy] pipeline exec failed, key=%s, dataSize=%d, err=%v", todayKey, len(data),
			err,
		)
		return err
	}
	log.InfoContextf(
		ctx, "[RedisStorage.BatchIncrBy] success, key=%s, dataSize=%d", todayKey, len(data),
	)
	return nil
}

// GetCountsByDay 获取指定天数的所有计数数据
// days: 天数，0表示今天，1表示昨天，2表示前天，以此类推
func (rs *RedisStorage) GetCountsByDay(ctx context.Context, days int) (map[string]int64, error) {
	targetDate := time.Now().AddDate(0, 0, -days)
	return rs.GetCountsByDate(ctx, targetDate)
}

// GetCountsByDate 获取指定日期的所有计数数据
// date: 日期时间，如果为零值则使用当前时间
func (rs *RedisStorage) GetCountsByDate(ctx context.Context, date time.Time) (map[string]int64, error) {
	redisKey := rs.getRedisKey(date)
	dateStr := date.Format("20060102")

	result, err := rs.redisClient.HGetAll(ctx, redisKey).Result()
	if err != nil {
		log.ErrorContextf(
			ctx, "[RedisStorage.GetCountsByDate] HGetAll failed, key=%s, dateStr=%s, err=%v", redisKey, dateStr, err,
		)
		return nil, err
	}
	log.InfoContextf(
		ctx, "[RedisStorage.GetCountsByDate] HGetAll success, key=%s, dateStr=%s, resultSize=%d", redisKey, dateStr,
		len(result),
	)

	counts := make(map[string]int64, len(result))
	for key, value := range result {
		var count int64
		_, err := fmt.Sscanf(value, "%d", &count)
		if err != nil {
			log.WarnContextf(
				ctx, "[RedisStorage.GetCountsByDate] parse count failed, key=%s, value=%s, err=%v", key, value, err,
			)
			continue
		}
		counts[key] = count
	}

	return counts, nil
}

// DeleteByDay 删除指定天数的所有计数数据
// days: 天数，0表示今天，1表示昨天，2表示前天，以此类推
func (rs *RedisStorage) DeleteByDay(ctx context.Context, days int) error {
	targetDate := time.Now().AddDate(0, 0, -days)
	return rs.DeleteByDate(ctx, targetDate)
}

// DeleteByDate 删除指定日期的所有计数数据
// date: 日期时间，如果为零值则使用当前时间
func (rs *RedisStorage) DeleteByDate(ctx context.Context, date time.Time) error {
	redisKey := rs.getRedisKey(date)
	dateStr := date.Format("20060102")

	err := rs.redisClient.Del(ctx, redisKey).Err()
	if err != nil {
		log.ErrorContextf(ctx, "[RedisStorage.DeleteByDate] delete failed, key=%s, dateStr=%s, err=%v", redisKey, dateStr, err)
		return err
	}
	log.InfoContextf(ctx, "[RedisStorage.DeleteByDate] success, key=%s, dateStr=%s", redisKey, dateStr)
	return nil
}
