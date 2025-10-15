package cachecount

import (
	"context"
	"fmt"
	"strconv"

	"git.code.oa.com/trpc-go/trpc-go/log"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// GetRedisSaveDBByTable 从Redis获取昨天的计数数据并通过table模式保存到数据库
// redisKeyPrefix: Redis key前缀
// redisClient: Redis客户端
// db: 数据库连接
// table: 数据库表名
// day: 代表昨天
// 返回影响的行数和错误
func GetRedisSaveDBByTable(
	ctx context.Context,
	redisKeyPrefix string,
	redisClient redis.UniversalClient,
	db *gorm.DB,
	table string,
	day int,
) (int64, error) {
	return getRedisSaveDB(ctx, redisKeyPrefix, redisClient, db, table, "", day)
}

// GetRedisSaveDBBySQL 从Redis获取昨天的计数数据并通过自定义SQL保存到数据库
// redisKeyPrefix: Redis key前缀
// redisClient: Redis客户端
// db: 数据库连接
// sql: 自定义SQL语句（格式: UPDATE xxx SET f_usage_count = f_usage_count + ?, f_usage_time = ? WHERE f_id = ?）
// day: 代表昨天
// 返回影响的行数和错误
func GetRedisSaveDBBySQL(
	ctx context.Context,
	redisKeyPrefix string,
	redisClient redis.UniversalClient,
	db *gorm.DB,
	sql string,
	day int,
) (int64, error) {
	return getRedisSaveDB(ctx, redisKeyPrefix, redisClient, db, "", sql, day)
}

// GetRedisSaveDB 从Redis获取昨天的计数数据并保存到数据库
// redisKeyPrefix: Redis key前缀
// redisClient: Redis客户端
// db: 数据库连接
// table: 数据库表名（如果使用table模式）
// day代表昨天
// sql: 自定义SQL语句（如果使用sql模式，格式: UPDATE xxx SET f_usage_count = f_usage_count + ?, f_usage_time = ? WHERE f_id = ?）
// 返回影响的行数和错误
func getRedisSaveDB(
	ctx context.Context, redisKeyPrefix string, redisClient redis.UniversalClient, db *gorm.DB, table string,
	sql string,
	day int,
) (int64, error) {
	// 1. 从Redis获取昨天的计数数据
	countRedis := NewRedisStorage(redisKeyPrefix, redisClient)
	lastDayCount, err := countRedis.GetCountsByDay(ctx, day) // -1 表示昨天
	if err != nil {
		log.ErrorContextf(ctx, "GetCountsByDay error: %v", err)
		return 0, fmt.Errorf("get counts from redis failed: %w", err)
	}

	// 2. 检查是否有数据
	if len(lastDayCount) == 0 {
		log.InfoContextf(ctx, "No data found in redis for yesterday")
		return 0, nil
	}

	// 3. 转换数据格式：map[string]int64 -> map[int]int64
	// Redis中的key是字符串类型的ID，需要转换为int类型
	usageData := make(map[int]int64, len(lastDayCount))
	for key, count := range lastDayCount {
		fid, err := strconv.Atoi(key)
		if err != nil {
			log.WarnContextf(ctx, "Invalid key format, skip: %s, error: %v", key, err)
			continue
		}
		usageData[fid] = count
	}

	// 4. 检查转换后是否有有效数据
	if len(usageData) == 0 {
		log.WarnContextf(ctx, "No valid data after conversion")
		return 0, nil
	}

	// 5. 构造CountInsertDb并批量更新数据库
	countInsertDb := &CountInsertDb{
		table:     table,
		db:        db,
		sql:       sql,
		usageData: usageData,
	}

	effectRow, err := countInsertDb.BatchUpdate(ctx)
	if err != nil {
		log.ErrorContextf(ctx, "BatchUpdate error: %v", err)
		return 0, fmt.Errorf("batch update to db failed: %w", err)
	}

	log.InfoContextf(ctx, "GetRedisSaveDB success, total items: %d, affected rows: %d", len(usageData), effectRow)
	return effectRow, nil
}
