package monitor

import (
	"context"
	"testing"

	"trpc.group/trpc-go/trpc-go/metrics"
)

type captureSink struct {
	recs []metrics.Record
}

func (s *captureSink) Name() string { return "capture" }

func (s *captureSink) Report(rec metrics.Record, _ ...metrics.Option) error {
	s.recs = append(s.recs, rec)
	return nil
}

func TestCostUsesHistogram(t *testing.T) {
	sink := &captureSink{}
	metrics.RegisterMetricsSink(sink)

	ReportCache("rid", "cache", "LoadMatchStat", true, 0.024)
	ReportPlatformTask("RefreshMatchStat", true, 0.024, "0 */10 * * * *")
	ReportData(context.Background(), nil, "invoke", false, false, 0.024)

	if len(sink.recs) != 3 {
		t.Fatalf("records = %d, want 3", len(sink.recs))
	}
	for _, rec := range sink.recs {
		var sawCost bool
		for _, m := range rec.GetMetrics() {
			if m.Policy() == metrics.PolicyAVG || m.Policy() == metrics.PolicyMAX {
				t.Errorf("%s/%s policy %d is dropped by prometheus", rec.Name, m.Name(), m.Policy())
			}
			if m.Name() != "cost" {
				continue
			}
			sawCost = true
			if m.Policy() != metrics.PolicyHistogram {
				t.Errorf("%s cost policy = %d, want histogram", rec.Name, m.Policy())
			}
		}
		if !sawCost {
			t.Errorf("%s missing cost", rec.Name)
		}
		if _, ok := metrics.GetHistogram(rec.Name + "_cost"); !ok {
			t.Errorf("histogram %s_cost not registered", rec.Name)
		}
	}
}
