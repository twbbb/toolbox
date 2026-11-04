package bloom

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	"github.com/twbbb/toolbox/component/bloom/bitset"
)

// GetBloom 获取并初始化bloom过滤器，如果不存在则创建.dataSize 数据大小，float 误差率
// 注意：v0.1.0和v0.0.5的hash方法不兼容
// https://krisives.github.io/bloom-calculator/
func GetBloom(redisClientCli *redis.Client, bloomKey string, dataSize uint, float float64) (
	*BloomFilter, bool, error,
) {
	if bloomKey == "" {
		return nil, false, errors.New("bloomKey is empty")
	}
	m, k := EstimateParameters(dataSize, float)
	bitSet := bitset.NewRedisBitSet(bloomKey, redisClientCli)
	b := New(m, k, bitSet)

	// 检查是否已经存在了
	flag, err := redisClientCli.Exists(context.Background(), bloomKey).Result()
	if err != nil {
		return nil, false, err
	}
	initFlag := false
	// 不存在则初始化
	if flag == 0 {
		err = redisClientCli.SetBit(context.Background(), bloomKey, int64(m), 0).Err()
		if err != nil {
			return nil, false, err
		}
		initFlag = true
	}

	return b, initFlag, nil
}
