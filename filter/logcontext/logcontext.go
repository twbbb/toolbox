// Package logcontext logcontext
package logcontext

import (
	"context"

	"git.code.oa.com/trpc-go/trpc-go/filter"
	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.woa.com/th-go/toolbox/trpcutil/trpcmeta"
)

// ServerFilter 被调拦截器的实现
func logContext() filter.ServerFilter {
	return func(ctx context.Context, req interface{}, handler filter.ServerHandleFunc) (rsp interface{}, err error) {
		requestId := trpcmeta.GetRequestId(ctx)
		sthid := trpcmeta.GetSthid(ctx)
		if requestId != "" {
			log.WithContextFields(ctx, "requestId", requestId)
		}
		if sthid != "" {
			log.WithContextFields(ctx, "sthid", sthid)
		}

		bizCheckInf, ok := req.(BizIdInf)
		if ok {
			log.WithContextFields(ctx, "bizId", bizCheckInf.GetBizId())
		}

		patientIdInf, ok := req.(PatientIdInf)
		if ok {
			log.WithContextFields(ctx, "patientId", patientIdInf.GetPatientId())
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
