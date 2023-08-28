package timeutil

import (
	"time"
)

// TimeStr ...
func TimeStr() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// CostTms ...
func CostTms(start time.Time) int64 {

	costDuration := time.Since(start)
	return costDuration.Milliseconds()
}
