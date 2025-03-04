// Package dateutil dateutil
package dateutil

import (
	"testing"
	"time"
)

// Test_DefaultFormat 检查字符串时间格式
func Test_DefaultFormat(t *testing.T) {
	t.Log(time.Now())
	aa := DefaultFormat(time.Now())
	date, _ := time.Parse("2006-01-02", "2025-01-02")
	t.Log(aa)
	t.Log(date)
}
