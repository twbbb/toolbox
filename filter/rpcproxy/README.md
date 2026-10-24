# rpcproxy 本地调用缓存

本地启动服务时，下游已经按 ip 端口直连。`rpcProxy` 是一个 client filter，把请求和结果缓存成 JSON 文件。命中缓存时直接返回，不再发起调用。未命中时走原来的下游调用，成功后再写入缓存。

这里不改寻址，也不把调用转发到测试环境。

`DOCKER_ENV` 非空（容器环境）时不注册，线上不生效。

## 1. 哪些调用会进缓存

| 调用 | 行为 |
|---|---|
| 普通下游 RPC | 命中缓存直接返回；未命中走原调用，成功后写缓存 |
| mysql / redis / mongo / kafka / tdmq | 服务名含这些关键字的直接放行，不缓存 |

## 2. 开启

业务配置结构里加一个 `rpcproxy.Config` 字段，加载配置后传给 rpcproxy。热更新时再调一次即可。

```go
import "github.com/twbbb/toolbox/filter/rpcproxy"

type Properties struct {
	RpcProxyArg rpcproxy.Config `yaml:"rpcProxyArg"`
}

rpcproxy.SetConfig(config.GetConf().RpcProxyArg)
```

未调用 `SetConfig`，或 `enable` 不是 true 时，filter 直接放行。

以下以 th_jingqi_server 的本地启动目录为例。

**框架配置** `localstart/config/local/trpc_local.yaml`：`rpcProxy` 放在 client filter 的最后一个。

```yaml
client:
  filter:
    - rpcProxy
```

**业务配置** `localstart/config/local/app.yaml`：

```yaml
rpcProxyArg:
  enable: true          # 必须显式开启
  localFilePath: /Users/.../th_jingqi_server/localstart/rpcproxy_cache   # 缓存目录，空则不缓存；用绝对路径
  force: false          # true：所有接口都不读缓存，调用后仍写缓存
  noCacheMethods:       # 这些接口不读也不写缓存，每次都真实调用
    - GetMemberList
    - trpc.a.b.S/GetMemberInfo
```

`localFilePath` 要写绝对路径。服务的工作目录是 `localstart/config/local/`，写相对路径会落到那下面。

缓存目录放在 `.gitignore` 里（`/localstart/rpcproxy_cache/`）。缓存里是下游返回的数据，不要提交。

改完配置要重启服务：`./localstart/startserver.sh restart`。改缓存文件本身不用重启。

## 3. 缓存文件

### 文件名

```
<服务名>.<方法名>_<md5>.json
例：trpc.a.b.S.GetMemberList_b5309a91a874ce82330471336652dd8a.json
```

方法名用 proto 里的方法名。`noCacheMethods` 也按这个写法填。

### 文件内容

请求和结果放在同一个文件：

```json
{
  "req": { ...请求... },
  "rsp": { ...结果... }
}
```

命中缓存时只用 `rsp`。改 `rsp` 就是改下游返回值。

`req` 只用来查看。改它不影响命中，因为文件名在写入时就已经算好了。

### 缓存 key

md5 由两部分组成：

1. 请求 JSON，去掉 `nonce`、`timestamp`、`signature` 这几个每次都会变的字段。
2. 用户身份 metadata：`sthid`、`openId`、`userId`、`HCUnionId`、`soid`。

所以：

- 请求一样、用户一样，命中同一个文件。`requestId` 变化不影响。
- 不同用户（不同 sthid）是不同的文件。
- 请求参数变了就是新文件，旧文件不会被覆盖。

### 什么时候写缓存

只有 `next` 不返回错误才写。下游业务错误码只要调用本身成功，也会照常缓存。

缓存不会自动过期。下游数据变了，删掉对应文件，或者临时开 `force`，或者把接口加进 `noCacheMethods`。

## 4. 日志

```bash
grep "rpcproxy cache" localstart/log/trpc.log | tail
```

| 日志 | 含义 |
|---|---|
| `rpcproxy cache\|miss\|<文件>` | 文件不存在，走原来的下游调用 |
| `rpcproxy cache\|write\|<文件>` | 调用成功，已写入这个文件 |
| `rpcproxy cache\|hit\|<文件>` | 用了这个文件的 `rsp`，没有调下游 |
| `rpcproxy cache\|force, skip read\|<文件>` | 开了 `force`，不读缓存，调用后仍写 |
| `rpcproxy cache\|disabled by noCacheMethods\|<服务>/<方法>` | 接口配置为不缓存 |
| `rpcproxy cache\|invalid, ignore\|<文件>\|<原因>` | 文件不是合法 JSON 或没有 `rsp`，改为走真实调用 |
| `rpcproxy cache\|rsp unmarshal fail, ignore\|<文件>\|<原因>` | `rsp` 的字段类型和接口对不上，改为走真实调用 |

后两种是手改文件写错时会出现的情况。请求会回退到真实调用，看日志里的原因修正即可。

## 5. 常见用法

### 改下游返回值调试

1. 调一次接口，日志出现 `cache|write|<文件>`。
2. 打开这个文件，改 `rsp` 里的字段。
3. 再调一次接口，日志出现 `cache|hit|<文件>`，返回的就是改过的数据。不用重启，每次请求都会重新读文件。

想恢复成真实数据，删掉这个文件。下次调用会重新请求下游并写入。

### 某个接口要实时数据

把它加进 `noCacheMethods` 后重启。只影响这个接口，其他接口照常缓存。

### 刷新全部缓存

清空缓存目录：`rm localstart/rpcproxy_cache/*.json`。也可以临时开 `force: true` 跑一遍，跑完关掉。`force` 期间仍会写缓存，所以跑一遍之后缓存就是最新的。

## 6. 排查

| 现象 | 原因 |
|---|---|
| 没有任何 `rpcproxy` 日志 | 接口本身没有下游调用，或者 `enable` 没开、`trpc_local.yaml` 没配 `rpcProxy` |
| 一直命中旧数据 | 缓存不过期。删文件，或用 `force` / `noCacheMethods` |
| 目录里没有缓存文件 | 没配 `localFilePath`、接口在 `noCacheMethods` 里，或者调用失败。失败不写缓存 |

## 7. 代码位置

`rpcproxy.go`：配置、filter 入口、缓存读写、缓存 key。
