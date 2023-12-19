// Package dateutil dateutil
package dateutil

import (
	"sync"
	"time"
)

// DefaultFormat 格式化
func DefaultFormat(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

var (
	beijingLocation *time.Location
	once            sync.Once
)

// GetBeijingLocation 获取北京时区
func GetBeijingLocation() *time.Location {
	once.Do(
		func() {
			var err error
			// 加载"Asia/Shanghai"时区
			beijingLocation, err = time.LoadLocation("Asia/Shanghai")
			if err != nil {
				// 处理错误
				panic(err)
			}
		},
	)
	return beijingLocation
}

// BeijingFormat 格式化为北京时间
func BeijingFormat(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	// 将时间转换为北京时间
	beijingTime := t.In(GetBeijingLocation())
	return beijingTime.Format("2006-01-02 15:04:05")
}

// BeijingByFormat 格式化为北京时间
func BeijingByFormat(t time.Time, format string) string {
	if t.IsZero() {
		return ""
	}
	// 将时间转换为北京时间
	beijingTime := t.In(GetBeijingLocation())
	return beijingTime.Format(format)
}
