// Package filters filters
package filters

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/twbbb/toolbox/filter/httpproxy/bodyparse"
	"trpc.group/trpc-go/trpc-go/log"

	"trpc.group/trpc-go/trpc-go/codec"
	"trpc.group/trpc-go/trpc-go/filter"
)

// HttpProxyArg 代理的参数
type HttpProxyArg struct {
	// 本地路径，请求参数和结果会存储在这里，下次相同参数可以直接用本地文件,如："/Users/tianwenbin/Documents/workspace"
	LocalFilePath string `yaml:"localFilePath"`
	BaseTrpcUrl   string `yaml:"baseTrpcUrl"` // trpc内网请求地址 如：https://proxy.th.woa.com/pretrpcapi
	Openid        string `yaml:"openid"`      // 内网鉴权openid 如：ooJQ34_lrbyQ-zo92AoqoPvkoXXX
	// trpc网关service名字和北极星名字映射关系，如："trpc.thuser.crowdruleserver.THCrowdRuleServer": "crowd_rule_server"
	ServiceMap map[string]string `yaml:"serviceMap"`
	Cookie     string            `yaml:"cookie"` // woa的cookie
	BodyParse  bodyparse.BodyParse
	HttpForce  bool `yaml:"httpForce"` // 是否强制走http请求
}

var httpProxyArg HttpProxyArg

// init 初始化
func init() {
	runInLocal := os.Getenv("DOCKER_ENV") == ""
	if !runInLocal {
		return
	}

	filter.Register("clientProxy", filter.NoopServerFilter, clientProxy())

}

// InitArg 初始化
func InitArg(arg HttpProxyArg) {
	if arg.BodyParse == nil {
		arg.BodyParse = bodyparse.GetSimpleBodyParseImpl()
	}
	httpProxyArg = arg

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

		requestFlag, httpRsp, err := httpRequest(ctx, servieName, method, req)
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

func httpRequest(ctx context.Context, servieName string, method string, req interface{}) (
	bool, []byte,
	error,
) {
	_, ok := httpProxyArg.ServiceMap[servieName]
	if !ok {
		lowerServiceName := strings.ToLower(servieName)
		if !strings.Contains(lowerServiceName, "mysql") &&
			!strings.Contains(lowerServiceName, "redis") &&
			!strings.Contains(lowerServiceName, "mongo") &&
			!strings.Contains(lowerServiceName, "kafka") &&
			!strings.Contains(lowerServiceName, "tdmq") {
			log.InfoContextf(ctx, "serviceMap not found|%s|%s", servieName, method)
		}
		return false, nil, nil
	}

	reqJsonStr, err := httpProxyArg.BodyParse.GetReqStr(ctx, req, servieName, method)
	if err != nil {
		return false, nil, err
	}
	md5Str := httpProxyArg.BodyParse.GetReqKey(ctx, req, reqJsonStr)
	fileReqName := servieName + method + "Req" + md5Str
	fileRspName := servieName + method + "Rsp" + md5Str
	exists, err := exists(httpProxyArg.LocalFilePath + fileRspName)

	requestFlag := true
	url := fmt.Sprintf("%s/%s/%s", httpProxyArg.BaseTrpcUrl, httpProxyArg.ServiceMap[servieName], method)
	if err != nil {
		return requestFlag, nil, err
	}
	if exists && !httpProxyArg.HttpForce {
		data, err := ioutil.ReadFile(httpProxyArg.LocalFilePath + fileRspName)
		return requestFlag, data, err
	} else {
		rsp, err := sendrequest(ctx, url, reqJsonStr, servieName, method)
		if err != nil {
			return requestFlag, nil, err
		}
		err = ioutil.WriteFile(httpProxyArg.LocalFilePath+fileReqName, []byte(reqJsonStr), 0644)
		if err != nil {
			return requestFlag, nil, err
		}
		err = ioutil.WriteFile(httpProxyArg.LocalFilePath+fileRspName, []byte(rsp), 0644)
		if err != nil {
			return requestFlag, nil, err
		}
		return requestFlag, []byte(rsp), err
	}
}

func sendrequest(ctx context.Context, url string, req, servieName, method string) (string, error) {

	reqBody := strings.NewReader(req)

	request, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("businessid", "tencent-health-mini-qq")
	request.Header.Set("Th-Auth-Type", "tencent-health-mini")
	request.Header.Set("Th-Appid", "wxb032bc789053daf4")
	request.Header.Set("Th-Auto-Openid", httpProxyArg.Openid)
	request.Header.Set("Referer", "http://hlego.woa.com/app/customer_support_server")
	request.Header.Set("User-Agent", "PostmanRuntime/7.31.0")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("Accept", "*/*")
	request.Header.Set("Cookie", httpProxyArg.Cookie)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	return httpProxyArg.BodyParse.GetRspStr(ctx, response, servieName, method)
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
