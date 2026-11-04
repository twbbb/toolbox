// Package monitor 上报
package monitor

import (
	"trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/log"
	"trpc.group/trpc-go/trpc-go/metrics"
)

// costBounds 耗时分桶，单位秒，和 task.getUseTime 一致。
// Prometheus 插件收下 SUM、SET、Histogram。耗时用直方图，均值取 cost 的 _sum/_count。
var costBounds = metrics.NewValueBounds(0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 5, 10, 30)

func appendCost(datas []*metrics.Metrics, reportName string, costSec float64) []*metrics.Metrics {
	datas = append(datas, metrics.NewMetrics("allCost", costSec, metrics.PolicySUM))
	histName := "cost"
	if reportName != "" {
		histName = reportName + "_cost"
	}
	// 分桶要在第一次 Observe 之前登记，sink 创建 Histogram 时按这个名字取。
	metrics.Histogram(histName, costBounds)
	return append(datas, metrics.NewMetrics("cost", costSec, metrics.PolicyHistogram))
}

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

	datas = appendCost(datas, reportName, costTms)

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
	dims = append(dims, &metrics.Dimension{Name: "app", Value: trpc.GlobalConfig().Server.App})
	dims = append(dims, &metrics.Dimension{Name: "server", Value: trpc.GlobalConfig().Server.Server})
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

	datas = appendCost(datas, "task", costTms)

	rec := metrics.NewMultiDimensionMetricsX("task", dims, datas)

	if err := metrics.Report(rec); err != nil {
		log.Errorf(
			"%s|Report Invoke Failed! reportName: %s, err: %+v",
			funcName, rec.Name, err,
		)
	}
}
