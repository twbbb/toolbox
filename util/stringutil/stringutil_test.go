package stringutil

import (
	"testing"
)

// Test_IsTimeStringYMDHMS 检查字符串时间格式
func Test_IsTimeStringYMDHMS(t *testing.T) {
	aa := IsTimeStringYMDHMS("2022-10-10 10:10:10")
	t.Log(aa)
}
