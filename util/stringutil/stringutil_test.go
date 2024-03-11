// Package stringutil stringutil
package stringutil

import (
	"testing"
)

// Test_IsTimeStringYMDHMS 检查字符串时间格式
func Test_IsTimeStringYMDHMS(t *testing.T) {
	aa := IsTimeStringYMDHMS("2022-10-10 10:10:10")
	t.Log(aa)
}

// Test_MD5Hash 检查字符串时间格式
func Test_MD5Hash(t *testing.T) {
	aa := MD5Hash("123")
	t.Log(aa)
}
