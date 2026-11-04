package bloom

import (
	"fmt"
	"testing"
)

// GetBloom 获取并初始化bloom过滤器，如果不存在则创建.dataSize 数据大小，float 误差率
// 注意：v0.1.0和v0.0.5的hash方法不兼容
// https://krisives.github.io/bloom-calculator/
func TestBloomHash(t *testing.T) {
	m, k := EstimateParameters(10000000, 0.007812)

	// 计算MB大小 (1字节 = 8比特，1MB = 1024*1024字节)
	sizeInBytes := float64(m) / 8
	sizeInMB := sizeInBytes / (1024 * 1024)

	fmt.Printf("布隆过滤器参数计算结果:\n")
	fmt.Printf("- 预期元素数量: 10,000,000\n")
	fmt.Printf("- 期望误判率: 1%%\n")
	fmt.Printf("- m (位数组大小): %d 比特 (%.2f MB)\n", m, sizeInMB)
	fmt.Printf("- k (哈希函数个数): %d 个 (每个元素需要用%d个不同的哈希函数计算位置)\n", k, k)
}

// 拿到对应hash次数时候，最小的P - 数学计算法
func TestBloomMinP(t *testing.T) {
	targetK := uint(7)
	n := uint(10000000) // 预期元素数量

	fmt.Printf("通过数学公式计算k=%d时的最优误判率:\n", targetK)
	fmt.Printf("========================================\n")

	ln2 := 0.6931471805599453 // math.Ln2

	// 计算k=targetK时的最优位数组大小
	optimalM := float64(targetK) * float64(n) / ln2

	// 计算对应的最优误判率: p = (0.5)^k
	optimalP := 1.0
	for i := uint(0); i < targetK; i++ {
		optimalP *= 0.5
	}

	fmt.Printf("📐 数学公式推导:\n")
	fmt.Printf("   k = ln(2) * m/n\n")
	fmt.Printf(
		"   当k=%d时，最优m = k*n/ln(2) = %d * %d / %.6f = %.0f\n",
		targetK, targetK, n, ln2, optimalM,
	)
	fmt.Printf("   最优误判率 p = (0.5)^k = (0.5)^%d = %.6f\n", targetK, optimalP)

}
