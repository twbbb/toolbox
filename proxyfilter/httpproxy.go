package filters

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"git.code.oa.com/trpc-go/trpc-go/log"

	"git.code.oa.com/trpc-go/trpc-go/codec"
	"git.code.oa.com/trpc-go/trpc-go/filter"
)

// 本地路径，请求参数和结果会存储在这里，下次相同参数可以直接用本地文件
var filePath = "/Users/tianwenbin/Documents/workspacejggl/jgglgo/medical_org_mgr/filters/request/"

// trpc内网请求地址
var baseTrpcUrl = "https://proxy.th.woa.com/pretrpcapi"

// 鉴权openid
var openid = "ooJQ34_lrbyQ-zo92AoqoPvkoTk4"

// var filePath = "/Users/tianwenbin/Documents/workspacejggl/jgglgo/medical_org_mgr/filters/foramlrequest/"
// var baseTrpcUrl = "https://proxy.th.woa.com/formaltrpcapi"

// trpc网关service名字和北极星名字映射关系
var serviceMap = map[string]string{
	"trpc.physical.pkg_server.Api":                            "InnerPkgServer",
	"trpc.dmp.readserver.DmpBaseRead":                         "westianTest",
	"trpc.tencent_health.th_coupon_server.BackTHCouponServer": "THCouponServer",
	"trpc.thuser.crowdruleserver.THCrowdRuleServer":           "crowd_rule_server",
	"trpc.physical.consumer_access.ServiceIntf":               "InnerConsumerAccessServer",
}

var cookie = ""

// Init 初始化
func Init(filePathInit, baseTrpcUrlInit, openidInit, cookieInit string, serviceMapInit map[string]string) {
	filePath = filePathInit
	baseTrpcUrl = baseTrpcUrlInit
	openid = openidInit
	serviceMap = serviceMapInit
	cookie = cookieInit
	filter.Register("clientProxy", filter.NoopClientFilter, clientProxy())

}

func httpRequest(servieName string, method string, req []byte) (bool, []byte, error) {
	reqJsonStr := fmt.Sprintf("%s", req)
	// 定义正则表达式
	re := regexp.MustCompile(`(?m)("nonce":\s*")[^"]*(".*"|\s*)(,\s*)?`)
	reqJsonStr = re.ReplaceAllString(reqJsonStr, `$1$2$3`)
	re = regexp.MustCompile(`(?m)("timestamp":\s*")[^"]*(".*"|\s*)(,\s*)?`)
	reqJsonStr = re.ReplaceAllString(reqJsonStr, `$1$2$3`)
	re = regexp.MustCompile(`(?m)("signature":\s*")[^"]*(".*"|\s*)(,\s*)?`)
	reqJsonStr = re.ReplaceAllString(reqJsonStr, `$1$2$3`)

	hash := md5.Sum([]byte(reqJsonStr))
	md5Str := fmt.Sprintf("%x", hash)

	fileReqName := servieName + md5Str + "Req"
	fileRspName := servieName + md5Str + "Rsp"
	exists, err := exists(filePath + fileRspName)
	_, ok := serviceMap[servieName]
	if !ok {
		return false, nil, nil
	}
	requestFlag := true
	url := fmt.Sprintf("%s/%s/%s", baseTrpcUrl, serviceMap[servieName], method)
	if err != nil {
		return requestFlag, nil, err
	}
	if exists {
		data, err := ioutil.ReadFile(filePath + fileRspName)
		return requestFlag, data, err
	} else {
		rsp, err := sendrequest(url, req)
		if err != nil {
			return requestFlag, nil, err
		}
		err = ioutil.WriteFile(filePath+fileReqName, req, 0644)
		if err != nil {
			return requestFlag, nil, err
		}
		err = ioutil.WriteFile(filePath+fileRspName, rsp, 0644)
		if err != nil {
			return requestFlag, nil, err
		}
		return requestFlag, rsp, err
	}
}
func sendrequest(url string, req []byte) ([]byte, error) {

	reqBody := strings.NewReader(string(req))

	request, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("businessid", "tencent-health-mini-qq")
	request.Header.Set("Th-Auth-Type", "tencent-health-mini")
	request.Header.Set("Th-Appid", "wxb032bc789053daf4")
	request.Header.Set("Th-Auto-Openid", openid)
	request.Header.Set("Referer", "http://hlego.woa.com/app/customer_support_server")
	request.Header.Set("User-Agent", "PostmanRuntime/7.31.0")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("Accept", "*/*")
	request.Header.Set("Cookie", cookie)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("response.statusCode is %d", response.StatusCode)
	}
	if strings.Contains(string(responseBody), "<!DOCTYPE html>") {
		return nil, fmt.Errorf("coockie失效")
	}

	return responseBody, nil
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// clientProxy 代理请求
func clientProxy() filter.ClientFilter {
	return func(ctx context.Context, req, rsp interface{}, handler filter.ClientHandleFunc) (err error) {

		begin := time.Now()

		msg := codec.Message(ctx)
		servieName := msg.CalleeServiceName()
		method := path.Base(msg.ClientRPCName())
		if method == "" {
			method = msg.CalleeMethod()

		}
		reqJSON, err := json.Marshal(req)
		// log.Infof("clientProxy start,|%s|%s|%s|", servieName, method, string(reqJSON))

		if err != nil {
			log.Errorf("json.Marshal err|%v|%s|", req, err.Error())
		}

		reqJSON = []byte(fmt.Sprintf("{ \"request\":%s}", reqJSON))

		requestFlag, httpRsp, err := httpRequest(servieName, method, reqJSON)
		// 有走http请求
		if requestFlag {
			if err != nil {
				log.Errorf("httpRequest err|%v|%s|", req, err.Error())
				return err
			}
			// log.Infof("clientProxy httpRsp,%s", string(httpRsp))
			err = json.Unmarshal(httpRsp, rsp)
		} else {
			err = handler(ctx, req, rsp) // 需要显式调用下一个filter
		}

		if err != nil {
			log.Errorf("clientProxy result|%v|%v|%v|%v", time.Since(begin), servieName, method, err.Error())
		} else {
			// log.Infof("clientProxy result|%v|%v|%v|%v", time.Since(begin), servieName, method, rsp)

		}
		return err
	}
}
