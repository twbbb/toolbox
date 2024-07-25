// Package report 上报
package report

import (
	"context"

	"git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.code.oa.com/trpc-go/trpc-go/metrics"
	"git.woa.com/th-go/toolbox/util/stringutil"
)

// ReportData 方法上报
func ReportData(ctx context.Context, dims []*metrics.Dimension, name string, isErr, isBizFail bool, costTms float64) {

	// dims := GetDefaultDimension(ctx)

	datas := make([]*metrics.Metrics, 0, 6)
	datas = append(datas, metrics.NewMetrics("allCount", 1, metrics.PolicySUM))
	if isErr {
		datas = append(datas, metrics.NewMetrics("errCount", 1, metrics.PolicySUM))
		datas = append(datas, metrics.NewMetrics("failCount", 1, metrics.PolicySUM))
	} else if isBizFail {
		datas = append(datas, metrics.NewMetrics("errCount", 0, metrics.PolicySUM))
		datas = append(datas, metrics.NewMetrics("failCount", 1, metrics.PolicySUM))
	} else {
		datas = append(datas, metrics.NewMetrics("errCount", 0, metrics.PolicySUM))
		datas = append(datas, metrics.NewMetrics("failCount", 0, metrics.PolicySUM))
	}

	datas = append(datas, metrics.NewMetrics("allCost", float64(costTms), metrics.PolicySUM))
	datas = append(datas, metrics.NewMetrics("avgCost", float64(costTms), metrics.PolicyAVG))
	datas = append(datas, metrics.NewMetrics("maxCost", float64(costTms), metrics.PolicyMAX))

	rec := metrics.NewMultiDimensionMetrics(dims, datas)
	rec.Name = name

	if err := metrics.Report(rec); err != nil {
		log.ErrorContext(
			ctx,
			"Report Invoke Failed! %s, err: %+v",
			stringutil.GetJsonStr(dims), err,
		)
	}

}

// GetDefaultDimension 获取默认的纬度
func GetDefaultDimension(ctx context.Context) []*metrics.Dimension {
	msg := trpc.Message(ctx)
	callerInterface := msg.CalleeServiceName()
	if callerInterface == "" {
		callerInterface = "none"
	}
	dims := []*metrics.Dimension{}
	dims = append(dims, &metrics.Dimension{Name: "srcContain", Value: trpc.GlobalConfig().Global.ContainerName})
	dims = append(dims, &metrics.Dimension{Name: "srcIP", Value: trpc.GlobalConfig().Global.LocalIP})
	dims = append(dims, &metrics.Dimension{Name: "func", Value: callerInterface})
	return dims
}

// AddDimension 增加纬度
func AddDimension(dims []*metrics.Dimension, name, value string) []*metrics.Dimension {
	dims = append(dims, &metrics.Dimension{Name: name, Value: value})
	return dims
}
