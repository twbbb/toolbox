// Package bitset bitset
package bitset

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

// SetBit 设置
func (r *RedisBitSet) SetBit(ctx context.Context, offsets []uint) error {
	if len(offsets) == 0 {
		return nil
	}

	// 按key分组offsets
	keyOffsets := make(map[string][]uint)
	for _, offset := range offsets {
		key, thisOffset := r.getKeyOffset(offset)
		keyOffsets[key] = append(keyOffsets[key], thisOffset)
	}

	// 对每个key执行批量BITFIELD操作
	for key, offsetList := range keyOffsets {
		args := []interface{}{}
		for _, offset := range offsetList {
			args = append(args, "SET", "u1", int64(offset), 1)
		}
		_, err := r.client.BitField(ctx, key, args...).Result()
		if err != nil {
			return err
		}
	}

	return nil
}

// Exists 测试是否存在
func (r *RedisBitSet) Exists(ctx context.Context, offsets []uint) (bool, error) {
	if len(offsets) == 0 {
		return true, nil
	}

	// 按key分组offsets
	keyOffsets := make(map[string][]uint)
	for _, offset := range offsets {
		key, thisOffset := r.getKeyOffset(offset)
		keyOffsets[key] = append(keyOffsets[key], thisOffset)
	}

	// 对每个key执行批量BITFIELD操作
	for key, offsetList := range keyOffsets {
		args := []interface{}{}
		for _, offset := range offsetList {
			args = append(args, "GET", "u1", int64(offset))
		}
		results, err := r.client.BitField(ctx, key, args...).Result()
		if err != nil {
			return false, err
		}
		// 检查所有位是否都为1
		for _, result := range results {
			if result == 0 {
				return false, nil
			}
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
