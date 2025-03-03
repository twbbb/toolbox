// Package bloom bloom
package bloom

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// RedisBitSet redis的bit过滤器
type RedisBitSet struct {
	keyPrefix string
	client    *redis.Client
}

// NewRedisBitSet 新增
func NewRedisBitSet(keyPrefix string, client *redis.Client) *RedisBitSet {
	return &RedisBitSet{keyPrefix, client}
}

// Set 设置
func (r *RedisBitSet) Set(ctx context.Context, offsets []uint) error {
	for _, offset := range offsets {
		key, thisOffset := r.getKeyOffset(offset)
		_, err := r.client.SetBit(ctx, key, int64(thisOffset), 1).Result()
		if err != nil {
			return err
		}
	}

	return nil
}

// Test 测试是否存在
func (r *RedisBitSet) Test(ctx context.Context, offsets []uint) (bool, error) {
	for _, offset := range offsets {
		key, thisOffset := r.getKeyOffset(offset)
		bitValue, err := r.client.GetBit(ctx, key, int64(thisOffset)).Result()
		if err != nil {
			return false, err
		}
		if bitValue == 0 {
			return false, nil
		}
	}

	return true, nil
}

// getKeyOffset 获取key位置
func (r *RedisBitSet) getKeyOffset(offset uint) (string, uint) {
	// const redisMaxLength = 8 * 1024 * 1024 * 512
	// n := uint(offset / redisMaxLength)
	// thisOffset := offset - n*redisMaxLength
	// key := fmt.Sprintf("%s:%d", r.keyPrefix, n)
	// return key, thisOffset
	return r.keyPrefix, offset
}
