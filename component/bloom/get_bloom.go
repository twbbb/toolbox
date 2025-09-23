package bloom

import (
	"context"

	"git.woa.com/th-go/toolbox/component/bloom/bitset"
	"github.com/redis/go-redis/v9"
)

// GetBloom 获取并初始化bloom过滤器，如果不存在则创建.dataSize 数据大小，float 误差率
// 注意：v0.1.0和v0.0.5的hash方法不兼容
func GetBloom(redisClientCli *redis.Client, bloomKey string, dataSize uint, float float64) (
	*BloomFilter, error,
) {
	m, k := EstimateParameters(dataSize, float)
	bitSet := bitset.NewRedisBitSet(bloomKey, redisClientCli)
	b := New(m, k, bitSet)

	// 检查是否已经存在了
	flag, err := redisClientCli.Exists(context.Background(), bloomKey).Result()
	if err != nil {
		return nil, err
	}
	// 不存在则初始化
	if flag == 0 {
		redisClientCli.SetBit(context.Background(), bloomKey, int64(m), 0)
	}

	return b, nil
}
