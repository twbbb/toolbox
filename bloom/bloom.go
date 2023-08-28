// Package bloom bloom
package bloom

//参考https://github.com/bculberson/bloom
import (
	"hash/fnv"
	"math"
)

//BitSetProvider bit存储提供者
type BitSetProvider interface {
	Set([]uint) error
	Test([]uint) (bool, error)
}

//BloomFilter 布隆过滤
type BloomFilter struct {
	m      uint
	k      uint
	bitSet BitSetProvider
}

//New 新增过滤器
func New(m uint, k uint, bitSet BitSetProvider) *BloomFilter {
	return &BloomFilter{m: m, k: k, bitSet: bitSet}
}

//EstimateParameters 计算大小m大小，k hash次数
func EstimateParameters(n uint, p float64) (uint, uint) {
	m := math.Ceil(float64(n) * math.Log(p) / math.Log(1.0/math.Pow(2.0, math.Ln2)))
	k := math.Ln2*m/float64(n) + 0.5

	return uint(m), uint(k)
}

//Add 添加
func (f *BloomFilter) Add(data []byte) error {
	locations := f.getLocations(data)
	err := f.bitSet.Set(locations)
	if err != nil {
		return err
	}
	return nil
}

// Exists 是否存在
func (f *BloomFilter) Exists(data []byte) (bool, error) {
	locations := f.getLocations(data)
	isSet, err := f.bitSet.Test(locations)
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
