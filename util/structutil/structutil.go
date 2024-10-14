// Package structutil structutil
package structutil

import "reflect"

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
