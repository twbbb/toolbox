package bodyparse

import (
	"context"
	"net/http"
)

// BodyParse 请求和响应解析
type BodyParse interface {
	// GetReqStr  获取请求参数
	GetReqStr(ctx context.Context, req interface{}, servieName, method string) (string, error)
	// GetReqKey 获取req存储key，本地存储
	GetReqKey(ctx context.Context, req string) string

	// GetRspStr  获取响应数据
	GetRspStr(ctx context.Context, rsp *http.Response, servieName, method string) (string, error)
}

var simpleBodyParse BodyParse = &SimpleBodyParseImpl{}

var proxyBodyParse BodyParse = &ProxyBodyParseImpl{}

// GetSimpleBodyParseImpl 获取一个用户信息对象
func GetSimpleBodyParseImpl() BodyParse {
	return simpleBodyParse
}

// GetProxyBodyParse 获取一个用户信息对象
func GetProxyBodyParse() BodyParse {
	return proxyBodyParse
}
