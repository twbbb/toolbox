// Package comimport 注册公网上有的 tRPC 插件。
package comimport

import (
	_ "trpc.group/trpc-go/trpc-filter/debuglog"
	_ "trpc.group/trpc-go/trpc-filter/recovery"
	_ "trpc.group/trpc-go/trpc-metrics-runtime"
	_ "trpc.group/trpc-go/trpc-naming-polarismesh"
)
