package errcode

import (
	"fmt"
	"strconv"
	"strings"

	"trpc.group/trpc-go/trpc-go/errs"
)

// BizErrType 业务类型，错误码 = 4客户端/5服务端 + 三位业务类型 + 三位编号。
// 分段：1～99 common_server 平台服务，100 起业务。
type BizErrType int

const (
	GATEWAY  BizErrType = 1   // 网关
	AUTH     BizErrType = 2   // 登录
	WXTOKEN  BizErrType = 3   // 微信凭证
	BizStart BizErrType = 100 // 业务起始
	MATHGAME BizErrType = 100 // 数学游戏
)

// GetBizErr 获取错误
func GetBizErr(bizType BizErrType) BizErr {
	bizErr := BizErr{
		BizType:   bizType,
		ErrMsgMap: map[int]string{},
	}

	for key, value := range errMsgMap {
		bizErr.ErrMsgMap[key] = value
	}
	return bizErr
}

// BizErr 配置
type BizErr struct {
	BizType   BizErrType
	ErrMsgMap map[int]string
}

// AddBizErrClient 添加客户端错误
func (be *BizErr) AddBizErrClient(code int, msg string) error {
	return be.addBizErr(4, code, msg)
}

// AddBizErrServer 添加服务端错误
func (be *BizErr) AddBizErrServer(code int, msg string) error {

	return be.addBizErr(5, code, msg)
}

// GetBizErrClient 获取客户端(前端)错误
func (be *BizErr) GetBizErrClient(code int, msg ...string) error {
	return be.getBizErr(4, code, msg...)
}

// GetBizErrServer 获取服务端(后端)错误
func (be *BizErr) GetBizErrServer(code int, msg ...string) error {
	return be.getBizErr(5, code, msg...)
}

// addBizErr 增加错误
func (be *BizErr) addBizErr(startCode int, code int, msg string) error {
	clientCodeStr := fmt.Sprintf("%d%03d%03d", startCode, be.BizType, code)
	clientCode, err := stringToInt(clientCodeStr)
	if err != nil {
		return err
	}
	return be.AddCodeMsg(clientCode, msg)
}

// AddCodeMsg 增加错误码和错误
func (be *BizErr) AddCodeMsg(clientCode int, msg string) error {
	if msg == "" {
		msg = "系统异常"
	}
	be.ErrMsgMap[clientCode] = msg
	return nil
}

// GetCodeErr 获取code对应的错误
func (be *BizErr) GetCodeErr(clientCode int, msg ...string) error {
	errMsg := be.ErrMsgMap[clientCode]
	if errMsg == "" {
		return GetErr(ErrCodeErr)
	}

	if len(msg) > 0 {
		for index, v := range msg {
			errMsg = strings.ReplaceAll(errMsg, fmt.Sprintf("{%d}", index), v)
		}
	}
	return errs.New(clientCode, errMsg)
}

// getBizErr 获取错误
func (be *BizErr) getBizErr(startCode int, code int, msg ...string) error {
	clientCodeStr := fmt.Sprintf("%d%03d%03d", startCode, be.BizType, code)
	clientCode, err := stringToInt(clientCodeStr)
	if err != nil {
		return GetErr(ErrCodeErr)
	}
	return be.GetCodeErr(clientCode, msg...)
}

// stringToInt string转int
func stringToInt(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return i, nil
}
