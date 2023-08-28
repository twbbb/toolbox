package structutil

import "reflect"

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
