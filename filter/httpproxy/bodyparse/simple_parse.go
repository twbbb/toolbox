package bodyparse

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"git.code.oa.com/trpc-go/trpc-go/log"
)

// SimpleBodyParseImpl 解析
type SimpleBodyParseImpl struct {
}

// GetReqKey 获取req存储key，本地存储
func (b SimpleBodyParseImpl) GetReqKey(ctx context.Context, req interface{}, reqJsonStr string) string {
	reqJSON, err := json.Marshal(req)

	if err != nil {
		log.Errorf("json.Marshal err|%v|%s|", req, err.Error())

	} else {
		reqJsonStr = string(reqJSON)
	}
	// 定义正则表达式
	re := regexp.MustCompile(`(?m)("nonce":\s*")[^"]*(".*"|\s*)(,\s*)?`)
	reqJsonStr = re.ReplaceAllString(reqJsonStr, `$1$2$3`)
	re = regexp.MustCompile(`(?m)("timestamp":\s*")[^"]*(".*"|\s*)(,\s*)?`)
	reqJsonStr = re.ReplaceAllString(reqJsonStr, `$1$2$3`)
	re = regexp.MustCompile(`(?m)("signature":\s*")[^"]*(".*"|\s*)(,\s*)?`)
	reqJsonStr = re.ReplaceAllString(reqJsonStr, `$1$2$3`)
	hash := md5.Sum([]byte(reqJsonStr))
	md5Str := fmt.Sprintf("%x", hash)
	return md5Str
}

// GetReqStr  获取请求参数
func (b SimpleBodyParseImpl) GetReqStr(ctx context.Context, req interface{}, servieName, method string) (
	string, error,
) {
	reqJSON, err := json.Marshal(req)

	if err != nil {
		log.Errorf("json.Marshal err|%v|%s|", req, err.Error())
		return "", err
	}

	reqJSON = []byte(fmt.Sprintf("{ \"request\":%s}", reqJSON))
	reqJsonStr := fmt.Sprintf("%s", req)

	return reqJsonStr, nil
}

// GetRspStr  获取响应数据
func (b SimpleBodyParseImpl) GetRspStr(ctx context.Context, response *http.Response, servieName, method string) (
	string, error,
) {
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	if response.StatusCode != 200 {
		return "", fmt.Errorf("response.statusCode is %d,%s", response.StatusCode, responseBody)
	}
	if strings.Contains(string(responseBody), "<!DOCTYPE html>") {
		return "", fmt.Errorf("coockie失效")
	}

	return string(responseBody), nil
}
