// Package logcontext logcontext
package logcontext

import (
	"context"
	"strconv"

	"github.com/twbbb/toolbox/trpcutil/trpcmeta"
	"trpc.group/trpc-go/trpc-go/filter"
	"trpc.group/trpc-go/trpc-go/log"
)

// logContext 被调拦截器的实现
func logContext() filter.ServerFilter {
	return func(ctx context.Context, req interface{}, handler filter.ServerHandleFunc) (rsp interface{}, err error) {
		if requestID := trpcmeta.GetRequestId(ctx); requestID != "" {
			ctx = log.WithContextFields(ctx, "requestId", requestID)
		}
		if userkey := trpcmeta.GetUserkey(ctx); userkey != "" {
			ctx = log.WithContextFields(ctx, "userkey", userkey)
		}
		if cid := trpcmeta.GetCid(ctx); cid != 0 {
			ctx = log.WithContextFields(ctx, "cid", strconv.FormatUint(cid, 10))
		}

		return handler(ctx, req)
	}
}

func init() {
	filter.Register("logContext", logContext(), filter.NoopClientFilter)
}
