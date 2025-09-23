// Package bloom bloom
package bloom

// 参考https://github.com/bculberson/bloom
import (
	"context"
	"hash/fnv"
	"math"

	"github.com/redis/go-redis/v9"
)

// BitSetProvider bit存储提供者
type BitSetProvider interface {
	Set(ctx context.Context, location []uint) error
	Test(ctx context.Context, location []uint) (bool, error)
}

// BloomFilter 布隆过滤
type BloomFilter struct {
	m      uint
	k      uint
	bitSet BitSetProvider
}

// GetBloom 获取并初始化bloom过滤器，如果不存在则创建.dataSize 数据大小，float 误差率
func GetBloom(redisClientCli *redis.Client, bloomKey string, dataSize uint, float float64) (
	*BloomFilter, error,
) {
	m, k := EstimateParameters(dataSize, float)
	bitSet := NewRedisBitSet(bloomKey, redisClientCli)
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

// New 新增过滤器
func New(m uint, k uint, bitSet BitSetProvider) *BloomFilter {
	return &BloomFilter{m: m, k: k, bitSet: bitSet}
}

// EstimateParameters 计算大小m大小，k hash次数
func EstimateParameters(n uint, p float64) (uint, uint) {
	m := math.Ceil(float64(n) * math.Log(p) / math.Log(1.0/math.Pow(2.0, math.Ln2)))
	k := math.Ln2*m/float64(n) + 0.5
	return uint(m), uint(k)
}

// Add 添加
func (f *BloomFilter) Add(ctx context.Context, data []byte) error {
	locations := f.getLocations(data)
	err := f.bitSet.Set(ctx, locations)
	if err != nil {
		return err
	}
	return nil
}

// Exists 是否存在
func (f *BloomFilter) Exists(ctx context.Context, data []byte) (bool, error) {
	locations := f.getLocations(data)
	isSet, err := f.bitSet.Test(ctx, locations)
	if err != nil {
		return false, err
	}
	if !isSet {
		return false, nil
	}

	return true, nil
}

// getLocations 获取位置
func (f *BloomFilter) getLocations(data []byte) []uint {
	locations := make([]uint, f.k)
	hasher := fnv.New64()
	hasher.Write(data)
	a := make([]byte, 1)
	for i := uint(0); i < f.k; i++ {
		a[0] = byte(i)
		hasher.Write(a)
		hashValue := hasher.Sum64()
		locations[i] = uint(hashValue % uint64(f.m))
	}
	return locations
}
