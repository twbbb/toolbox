package arrayutil

import (
	"fmt"
	"reflect"
	"strings"
)

// FilterEmptyStrings 过滤掉数组中空字符串
func FilterEmptyStrings(arr []string) []string {
	var result = []string{}
	for _, s := range arr {
		if s != "" {
			result = append(result, s)
		}
	}
	return result
}

// DeDuplication 去重
func DeDuplication(arr []string) []string {
	if arr == nil {
		return nil
	}
	obj := make(map[string]string)
	result := []string{}
	for _, item := range arr {
		if _, ok := obj[item]; ok {
			continue
		}
		obj[item] = item
		result = append(result, item)
	}
	return result
}

func Contains[T ~string | ~int | ~int32](slice []T, value T) bool {
	if slice == nil {
		return false
	}
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// GetStructAttrArr 将结构体数组，转换为只获取一个属性的string数组,fieldName可以为嵌套格式"BasicInfo.OrgId"
func GetStructAttrArr(items interface{}, fieldName string) ([]string, error) {
	v := reflect.ValueOf(items)
	ids := make([]string, v.Len())
	fieldNames := strings.Split(fieldName, ".")

	for i := 0; i < v.Len(); i++ {
		item := v.Index(i).Interface()

		// 如果是指针获取对应的值
		reqVal := reflect.ValueOf(item)
		if reqVal.Kind() == reflect.Ptr {
			if reqVal.IsNil() {
				// 处理 nil 指针的情况
				return nil, fmt.Errorf("nil pointer for item %v", v)
			}
			reqVal = reqVal.Elem()
		}

		// 循环获取每个嵌套结构体的值
		value := reqVal
		for _, fn := range fieldNames {
			value = value.FieldByName(fn)

			if !value.IsValid() {
				return nil, fmt.Errorf("field not found: %s", fieldName)
			}
			if value.Kind() == reflect.Ptr {
				if value.IsNil() {
					// 处理 nil 指针的情况
					return nil, fmt.Errorf("nil pointer for item, %s %v", fieldName, v)
				}
				value = value.Elem()
			}

		}

		ids[i] = value.String()
	}

	return ids, nil
}

// ArrToMapSlice 数组转换为map，key为string，值为结构体的切片
func ArrToMapSlice[T any](arr []T, keyName string) (map[string][]T, error) {
	m := map[string][]T{}
	for _, value := range arr {
		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				// 处理 nil 指针的情况
				return nil, fmt.Errorf("nil pointer for item %+v", value)
			}
			v = v.Elem()
		}

		field := v.FieldByName(keyName)
		if !field.IsValid() {
			return nil, fmt.Errorf("field not found: %s", keyName)
		}
		key := fmt.Sprintf("%v", field.Interface())
		AddToMapSlice(m, key, value)
	}
	return m, nil
}

// ArrToMapMapSlice 数组转换为map，key为keyName1，值为map2[keyName2]结构体的切片
func ArrToMapMapSlice[T any](arr []T, keyName1 string, keyName2 string) (map[string]map[string][]T, error) {
	m := map[string]map[string][]T{}
	for _, value := range arr {
		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				// 处理 nil 指针的情况
				return nil, fmt.Errorf("nil pointer for item %+v", value)
			}
			v = v.Elem()
		}

		field1 := v.FieldByName(keyName1)
		if !field1.IsValid() {
			return nil, fmt.Errorf("field not found: %s", keyName1)
		}
		field2 := v.FieldByName(keyName2)
		if !field2.IsValid() {
			return nil, fmt.Errorf("field not found: %s", keyName2)
		}
		key1 := fmt.Sprintf("%v", field1.Interface())
		key2 := fmt.Sprintf("%v", field2.Interface())
		mm, ok := m[key1]
		if !ok {
			mm = map[string][]T{}
			m[key1] = mm
		}
		AddToMapSlice(mm, key2, value)
	}
	return m, nil
}

// ArrToMapStruct 数组转换为map，key为string，值为结构体
func ArrToMapStruct[T any](arr []T, keyName string) (map[string]T, error) {
	m := map[string]T{}
	for _, value := range arr {
		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				// 处理 nil 指针的情况
				return nil, fmt.Errorf("nil pointer for item %+v", value)
			}
			v = v.Elem()
		}

		field := v.FieldByName(keyName)
		if !field.IsValid() {
			return nil, fmt.Errorf("field not found: %s", keyName)
		}
		key := fmt.Sprintf("%v", field.Interface())
		m[key] = value
	}
	return m, nil
}

// AddToMapSlice map的key为切片，如果没有就创建切片，否则直接放入
func AddToMapSlice[T any](m map[string][]T, key string, value T) {
	_, ok := m[key]
	if !ok {
		m[key] = []T{}
	}
	m[key] = append(m[key], value)
}

// AddToMapSliceDedup map的key为切片，如果没有就创建切片，否则直接放入，并且去重
func AddToMapSliceDedup(m map[string][]string, key string, value string) {

	_, ok := m[key]
	if !ok {
		m[key] = []string{}
	}
	if Contains(m[key], value) {
		return
	}
	m[key] = append(m[key], value)
}

// GetMax 获取最大的数字
func GetMax(nums ...int) int {
	if len(nums) == 0 {
		return 0
	}
	max := nums[0]
	for i := 1; i < len(nums); i++ {
		if nums[i] > max {
			max = nums[i]
		}
	}
	return max
}

// Intersect 两个数组取交集
func Intersect(a, b []string) []string {
	if a == nil || len(a) == 0 {
		return b
	}

	if b == nil || len(b) == 0 {
		return a
	}
	m := make(map[string]bool)
	for _, s := range a {
		m[s] = true
	}
	var result []string
	for _, s := range b {
		if m[s] {
			result = append(result, s)
		}
	}
	return result
}

// TrimSplitFilterEmpty 去空格分割并且去掉空的数据
func TrimSplitFilterEmpty(str, sep string) []string {
	str = strings.TrimSpace(str)
	if str == "" {
		return []string{}
	}
	return FilterEmptyStrings(strings.Split(str, sep))
}

// SplitFilterEmpty 分割并且去掉空的数据
func SplitFilterEmpty(str, sep string) []string {

	return FilterEmptyStrings(strings.Split(str, sep))
}
