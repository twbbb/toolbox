// Package stringutil stringutil
package stringutil

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// TruncateString trims the input string to the specified length.
func TruncateString(input string, maxLength int) string {
	if len(input) > maxLength {
		return input[:maxLength]
	}
	return input
}

// MD5Hash 获取md5的数据
func MD5Hash(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	hashValue := hash.Sum(nil)
	hashString := fmt.Sprintf("%x", hashValue)
	return hashString
}

// GetFloat64Val 分割并且去掉空的数据
func GetFloat64Val(str string, defaultVal float64) float64 {

	if val, err := strconv.ParseFloat(str, 64); err == nil {
		return val
	}
	return defaultVal
}

// GetUrlParams  urgl字符串解析出参数
func GetUrlParams(urlStr string) (map[string]string, error) {
	if urlStr == "" {
		return nil, errors.New("urlStr is empty")
	}
	url, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	strParams := url.Query().Encode()
	if len(strParams) == 0 {
		return nil, errors.New("no query parameters found")
	}
	params := strings.Split(strParams, "&")
	kvs := make(map[string]string)
	for _, p := range params {
		kv := strings.Split(p, "=")
		if len(kv) != 2 {
			continue
		}
		kvs[kv[0]] = kv[1]
	}
	if len(kvs) == 0 {
		return nil, errors.New("no key-value pairs found in query parameters")
	}
	return kvs, nil
}

// StringToInt string转int
func StringToInt(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return i, nil
}

// IntToString int转string
func IntToString(n interface{}, defaultValue string) string {
	switch n := n.(type) {
	case int:
		return strconv.Itoa(n)
	case int8:
		return strconv.Itoa(int(n))
	case int16:
		return strconv.Itoa(int(n))
	case int32:
		return strconv.Itoa(int(n))
	case int64:
		return strconv.FormatInt(n, 10)
	default:
		return defaultValue
	}
}

// GetPhoneMaskStr 获取脱敏数据
func GetPhoneMaskStr(phoneNumber string) string {

	runes := []rune(phoneNumber)
	runesLen := len(runes)
	if runesLen <= 1 {
		return phoneNumber
	}
	// 11位长度手机号
	if runesLen == 11 {
		masked := []rune(phoneNumber)
		for i := 3; i < 7; i++ {
			masked[i] = '*'
		}
		return string(masked)
	}

	// 截取一半显示
	showLen := runesLen * 1 / 2
	showRune := runes[0:showLen]
	showStr := string(showRune)
	maskStr := showStr + strings.Repeat("*", runesLen-showLen)
	return maskStr
}

// CheckMobile 检验手机号
func CheckMobile(phone string) bool {
	regRuler := `^1[3456789]\d{9}$`
	// 正则调用规则
	reg := regexp.MustCompile(regRuler)
	return reg.MatchString(phone)

}

// IsTimeStringYMDHMS 检查字符串时间格式
func IsTimeStringYMDHMS(str string) bool {
	return IsTimeStringValid(str, "2006-01-02 15:04:05")
}

// IsTimeStringYMD 检查字符串时间格式
func IsTimeStringYMD(str string) bool {
	return IsTimeStringValid(str, "2006-01-02")
}

// IsTimeStringValid 检查字符串时间格式
func IsTimeStringValid(timeString, format string) bool {
	_, err := time.Parse(format, timeString)
	if err != nil {
		return false
	}
	return true
}

// IsNil 判断是否为空
func IsNil(i interface{}) bool {
	if i == nil {
		return false
	}
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	return false
}

// GetJsonStrNilEmpty 获取json字符串
func GetJsonStrNilEmpty(obj interface{}) string {
	if IsNil(obj) {
		return ""
	}
	c, _ := json.Marshal(obj)
	return string(c)
}

// GetJsonStr 获取json字符串
func GetJsonStr(obj interface{}) string {
	c, _ := json.Marshal(obj)
	return string(c)
}

// GetJsonStrNoEscape 获取不转义的json字符串
func GetJsonStrNoEscape(obj interface{}) string {
	buffer := &strings.Builder{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.Encode(obj)
	result := buffer.String()
	// 去掉末尾的换行符
	return result
}

// GetMaskStr 获取脱敏数据
func GetMaskStr(str string, plainType int) string {

	runes := []rune(str)
	runesLen := len(runes)
	if runesLen <= 1 {
		return str
	}
	if plainType == 100 && strings.Contains(str, "@") {
		result := ""
		res := strings.Split(str, "@")
		if len(res) == 2 {
			if len(res[0]) < 3 {
				resString := "***"
				result = resString + "@" + res[1]
			} else {
				res2 := Substr2(str, 0, 3)
				resString := res2 + "***"
				result = resString + "@" + res[1]
			}
			return result
		}
	}
	// 截取一半显示
	showLen := runesLen * 1 / 2
	showRune := runes[0:showLen]
	showStr := string(showRune)
	maskStr := showStr + strings.Repeat("*", runesLen-showLen)
	return maskStr
}

// Substr2 截取字符串
func Substr2(str string, start int, end int) string {
	rs := []rune(str)
	return string(rs[start:end])
}
