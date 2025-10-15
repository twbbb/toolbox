package cachecount

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

// 创建测试用的Redis客户端
func createTestRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // 根据实际情况修改
		DB:       0,
	})
}

// TestCacheCount_Incr 测试基本的计数功能
func TestCacheCount_Incr(t *testing.T) {
	client := createTestRedisClient()
	defer client.Close()

	// 清理测试数据
	ctx := context.Background()
	testKey := "test:cache:count"
	client.Del(ctx, testKey)

	cc := New(testKey, client)
	defer cc.Close()

	// 测试单次增加
	cc.Incr("user:1")
	if count := cc.GetLocalCount("user:1"); count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}

	// 测试多次增加
	cc.Incr("user:1")
	cc.Incr("user:1")
	if count := cc.GetLocalCount("user:1"); count != 3 {
		t.Errorf("Expected count 3, got %d", count)
	}

	// 测试不同的key
	cc.Incr("user:2")
	if count := cc.GetLocalCount("user:2"); count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}
}

// TestCacheCount_IncrBy 测试指定增量的计数功能
func TestCacheCount_IncrBy(t *testing.T) {
	client := createTestRedisClient()
	defer client.Close()

	ctx := context.Background()
	testKey := "test:cache:count:incrby"
	client.Del(ctx, testKey)

	cc := New(testKey, client)
	defer cc.Close()

	cc.IncrBy("user:1", 5)
	if count := cc.GetLocalCount("user:1"); count != 5 {
		t.Errorf("Expected count 5, got %d", count)
	}

	cc.IncrBy("user:1", 10)
	if count := cc.GetLocalCount("user:1"); count != 15 {
		t.Errorf("Expected count 15, got %d", count)
	}
}

// TestCacheCount_Concurrent 测试并发场景
func TestCacheCount_Concurrent(t *testing.T) {
	client := createTestRedisClient()
	defer client.Close()

	ctx := context.Background()
	testKey := "test:cache:count:concurrent"
	client.Del(ctx, testKey)

	cc := New(testKey, client)
	defer cc.Close()

	// 并发增加计数
	var wg sync.WaitGroup
	goroutines := 100
	increments := 100

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < increments; j++ {
				cc.Incr("user:concurrent")
			}
		}()
	}

	wg.Wait()

	expectedCount := int64(goroutines * increments)
	if count := cc.GetLocalCount("user:concurrent"); count != expectedCount {
		t.Errorf("Expected count %d, got %d", expectedCount, count)
	}
}

// TestCacheCount_Flush 测试刷新到Redis功能
func TestCacheCount_Flush(t *testing.T) {
	client := createTestRedisClient()
	defer client.Close()

	ctx := context.Background()
	testKey := "test:cache:count:flush"
	client.Del(ctx, testKey)

	cc := New(testKey, client)
	defer cc.Close()

	// 增加计数
	cc.Incr("user:1")
	cc.Incr("user:1")
	cc.IncrBy("user:2", 5)

	// 手动刷新
	cc.Flush()

	// 等待刷新完成
	time.Sleep(100 * time.Millisecond)

	// 检查本地计数是否已清空
	if count := cc.GetLocalCount("user:1"); count != 0 {
		t.Errorf("Expected local count 0 after flush, got %d", count)
	}

	// 检查Redis中的值
	redisCount1, err := client.HGet(ctx, testKey, "user:1").Int64()
	if err != nil {
		t.Fatalf("Failed to get redis count: %v", err)
	}
	if redisCount1 != 2 {
		t.Errorf("Expected redis count 2, got %d", redisCount1)
	}

	redisCount2, err := client.HGet(ctx, testKey, "user:2").Int64()
	if err != nil {
		t.Fatalf("Failed to get redis count: %v", err)
	}
	if redisCount2 != 5 {
		t.Errorf("Expected redis count 5, got %d", redisCount2)
	}
}

// TestCacheCount_AutoFlush 测试自动刷新功能
func TestCacheCount_AutoFlush(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping auto flush test in short mode")
	}

	client := createTestRedisClient()
	defer client.Close()

	ctx := context.Background()
	testKey := "test:cache:count:autoflush"
	client.Del(ctx, testKey)

	cc := New(testKey, client)
	defer cc.Close()

	// 增加计数
	cc.Incr("user:1")
	cc.IncrBy("user:2", 3)

	// 等待自动刷新（2分钟+缓冲时间）
	time.Sleep(2*time.Minute + 5*time.Second)

	// 检查本地计数是否已清空
	if count := cc.GetLocalCount("user:1"); count != 0 {
		t.Errorf("Expected local count 0 after auto flush, got %d", count)
	}

	// 检查Redis中的值
	redisCount, err := client.HGet(ctx, testKey, "user:1").Int64()
	if err != nil {
		t.Fatalf("Failed to get redis count: %v", err)
	}
	if redisCount != 1 {
		t.Errorf("Expected redis count 1, got %d", redisCount)
	}
}

// TestCacheCount_GetTotalCount 测试获取总计数功能
func TestCacheCount_GetTotalCount(t *testing.T) {
	client := createTestRedisClient()
	defer client.Close()

	ctx := context.Background()
	testKey := "test:cache:count:total"
	client.Del(ctx, testKey)

	// 先在Redis中设置一些初始值
	client.HSet(ctx, testKey, "user:1", 10)

	cc := New(testKey, client)
	defer cc.Close()

	// 增加本地计数
	cc.IncrBy("user:1", 5)

	// 获取总计数
	totalCount, err := cc.GetTotalCount("user:1")
	if err != nil {
		t.Fatalf("Failed to get total count: %v", err)
	}

	if totalCount != 15 {
		t.Errorf("Expected total count 15, got %d", totalCount)
	}
}

// BenchmarkCacheCount_Incr 性能测试
func BenchmarkCacheCount_Incr(b *testing.B) {
	client := createTestRedisClient()
	defer client.Close()

	ctx := context.Background()
	testKey := "test:cache:count:bench"
	client.Del(ctx, testKey)

	cc := New(testKey, client)
	defer cc.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			cc.Incr("user:bench")
			i++
		}
	})
}
