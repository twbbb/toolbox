本目录是可单独引用的组件：
aesutil 用 AES 加解密字符串，密文是 base64
bloom 布隆过滤器，位图放在 Redis
cachecount 内存里累加计数，定时批量刷到 Redis
redislock 用 Redis SetNX 加锁，删 key 解锁
