// Package errcode errcode
package errcode

import (
	"git.code.oa.com/trpc-go/trpc-go/errs"
)

// 4客户端错误, 5服务器错误,
// 中间三位 000 公共错误码，001，经期
// 后三位，错误码对应编号。

const (
	// SUCCESS 成功
	SUCCESS = 0

	// ErrReqParams 请求参数异常
	ErrReqParams = 4000000

	// ErrEmptySthid sthid为空
	ErrEmptySthid = 4000001

	// ErrDuplicateReq 重复请求
	ErrDuplicateReq = 4000002

	// ErrServer 服务内部错误
	ErrServer = 5000000

	// ErrLogic 逻辑错误
	ErrLogic = 5000001

	// ErrDBExec db exec error
	ErrDBExec = 5000002

	// ErrRPC 第三方调用返回错误
	ErrRPC = 5000003

	// ErrCodeErr 错误码异常
	ErrCodeErr = 5000004
)

var errMsgMap map[int]string

func init() {
	errMsgMap = map[int]string{
		ErrDBExec:       "数据操作异常",
		ErrEmptySthid:   "sthid为空",
		ErrReqParams:    "请求异常",
		ErrServer:       "系统异常",
		ErrRPC:          "远程调用异常",
		ErrLogic:        "逻辑异常",
		SUCCESS:         "成功",
		ErrCodeErr:      "错误异常",
		ErrDuplicateReq: "请求还在处理中",
	}
}

// GetServerError 逻辑有问题
func GetServerError(errmsg string) error {
	defaulMsg := errMsgMap[ErrServer]
	if errmsg != "" {
		defaulMsg = errmsg
	}
	return errs.New(ErrServer, defaulMsg)
}

// GetLogicError 逻辑有问题
func GetLogicError(errmsg string) error {
	defaulMsg := errMsgMap[ErrLogic]
	if errmsg != "" {
		defaulMsg = errmsg
	}
	return errs.New(ErrLogic, defaulMsg)
}

// GetRequestError 获取请求异常error
func GetRequestError(errmsg string) error {
	defaulMsg := errMsgMap[ErrReqParams]
	if errmsg != "" {
		defaulMsg = errmsg
	}
	return errs.New(ErrReqParams, defaulMsg)
}

// GetRpcError 获取请求异常error
func GetRpcError(errcode ...string) error {
	defaulMsg := errMsgMap[ErrRPC]
	if len(errcode) > 0 {
		defaulMsg = errcode[0]
	}
	return errs.New(ErrRPC, defaulMsg)
}

// GetErrMsg 获取错误信息
func GetErrMsg(code int) string {
	msg := errMsgMap[code]
	if msg == "" {
		msg = "系统繁忙,请稍后再试"
	}
	return msg
}

// GetErr 获取错误
func GetErr(code int) error {
	return errs.New(code, GetErrMsg(code))
}
