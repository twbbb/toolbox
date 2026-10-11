// Package bloom bloom
package bloom

// 参考https://github.com/bculberson/bloom
import (
	"context"
	"hash/fnv"
	"math"

	"github.com/twbbb/toolbox/component/bloom/bitset"
)

// BloomFilter 布隆过滤
type BloomFilter struct {
	m      uint
	k      uint
	bitSet bitset.BitSetApi
}

// New 新增过滤器
func New(m uint, k uint, bitSet bitset.BitSetApi) *BloomFilter {
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
	err := f.bitSet.SetBit(ctx, locations)
	if err != nil {
		return err
	}
	return nil
}

// Exists 是否存在
func (f *BloomFilter) Exists(ctx context.Context, data []byte) (bool, error) {
	locations := f.getLocations(data)
	isSet, err := f.bitSet.Exists(ctx, locations)
	if err != nil {
		return false, err
	}
	if !isSet {
		return false, nil
	}

	return true, nil
}

// getLocations 获取位置 - 使用双hash方法
func (f *BloomFilter) getLocations(data []byte) []uint {
	locations := make([]uint, f.k)

	// 使用两个独立的hash函数
	hasher1 := fnv.New64()
	hasher1.Write(data)
	hash1 := hasher1.Sum64()

	hasher2 := fnv.New64a()
	hasher2.Write(data)
	hash2 := hasher2.Sum64()

	// 使用双hash公式: h1(x) + i * h2(x)
	for i := uint(0); i < f.k; i++ {
		combinedHash := hash1 + uint64(i)*hash2
		locations[i] = uint(combinedHash % uint64(f.m))
	}

	return locations
}

//
// // getLocations 获取位置
// func (f *BloomFilter) getLocations(data []byte) []uint {
// 	locations := make([]uint, f.k)
// 	hasher := fnv.New64()
// 	hasher.Write(data)
// 	a := make([]byte, 1)
// 	for i := uint(0); i < f.k; i++ {
// 		a[0] = byte(i)
// 		hasher.Write(a)
// 		hashValue := hasher.Sum64()
// 		locations[i] = uint(hashValue % uint64(f.m))
// 	}
// 	return locations
// }
