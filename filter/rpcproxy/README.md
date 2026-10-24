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
  readOnly: false       # true：只读不写。本地用例用：缓存目录里只有用例写的替身，没命中的照常真实调用
  noCacheMethods:       # 这些接口不读也不写缓存，每次都真实调用
    - GetMemberList
    - trpc.a.b.S/GetMemberInfo
```

`localFilePath` 要写绝对路径。服务的工作目录是 `localstart/config/local/`，写相对路径会落到那下面。

缓存目录放在 `.gitignore` 里（`/localstart/rpcproxy_cache/`）。缓存里是下游返回的数据，不要提交。

改完配置要重启服务：`./localstart/startserver.sh restart`。改缓存文件本身不用重启。

### HTTP 下游

用 `net/http` 直接调的外部接口（如微信），构造 `http.Client` 时套一层 `Transport`，和 trpc 版共用同一份 `SetConfig` 配置、缓存目录和日志：

```go
client := &http.Client{Timeout: 5 * time.Second, Transport: rpcproxy.Transport(nil)} // nil 用 http.DefaultTransport
```

`DOCKER_ENV` 非空时 `Transport` 直接返回传入的 base，线上不生效。HTTP 版没有用户身份 metadata；只要调用本身没出错就写缓存，HTTP 状态码不是 2xx 也照常写。

## 3. 缓存文件

### 文件名

```
<服务名>.<方法名>_<md5>.json
例：trpc.a.b.S.GetMemberList_b5309a91a874ce82330471336652dd8a.json
```

方法名用 proto 里的方法名。`noCacheMethods` 也按这个写法填。

HTTP 下游的服务名是 `http.<host>`，方法名是 path 去掉首尾 `/` 后把 `/` 换成 `_`：

```
http.api.weixin.qq.com.cgi-bin_stable_token_<md5>.json
noCacheMethods: [cgi-bin_stable_token]  或  [http.api.weixin.qq.com/cgi-bin_stable_token]
```

### 文件内容

请求和结果放在同一个文件：

```json
{
  "req": { ...请求... },
  "rsp": { ...结果... }
}
```

HTTP 下游的文件内容：

```json
{
  "req": {"method": "POST", "url": "https://api.weixin.qq.com/cgi-bin/stable_token", "body": {...}},
  "rsp": {"status": 200, "header": {"Content-Type": "application/json"}, "body": {...}}
}
```

body 是 JSON 时原样放在 `body`，不是 JSON 时放在 `bodyText` 字符串里。手写替身时 `rsp` 只写 `body` 也行，`status` 缺省为 200。

命中缓存时只用 `rsp`。改 `rsp` 就是改下游返回值。

`req` 只用来查看。改它不影响命中，因为文件名在写入时就已经算好了。

### 缓存 key

md5 由两部分组成：

1. 请求 JSON，去掉 `nonce`、`timestamp`、`signature` 这几个每次都会变的字段。
2. 用户身份 metadata：`sthid`、`openId`、`userId`、`HCUnionId`、`soid`。

HTTP 下游的 md5 由 method、host、path、query（去掉 `nonce`、`timestamp`、`signature`，按参数名排序）和 body（同样去掉这三个字段）组成。`access_token` 这类参数参与 key：同一个接口带旧 token 和新 token 命中不同的文件，可以分别写替身。

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

### 本地用例写替身

本地用例开 `readOnly: true`，缓存目录里只放用例自己写的替身：外部下游命中替身，兄弟服务没有替身就真实调用，也不会被写进缓存。替身文件名用和 filter 同一份算法算，一定命中：

```go
// trpc 下游：被调服务名、方法名、被测服务发出的请求、用户身份 metadata
file, _ := rpcproxy.StubFile(dir, "trpc.common.auth.Auth", "Login", req, nil)
_ = rpcproxy.WriteStub(file, req, rsp)

// HTTP 下游：method、带 query 的 URL、请求体
file, _ = rpcproxy.HTTPStubFile(dir, "POST", "https://api.weixin.qq.com/cgi-bin/stable_token", body)
_ = rpcproxy.WriteStub(file, &rpcproxy.HTTPReq{Method: "POST", URL: u}, &rpcproxy.HTTPRsp{Body: json.RawMessage(`{"errcode":0}`)})
```

日志出现 `cache|hit|<该文件>`，说明被测服务实际发出的请求和写替身时给的请求一致。

## 6. 排查

| 现象 | 原因 |
|---|---|
| 没有任何 `rpcproxy` 日志 | 接口本身没有下游调用，或者 `enable` 没开、`trpc_local.yaml` 没配 `rpcProxy` |
| 一直命中旧数据 | 缓存不过期。删文件，或用 `force` / `noCacheMethods` |
| 目录里没有缓存文件 | 没配 `localFilePath`、接口在 `noCacheMethods` 里，或者调用失败。失败不写缓存 |

## 7. 代码位置

`rpcproxy.go`：配置、filter 入口、缓存读写、缓存 key。

`http.go`：HTTP 下游的 `Transport`、缓存 key、响应拼装。

`stub.go`：本地用例写替身用的 `StubFile`、`HTTPStubFile`、`WriteStub`。
