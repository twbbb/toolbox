// Package report report
package report

import (
	"git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.code.oa.com/trpc-go/trpc-go/metrics"
)

// ReportCache ...
func ReportCache(requestId string, reportName string, biz string, suc bool, costTms float64) {

	dims := make([]*metrics.Dimension, 0, 4)
	dims = append(dims, &metrics.Dimension{Name: "env", Value: trpc.GlobalConfig().Global.EnvName})
	dims = append(dims, &metrics.Dimension{Name: "srcContain", Value: trpc.GlobalConfig().Global.ContainerName})
	dims = append(dims, &metrics.Dimension{Name: "srcIP", Value: trpc.GlobalConfig().Global.LocalIP})
	dims = append(dims, &metrics.Dimension{Name: "biz", Value: biz})

	datas := make([]*metrics.Metrics, 0, 5)
	datas = append(datas, metrics.NewMetrics("allCount", 1, metrics.PolicySUM))
	if !suc {
		datas = append(datas, metrics.NewMetrics("failCount", 1, metrics.PolicySUM))
	} else {
		datas = append(datas, metrics.NewMetrics("failCount", 0, metrics.PolicySUM))
	}

	datas = append(datas, metrics.NewMetrics("allCost", costTms, metrics.PolicySUM))
	datas = append(datas, metrics.NewMetrics("avgCost", costTms, metrics.PolicyAVG))
	datas = append(datas, metrics.NewMetrics("maxCost", costTms, metrics.PolicyMAX))

	rec := metrics.NewMultiDimensionMetricsX(reportName, dims, datas)

	if err := metrics.Report(rec); err != nil {
		log.Errorf(
			"%s|%s|Report Invoke Failed! reportName: %s, err: %+v",
			requestId, biz, rec.Name, err,
		)
	}
}

func calcDistr(distr []int64, distrLen int, val int64) []int {

	result := make([]int, distrLen)

	idx := 0
	for ; idx < distrLen-1; idx++ {
		if len(distr) > idx && val <= distr[idx] {
			break
		}
	}
	result[idx] = 1

	return result
}
