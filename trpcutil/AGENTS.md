本目录是 tRPC 服务里常用的辅助：
entity 库、下游、Redis 和定时任务的公共结构
initdata 按配置拼 tRPC 和 TAF 的 client 选项
monitor 把调用次数、失败和耗时报到智研
ping 网关探活用的空 Ping，只说明进程还活着
task 本机缓存刷新、分布式定时任务，以及对接定时任务平台
trpcmeta 从 context 读写网关透传的用户身份
