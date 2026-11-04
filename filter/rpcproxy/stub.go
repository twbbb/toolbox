package rpcproxy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"trpc.group/trpc-go/trpc-go/codec"
)

// 本地用例写替身用：文件名和 filter、Transport 用同一份算法，替身一定命中。

// StubFile trpc 下游一次调用的替身文件路径。serviceName 是被调服务名（如 trpc.common.auth.Auth），
// method 是方法名，req 是被测服务发出的请求，md 是它带的用户身份 metadata（没有就传 nil）。
func StubFile(dir, serviceName, method string, req interface{}, md map[string]string) (string, error) {
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("marshal req: %w", err)
	}
	meta := codec.MetaData{}
	for k, v := range md {
		meta[k] = []byte(v)
	}
	return filepath.Join(dir, serviceName+"."+method+"_"+cacheKey(reqJSON, meta)+".json"), nil
}

// HTTPStubFile HTTP 下游一次调用的替身文件路径。rawURL 带 query，body 是请求体原文。
func HTTPStubFile(dir, method, rawURL string, body []byte) (string, error) {
	req, err := http.NewRequest(strings.ToUpper(method), rawURL, nil)
	if err != nil {
		return "", fmt.Errorf("parse url: %w", err)
	}
	name := httpServicePrefix + req.URL.Host + "." + httpMethodName(req.URL.Path)
	return filepath.Join(dir, name+"_"+httpCacheKey(req, body)+".json"), nil
}

// WriteStub 把 req、rsp 写成替身文件，格式和缓存文件一样。HTTP 下游的 rsp 用 *HTTPRsp。
func WriteStub(file string, req, rsp interface{}) error {
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal req: %w", err)
	}
	rspJSON, err := json.Marshal(rsp)
	if err != nil {
		return fmt.Errorf("marshal rsp: %w", err)
	}
	data, err := json.MarshalIndent(&cacheEntry{Req: reqJSON, Rsp: rspJSON}, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		return err
	}
	return os.WriteFile(file, data, 0644)
}
