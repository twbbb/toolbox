// Package rpcproxy 给本地下游调用加一层 JSON 文件缓存。
// 命中缓存时直接返回，不再发起调用。未命中时走原来的下游调用，成功后再写入缓存。
// 本地调用已经是 ip 端口，这里不改寻址，也不转发。
package rpcproxy

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync/atomic"

	"trpc.group/trpc-go/trpc-go/codec"
	"trpc.group/trpc-go/trpc-go/filter"
	"trpc.group/trpc-go/trpc-go/log"
)

// FilterName 在 trpc_local.yaml 的 client.filter 中配置，需放在最后一个
const FilterName = "rpcProxy"

// skipKeywords 存储类客户端不缓存，直接放行
var skipKeywords = []string{"mysql", "redis", "mongo", "kafka", "tdmq"}

var cacheKeyIgnoreRegs = []*regexp.Regexp{
	regexp.MustCompile(`(?m)("nonce":\s*")[^"]*(".*"|\s*)(,\s*)?`),
	regexp.MustCompile(`(?m)("timestamp":\s*")[^"]*(".*"|\s*)(,\s*)?`),
	regexp.MustCompile(`(?m)("signature":\s*")[^"]*(".*"|\s*)(,\s*)?`),
}

// cacheKeyMetaKeys 参与缓存 key 的用户身份 metadata，requestId 等每次都变的不参与
var cacheKeyMetaKeys = []string{"sthid", "openId", "userId", "HCUnionId", "soid"}

// Config rpcproxy filter 配置
type Config struct {
	Enable        bool   `yaml:"enable"`        // 是否启用缓存，默认关，显式配 true 才生效
	LocalFilePath string `yaml:"localFilePath"` // 请求和结果缓存目录，空则不缓存
	Force         bool   `yaml:"force"`         // true 时不读缓存，调用后仍写缓存
	// NoCacheMethods 不走缓存的接口，每次都真实调用且不写缓存；
	// 写方法名对所有服务生效，写 "服务名/方法名" 只对该服务生效
	NoCacheMethods []string `yaml:"noCacheMethods"`
}

// conf 当前生效的配置，存 *Config；未设置时不缓存
var conf atomic.Value

// SetConfig 业务服务加载配置后设置，如 rpcproxy.SetConfig(config.GetConf().RpcProxyArg)；
// 业务配置热更新时再调一次即可生效
func SetConfig(c Config) {
	conf.Store(&c)
}

func init() {
	if os.Getenv("DOCKER_ENV") != "" {
		return
	}
	filter.Register(FilterName, filter.NoopServerFilter, clientFilter)
}

func clientFilter(ctx context.Context, req, rsp interface{}, next filter.ClientHandleFunc) error {
	msg := codec.Message(ctx)
	serviceName := msg.CalleeServiceName()
	if isSkipService(serviceName) {
		return next(ctx, req, rsp)
	}
	c, _ := conf.Load().(*Config)
	if c == nil || !c.Enable {
		return next(ctx, req, rsp)
	}
	arg := *c
	method := methodName(msg)

	reqJSON, err := json.Marshal(req)
	if err != nil {
		log.ErrorContextf(ctx, "rpcproxy cache|marshal req fail, skip cache|%s|%s|%v", serviceName, method, err)
		return next(ctx, req, rsp)
	}
	file := cacheFile(ctx, &arg, serviceName, method, reqJSON, msg.ClientMetaData())
	if data, ok := readCache(ctx, &arg, file); ok {
		if err := json.Unmarshal(data, rsp); err != nil {
			log.ErrorContextf(ctx, "rpcproxy cache|rsp unmarshal fail, ignore|%s|%v", file, err)
		} else {
			return nil
		}
	}
	if err := next(ctx, req, rsp); err != nil {
		return err
	}
	if file == "" {
		return nil
	}
	rspJSON, err := json.Marshal(rsp)
	if err != nil {
		log.ErrorContextf(ctx, "rpcproxy cache|marshal rsp fail, skip cache|%s|%v", file, err)
		return nil
	}
	writeCache(ctx, &arg, file, reqJSON, rspJSON)
	return nil
}

func methodName(msg codec.Msg) string {
	method := path.Base(msg.ClientRPCName())
	if method == "" || method == "." || method == "/" {
		return msg.CalleeMethod()
	}
	return method
}

func isSkipService(serviceName string) bool {
	lower := strings.ToLower(serviceName)
	for _, k := range skipKeywords {
		if strings.Contains(lower, k) {
			return true
		}
	}
	return false
}

// cacheEntry 缓存文件内容，请求和结果放在同一个文件里
type cacheEntry struct {
	Req json.RawMessage `json:"req"`
	Rsp json.RawMessage `json:"rsp"`
}

// cacheFile LocalFilePath 为空或接口配置为不缓存时返回空路径，表示不缓存
func cacheFile(
	ctx context.Context, arg *Config, serviceName, method string, reqJSON []byte, md codec.MetaData,
) string {
	if arg.LocalFilePath == "" {
		return ""
	}
	if isNoCache(arg, serviceName, method) {
		log.InfoContextf(ctx, "rpcproxy cache|disabled by noCacheMethods|%s/%s", serviceName, method)
		return ""
	}
	return filepath.Join(arg.LocalFilePath, serviceName+"."+method+"_"+cacheKey(reqJSON, md)+".json")
}

// isNoCache 配置项为方法名时匹配所有服务，为 "服务名/方法名" 时只匹配该服务
func isNoCache(arg *Config, serviceName, method string) bool {
	full := serviceName + "/" + method
	for _, m := range arg.NoCacheMethods {
		if m == method || m == full {
			return true
		}
	}
	return false
}

// readCache 返回缓存中的结果 JSON；日志带完整文件路径，调试时可直接改其中的 rsp
func readCache(ctx context.Context, arg *Config, file string) ([]byte, bool) {
	if file == "" {
		return nil, false
	}
	if arg.Force {
		log.InfoContextf(ctx, "rpcproxy cache|force, skip read|%s", file)
		return nil, false
	}
	data, err := os.ReadFile(file)
	if err != nil {
		log.InfoContextf(ctx, "rpcproxy cache|miss|%s", file)
		return nil, false
	}
	entry := &cacheEntry{}
	if err := json.Unmarshal(data, entry); err != nil || len(entry.Rsp) == 0 {
		log.ErrorContextf(ctx, "rpcproxy cache|invalid, ignore|%s|%v", file, err)
		return nil, false
	}
	log.InfoContextf(ctx, "rpcproxy cache|hit|%s", file)
	return entry.Rsp, true
}

func writeCache(ctx context.Context, arg *Config, file string, reqJSON, rspJSON []byte) {
	if file == "" {
		return
	}
	data, err := json.MarshalIndent(&cacheEntry{Req: reqJSON, Rsp: rspJSON}, "", "  ")
	if err != nil {
		log.ErrorContextf(ctx, "rpcproxy cache|marshal entry fail|%s|%v", file, err)
		return
	}
	if err := os.MkdirAll(arg.LocalFilePath, 0755); err != nil {
		log.ErrorContextf(ctx, "rpcproxy cache|mkdir fail|%s|%v", file, err)
		return
	}
	if err := os.WriteFile(file, data, 0644); err != nil {
		log.ErrorContextf(ctx, "rpcproxy cache|write fail|%s|%v", file, err)
		return
	}
	log.InfoContextf(ctx, "rpcproxy cache|write|%s", file)
}

// cacheKey 去掉每次都会变的签名字段后，连同用户身份 metadata 取 md5；
// 下游按 metadata 里的 sthid 等识别用户，不带上会让不同用户命中同一份缓存
func cacheKey(reqJSON []byte, md codec.MetaData) string {
	s := string(reqJSON)
	for _, re := range cacheKeyIgnoreRegs {
		s = re.ReplaceAllString(s, `$1$2$3`)
	}
	for _, k := range cacheKeyMetaKeys {
		if v, ok := md[k]; ok {
			s += "|" + k + "=" + string(v)
		}
	}
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
