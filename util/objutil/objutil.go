package objectutil

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"trpc.group/trpc-go/trpc-go"
)

// ExtractMethodName 提取最后一个 '/' 后的元素
func ExtractMethodName(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path is empty")
	}
	// Split by '/' and get the last non-empty part
	parts := strings.Split(strings.TrimRight(path, "/"), "/")
	if len(parts) == 0 {
		return "", fmt.Errorf("invalid path")
	}
	return parts[len(parts)-1], nil
}

// GetAndSetFirstReturnType 根据方法名获取第一个返回体并实例化
//
//	GetAndSetFirstReturnType(ctx, (*pb.ThMedicalCenterGoService)(nil), &pb.Result{
//			Msg: "suc",
//		}, "Result")
func GetAndSetFirstReturnType(ctx context.Context, tempService interface{}, addElem interface{}, addElemName string) (interface{}, error) {
	methodName, _ := ExtractMethodName(trpc.Message(ctx).CalleeMethod())
	//methodName := "Ping"
	// 获取接口的反射类型
	ifaceType := reflect.TypeOf(tempService).Elem()

	// 查找方法
	method, exists := ifaceType.MethodByName(methodName)
	if !exists {
		return nil, fmt.Errorf("method %s not found in ThMedicalCenterGoService", methodName)
	}

	// 获取方法返回类型（第一个返回参数）
	if method.Type.NumOut() < 1 {
		return nil, fmt.Errorf("method %s has no return values", methodName)
	}

	// 获取第一个返回类型的具体类型
	firstReturnType := method.Type.Out(0)

	// 检查是否是指针类型
	if firstReturnType.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("first return type of %s is not a pointer", methodName)
	}
	// 创建实例
	instance := reflect.New(firstReturnType.Elem()).Interface()

	instanceValue := reflect.ValueOf(instance)
	if instanceValue.Kind() != reflect.Ptr || instanceValue.IsNil() {
		return nil, fmt.Errorf("instance is not a valid pointer")
	}

	structValue := instanceValue.Elem()
	if structValue.Kind() != reflect.Struct {
		return nil, fmt.Errorf("instance is not a struct")
	}
	resultField := structValue.FieldByName(addElemName)
	if !resultField.IsValid() {
		return nil, fmt.Errorf("Result field not found in %s", firstReturnType.Name())
	}

	if resultField.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("result field is not a pointer")
	}
	if resultField.IsNil() {
		resultType := resultField.Type().Elem()
		newResult := reflect.New(resultType)
		resultField.Set(newResult)
	}
	resultField.Set(reflect.ValueOf(addElem))
	return instance, nil
}
