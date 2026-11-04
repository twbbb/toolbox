// Package logcontext logcontext
package logcontext

import (
	"context"

	"github.com/twbbb/toolbox/trpcutil/trpcmeta"
	"trpc.group/trpc-go/trpc-go/filter"
	"trpc.group/trpc-go/trpc-go/log"
)

// ServerFilter 被调拦截器的实现
func logContext() filter.ServerFilter {
	return func(ctx context.Context, req interface{}, handler filter.ServerHandleFunc) (rsp interface{}, err error) {
		requestId := trpcmeta.GetRequestId(ctx)
		if requestId != "" {
			log.WithContextFields(ctx, "requestId", requestId)
		}

		bizCheckInf, ok := req.(BizIdInf)
		if ok {
			log.WithContextFields(ctx, "bizId", bizCheckInf.GetBizId())
		}

		return handler(ctx, req)
	}
}

// BizIdInf 定义一个接口
type BizIdInf interface {
	GetBizId() string
}

// PatientIdInf 定义一个接口
type PatientIdInf interface {
	GetPatientId() string
}

func init() {
	filter.Register("logContext", logContext(), filter.NoopClientFilter)
}
