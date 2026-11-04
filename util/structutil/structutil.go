// Package structutil structutil
package structutil

import (
	"reflect"
	"strings"

	"github.com/jinzhu/copier"
	"trpc.group/trpc-go/trpc-go/log"
)

// GetStructArrName 提取类数组名称
func GetStructArrName(slice interface{}) string {
	// 获取输入参数的反射值
	val := reflect.ValueOf(slice)

	// 检查输入是否为切片
	if val.Kind() != reflect.Slice {
		return ""
	}

	// 获取切片长度
	length := val.Len()

	// 创建一个字符串切片用于存储类名
	names := make([]string, length)

	// 遍历切片的每个元素
	for i := 0; i < length; i++ {
		elem := val.Index(i).Interface()
		names[i] = GetStructName(elem)
	}

	// 使用逗号连接类名并返回
	return strings.Join(names, ",")
}

// GetStructName 提取类名称
func GetStructName(i interface{}) string {
	// 通过 reflect.TypeOf 获取接口的动态类型
	t := reflect.TypeOf(i)

	// 检查是否为指针，如果是，获取指针指向的元素类型
	if t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	}
	return t.Name()
}

// CopyFields 复制数据
func CopyFields(src interface{}, dst interface{}) {
	if src == nil {
		return
	}
	srcValue := reflect.ValueOf(src)
	dstValue := reflect.ValueOf(dst)

	if srcValue.Kind() == reflect.Ptr {
		srcValue = srcValue.Elem()
	}

	if dstValue.Kind() == reflect.Ptr {
		dstValue = dstValue.Elem()
	}

	for i := 0; i < srcValue.NumField(); i++ {
		fieldName := srcValue.Type().Field(i).Name
		fieldValue := srcValue.Field(i)
		dstField := dstValue.FieldByName(fieldName)

		if dstField.IsValid() && dstField.CanSet() {
			dstField.Set(fieldValue)
		}
	}
}

// CopyFieldsByParamKey 复制数据
func CopyFieldsByParamKey(src interface{}, dst interface{}, reqParamKey []string) []string {
	errKey := []string{}
	if src == nil {
		return errKey
	}
	srcValue := reflect.ValueOf(src)
	dstValue := reflect.ValueOf(dst)

	if srcValue.Kind() == reflect.Ptr {
		srcValue = srcValue.Elem()
	}

	if dstValue.Kind() == reflect.Ptr {
		dstValue = dstValue.Elem()
	}

	for _, key := range reqParamKey {
		srcField := srcValue.FieldByName(key)
		dstField := dstValue.FieldByName(key)

		if srcField.IsValid() && dstField.IsValid() && dstField.CanSet() {
			// 检查字段类型是否相同
			if srcField.Type() == dstField.Type() {
				dstField.Set(srcField)
			} else {
				errKey = append(errKey, key)
			}
		} else {
			errKey = append(errKey, key)
		}
	}
	return errKey
}

// Copy ...
func Copy(from interface{}, to interface{}) {
	err := copier.CopyWithOption(to, from, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	if err != nil {
		log.Errorf("conv failed, err=%+v", err)
	}
}
