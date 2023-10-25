package bodyparse

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"git.code.oa.com/trpc-go/trpc-go/log"
)

// ProxyBodyParseImpl 解析
type ProxyBodyParseImpl struct {
	SimpleBodyParseImpl
}

// GetReqStr  获取请求参数
func (b ProxyBodyParseImpl) GetReqStr(ctx context.Context, req interface{}, servieName, method string) (
	string, error,
) {
	reqJSON, err := json.Marshal(req)

	if err != nil {
		log.Errorf("json.Marshal err|%v|%s|", req, err.Error())
		return "", err
	}

	proxyMap := make(map[string]string)
	proxyMap["reqStr"] = string(reqJSON)

	return b.SimpleBodyParseImpl.GetReqStr(ctx, proxyMap, servieName, method)
}

// GetRspStr  获取响应数据
func (b ProxyBodyParseImpl) GetRspStr(ctx context.Context, response *http.Response, servieName, method string) (
	string, error,
) {

	rsp, err := b.SimpleBodyParseImpl.GetReqStr(ctx, response, servieName, method)
	if err != nil {
		return rsp, err
	}
	proxyMap := make(map[string]interface{})
	if err := json.Unmarshal([]byte(rsp), &proxyMap); err != nil {
		fmt.Println("JSON unmarshal error:", err)
		return "", err
	}
	rspStr := proxyMap["rspStr"]
	if rspStr == nil {
		return rsp, nil
	}
	result, err := json.Marshal(rspStr)
	if err != nil {
		fmt.Println("JSON marshal error:", err)
		return "", err
	}
	return string(result), nil

}
