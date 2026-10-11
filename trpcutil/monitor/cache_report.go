// Package monitor 上报
package monitor

import (
	"trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/log"
	"trpc.group/trpc-go/trpc-go/metrics"
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

func ReportPlatformTask(funcName string, suc bool, costTms float64, cornTime string) {
	dims := make([]*metrics.Dimension, 0, 5)
	dims = append(dims, &metrics.Dimension{Name: "env", Value: trpc.GlobalConfig().Global.EnvName})
	dims = append(dims, &metrics.Dimension{Name: "srcContain", Value: trpc.GlobalConfig().Global.ContainerName})
	dims = append(dims, &metrics.Dimension{Name: "srcIP", Value: trpc.GlobalConfig().Global.LocalIP})
	dims = append(dims, &metrics.Dimension{Name: "funName", Value: funcName})
	dims = append(dims, &metrics.Dimension{Name: "cornTime", Value: cornTime})

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

	rec := metrics.NewMultiDimensionMetricsX("task", dims, datas)

	if err := metrics.Report(rec); err != nil {
		log.Errorf(
			"%s|Report Invoke Failed! reportName: %s, err: %+v",
			funcName, rec.Name, err,
		)
	}
}
