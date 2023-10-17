// Package dateutil dateutil
package dateutil

import (
	"time"
)

// DefaultFormat 格式化
func DefaultFormat(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
