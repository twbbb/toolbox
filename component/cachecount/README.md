# CacheCount - 高并发本地缓存计数组件

一个专为高并发场景设计的本地缓存使用次数统计组件。

## 特性

1. **高性能**: 使用本地内存缓存计数，避免频繁访问Redis
2. **并发安全**: 使用读写锁保证并发场景下的数据安全
3. **自动刷新**: 每隔2分钟自动将本地计数刷新到Redis
4. **批量写入**: 使用Redis Pipeline批量写入，提高性能
5. **优雅关闭**: 支持优雅关闭，确保数据不丢失
6. **失败重试**: 写入Redis失败时，数据会重新加回本地缓存

## 设计说明

1. 结构体初始化时，入参是Redis的key和Redis实例，并初始化map和RWMutex
2. 提供增加计数方法，入参是一个key，将计数存储到本地map
3. 每隔2分钟，将map的计数批量刷到Redis（使用Hash结构），并清理本地map
4. 完全支持高并发场景，使用读写锁保证线程安全

## 安装

```bash
go get github.com/twbbb/toolbox
```

## 快速开始

```go
package main

import (
    "github.com/redis/go-redis/v9"
    "github.com/twbbb/toolbox/component/cachecount"
)

func main() {
    // 创建Redis客户端
    redisClient := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
    defer redisClient.Close()
    
    // 创建CacheCount实例
    cc := cachecount.New("api:access:count", redisClient)
    defer cc.Close() // 确保程序退出时关闭
    
    // 增加计数
    cc.Incr("api:/user/login")
    cc.IncrBy("api:/order/list", 5)
    
    // 获取本地计数
    localCount := cc.GetLocalCount("api:/user/login")
    
    // 获取总计数（本地+Redis）
    totalCount, _ := cc.GetTotalCount("api:/user/login")
    
    // 手动触发刷新（可选）
    cc.Flush()
}
```

## API 说明

### New(redisKey string, redisClient *redis.Client) *CacheCount
创建一个新的CacheCount实例。
- `redisKey`: Redis中存储计数的key（使用Hash结构）
- `redisClient`: Redis客户端实例

### Incr(key string)
增加指定key的计数（+1）。
- `key`: 业务key

### IncrBy(key string, delta int64)
增加指定key的计数（指定增量）。
- `key`: 业务key
- `delta`: 增量值

### GetLocalCount(key string) int64
获取本地缓存中指定key的计数（不包含Redis中的值）。

### GetRedisCount(key string) (int64, error)
从Redis中获取指定key的计数。

### GetTotalCount(key string) (int64, error)
获取指定key的总计数（本地+Redis）。

### Flush()
手动触发刷新到Redis。

### Close()
关闭CacheCount，停止后台刷新并执行最后一次刷新。

## 使用场景

- API访问次数统计
- 用户行为计数
- 资源使用统计
- 任何需要高并发计数的场景

## 注意事项

1. 必须在程序退出时调用`Close()`方法，确保数据刷新到Redis
2. Redis中使用Hash结构存储，field为业务key，value为计数
3. 默认每2分钟自动刷新一次，可以通过`Flush()`手动触发
4. 写入Redis失败时，数据会保留在本地缓存，下次刷新时重试

## 性能

在高并发场景下，相比直接操作Redis，性能提升显著：
- 本地内存操作，延迟极低
- 批量写入Redis，减少网络开销
- 读写锁优化，支持高并发读取