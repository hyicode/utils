package invoke

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func InvokeByJson(obj interface{}, methodName string, jsonData []byte) ([]reflect.Value, error) {
	return invokeByJson(reflect.ValueOf(obj), methodName, jsonData)
}

// 不支持调用参数中含有 interface/循环引用 结构体的 Method
func invokeByJson(obj reflect.Value, methodName string, jsonData []byte) ([]reflect.Value, error) {
	// 如果传入的是指针，获取其底层值
	if obj.Kind() == reflect.Ptr {
		obj = obj.Elem()
	}

	method := obj.MethodByName(methodName)
	if !method.IsValid() {
		// 如果方法在值接收者上找不到，尝试在指针接收者上找
		if obj.CanAddr() {
			method = obj.Addr().MethodByName(methodName)
			if !method.IsValid() {
				return nil, fmt.Errorf("method %s not found", methodName)
			}
		} else {
			return nil, fmt.Errorf("method %s not found and cannot get address", methodName)
		}
	}

	t := method.Type()
	argsNum := t.NumIn()
	args := make([]reflect.Value, 0, argsNum)

	var jsonArgs []json.RawMessage
	if err := json.Unmarshal(jsonData, &jsonArgs); err != nil {
		var raw json.RawMessage
		if err := json.Unmarshal(jsonData, &raw); err != nil {
			return nil, err
		}
		jsonArgs = []json.RawMessage{raw}
	}

	if t.IsVariadic() {
		if len(jsonArgs) < argsNum-1 {
			return nil, fmt.Errorf("variadic method requires at least %d arguments, but got %d", argsNum-1, len(jsonArgs))
		}
	} else if len(jsonArgs) != argsNum {
		return nil, fmt.Errorf("method requires %d arguments, but got %d", argsNum, len(jsonArgs))
	}

	for i := 0; i < argsNum; i++ {
		if i >= len(jsonArgs) {
			break
		}
		argType := t.In(i)
		// 如果是可变长参数，则将jsonArgs中的数据转换为切片
		if i == argsNum-1 && t.IsVariadic() {
			sliceType := argType.Elem()
			for j := i; j < len(jsonArgs); j++ {
				e := reflect.New(sliceType).Elem()
				if err := json.Unmarshal(jsonArgs[j], e.Addr().Interface()); err != nil {
					return nil, err
				}
				args = append(args, e)
			}
			break
		}
		argValue := reflect.New(argType).Elem()
		if err := json.Unmarshal(jsonArgs[i], argValue.Addr().Interface()); err != nil {
			return nil, err
		}
		args = append(args, argValue)
	}

	return method.Call(args), nil
}
