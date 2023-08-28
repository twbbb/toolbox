package stringutil

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
)

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
