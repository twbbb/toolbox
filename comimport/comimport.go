// Package comimport 公共导入
package comimport

import (
	_ "git.code.oa.com/trpc-go/trpc-config-tconf"
	_ "git.code.oa.com/trpc-go/trpc-filter/debuglog"
	_ "git.code.oa.com/trpc-go/trpc-filter/recovery"
	_ "git.code.oa.com/trpc-go/trpc-metrics-runtime"
	_ "git.code.oa.com/trpc-go/trpc-naming-polaris"
	_ "git.woa.com/opentelemetry/opentelemetry-go-ecosystem/instrumentation/otelcarriers/oteldatabase"
	_ "git.woa.com/opentelemetry/opentelemetry-go-ecosystem/instrumentation/oteltrpc"
	_ "git.woa.com/th-cloud/trpc-config-configmapfile"
	_ "git.woa.com/trpc-go/trpc-metrics-zhiyan/v2"
	_ "git.woa.com/wedev/trpc-route-propagator"
)
