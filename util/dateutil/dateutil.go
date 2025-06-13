// Package dateutil dateutil
package dateutil

import (
	"sync"
	"time"
)

var (
	beijingLocation *time.Location
	once            sync.Once
)

func init() {
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
}

// GetBeijingLocation 获取北京时区
func GetBeijingLocation() *time.Location {

	return beijingLocation
}

// DefaultFormat 格式化
func DefaultFormat(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// DefaultFormatYMD 格式化
func DefaultFormatYMD(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
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

// BeijingFormatYMD 格式化为北京时间
func BeijingFormatYMD(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	// 将时间转换为北京时间
	beijingTime := t.In(GetBeijingLocation())
	return beijingTime.Format("2006-01-02")
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

// BeijingParseYMD 格式化
func BeijingParseYMD(t string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02", t, beijingLocation)
}

// BeijingDefaultParse 格式化
func BeijingDefaultParse(t string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02 15:04:05", t, beijingLocation)
}

// BeijingParseByFormat 格式化
func BeijingParseByFormat(format string, t string) (time.Time, error) {
	return time.ParseInLocation(format, t, beijingLocation)
}

// GetDayAddDay 获取下一天
func GetDayAddDay(lastJqEndDay string, dayNum int) (
	string, error,
) {
	// 将结束日期加一天
	endDay, err := BeijingParseYMD(lastJqEndDay)
	if err != nil {
		return "", err
	}
	lastJqEndNextDay := BeijingFormatYMD(endDay.AddDate(0, 0, dayNum))
	return lastJqEndNextDay, nil
}
