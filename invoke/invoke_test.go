package invoke

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"
)

// TestStruct 用于测试的结构体
type TestStruct struct {
	Name string
	Age  int
}

// NoParam 无参数方法
func (t *TestStruct) NoParam() string {
	return "Hello World"
}

// SingleParam 单参数方法
func (t *TestStruct) SingleParam(data map[string]interface{}) string {
	if data == nil {
		return ""
	}
	if msg, ok := data["message"]; ok {
		return msg.(string)
	}
	return ""
}

// MultiParam 多参数方法
func (t *TestStruct) MultiParam(name string, age int) {
	t.Name = name
	t.Age = age
}

// ComplexTypes 测试复杂类型参数
func (t *TestStruct) ComplexTypes(str string, num int, b bool, f float64, m map[string]interface{}, s []interface{}) {
	t.Name = str
	t.Age = num
}

// PointerParam 测试指针类型参数
func (t *TestStruct) PointerParam(name *string, age *int) {
	if name != nil {
		t.Name = *name
	}
	if age != nil {
		t.Age = *age
	}
}

// SliceParam 测试切片类型参数
func (t *TestStruct) SliceParam(names []string) {
	if len(names) > 0 {
		t.Name = names[0]
	}
}

// MapParam 测试map类型参数
func (t *TestStruct) MapParam(data map[string]string) {
	if name, ok := data["name"]; ok {
		t.Name = name
	}
	if _, ok := data["age"]; ok {
		t.Age = 25 // 固定值用于测试
	}
}

// 值接收者方法
func (t TestStruct) ValueReceiverMethod(name string) string {
	return "Hello, " + name
}

// 指针接收者方法
func (t *TestStruct) PointerReceiverMethod(name string) string {
	t.Name = name
	return "Hello, " + name
}

// 结构体参数方法
func (t *TestStruct) StructParamMethod(param TestStruct) string {
	return fmt.Sprintf("Name: %s, Age: %d", param.Name, param.Age)
}

// 结构体指针参数方法
func (t *TestStruct) StructPointerParamMethod(param *TestStruct) string {
	return fmt.Sprintf("Name: %s, Age: %d", param.Name, param.Age)
}

// NoReturn 无返回值方法
func (t *TestStruct) NoReturn(name string) {
	t.Name = name
}

// VariadicParam 可变长参数方法
func (t *TestStruct) VariadicParam(name string, ages ...int) string {
	t.Name = name
	if len(ages) > 0 {
		t.Age = ages[0]
	}
	return fmt.Sprintf("Name: %s, Ages: %v", name, ages)
}

// FixedAndVariadicParam 固定参数和可变长参数组合的方法
func (t *TestStruct) FixedAndVariadicParam(name string, age int, scores ...int) string {
	t.Name = name
	t.Age = age
	if len(scores) > 0 {
		return fmt.Sprintf("Name: %s, Age: %d, Scores: %v", name, age, scores)
	}
	return fmt.Sprintf("Name: %s, Age: %d, Scores: []", name, age)
}

// InterfaceParam 接口参数
type InterfaceParam interface {
	GetName() string
}

// InterfaceImpl 接口实现
type InterfaceImpl struct {
	Name string
}

func (i InterfaceImpl) GetName() string {
	return i.Name
}

// InterfaceParamMethod 使用接口作为参数的方法
func (t *TestStruct) InterfaceParamMethod(param InterfaceParam) string {
	t.Name = param.GetName()
	return fmt.Sprintf("Name: %s", param.GetName())
}

// CircularStructA 循环引用结构体A
type CircularStructA struct {
	Name string
	B    *CircularStructB
}

// CircularStructB 循环引用结构体B
type CircularStructB struct {
	Name string
	A    *CircularStructA
}

// CircularStructParam 使用循环引用结构体作为参数的方法
func (t *TestStruct) CircularStructParam(param CircularStructA) string {
	t.Name = param.Name
	return fmt.Sprintf("Name: %s, B: %+v", param.Name, param.B)
}

func TestInvokeByJson(t *testing.T) {
	tests := []struct {
		name     string
		obj      interface{}
		method   string
		jsonStr  string
		wantErr  bool
		validate func(t *TestStruct, results []reflect.Value) bool
	}{
		{
			name:    "测试无参数方法",
			obj:     &TestStruct{},
			method:  "NoParam",
			jsonStr: `[]`,
			wantErr: false,
			validate: func(t *TestStruct, results []reflect.Value) bool {
				return len(results) == 1 && results[0].Interface() == "Hello World"
			},
		},
		{
			name:    "测试单参数方法",
			obj:     &TestStruct{},
			method:  "SingleParam",
			jsonStr: `[{"message": "测试消息"}]`,
			wantErr: false,
			validate: func(t *TestStruct, results []reflect.Value) bool {
				return len(results) == 1 && results[0].Interface() == "测试消息"
			},
		},
		{
			name:    "测试多参数方法",
			obj:     &TestStruct{},
			method:  "MultiParam",
			jsonStr: `["张三", 25]`,
			wantErr: false,
			validate: func(t *TestStruct, results []reflect.Value) bool {
				return len(results) == 0 && t.Name == "张三" && t.Age == 25
			},
		},
		{
			name:    "测试不存在的方法",
			obj:     &TestStruct{},
			method:  "NotExistMethod",
			jsonStr: "{}",
			wantErr: true,
			validate: func(t *TestStruct, results []reflect.Value) bool {
				return len(results) == 0
			},
		},
		{
			name:    "测试无效的JSON",
			obj:     &TestStruct{},
			method:  "MultiParam",
			jsonStr: `{invalid json}`,
			wantErr: true,
			validate: func(t *TestStruct, results []reflect.Value) bool {
				return len(results) == 0
			},
		},
		{
			name:    "测试复杂类型参数",
			obj:     &TestStruct{},
			method:  "ComplexTypes",
			jsonStr: `["测试", 100, true, 3.14, {"key": "value"}, [1, 2, 3]]`,
			wantErr: false,
			validate: func(t *TestStruct, results []reflect.Value) bool {
				return len(results) == 0 && t.Name == "测试" && t.Age == 100
			},
		},
		{
			name:    "测试指针类型参数",
			obj:     &TestStruct{},
			method:  "PointerParam",
			jsonStr: `["指针测试", 30]`,
			wantErr: false,
			validate: func(t *TestStruct, results []reflect.Value) bool {
				return len(results) == 0 && t.Name == "指针测试" && t.Age == 30
			},
		},
		{
			name:    "测试切片类型参数",
			obj:     &TestStruct{},
			method:  "SliceParam",
			jsonStr: `[["张三", "李四", "王五"]]`,
			wantErr: false,
			validate: func(t *TestStruct, results []reflect.Value) bool {
				return len(results) == 0 && t.Name == "张三"
			},
		},
		{
			name:    "测试map类型参数",
			obj:     &TestStruct{},
			method:  "MapParam",
			jsonStr: `[{"name": "map测试", "age": "25"}]`,
			wantErr: false,
			validate: func(t *TestStruct, results []reflect.Value) bool {
				return len(results) == 0 && t.Name == "map测试" && t.Age == 25
			},
		},
		{
			name:    "测试值接收者方法",
			obj:     TestStruct{Name: "test", Age: 18},
			method:  "ValueReceiverMethod",
			jsonStr: `["world"]`,
			wantErr: false,
			validate: func(t *TestStruct, results []reflect.Value) bool {
				return len(results) == 1 && results[0].Interface() == "Hello, world"
			},
		},
		{
			name:    "测试指针接收者方法",
			obj:     &TestStruct{Name: "test", Age: 18},
			method:  "PointerReceiverMethod",
			jsonStr: `["world"]`,
			wantErr: false,
			validate: func(t *TestStruct, results []reflect.Value) bool {
				return len(results) == 1 && results[0].Interface() == "Hello, world" && t.Name == "world"
			},
		},
		{
			name:    "测试结构体参数",
			obj:     &TestStruct{Name: "test", Age: 18},
			method:  "StructParamMethod",
			jsonStr: `[{"Name": "param", "Age": 20}]`,
			wantErr: false,
			validate: func(t *TestStruct, results []reflect.Value) bool {
				return len(results) == 1 && results[0].Interface() == "Name: param, Age: 20"
			},
		},
		{
			name:    "测试结构体指针参数",
			obj:     &TestStruct{Name: "test", Age: 18},
			method:  "StructPointerParamMethod",
			jsonStr: `[{"Name": "param", "Age": 20}]`,
			wantErr: false,
			validate: func(t *TestStruct, results []reflect.Value) bool {
				return len(results) == 1 && results[0].Interface() == "Name: param, Age: 20"
			},
		},
		{
			name:    "测试无返回值方法",
			obj:     &TestStruct{},
			method:  "NoReturn",
			jsonStr: `["测试无返回值"]`,
			wantErr: false,
			validate: func(t *TestStruct, results []reflect.Value) bool {
				return len(results) == 0 && t.Name == "测试无返回值"
			},
		},
		{
			name:    "测试可变长参数方法-无可变参数",
			obj:     &TestStruct{},
			method:  "VariadicParam",
			jsonStr: `["张三"]`,
			wantErr: false,
			validate: func(t *TestStruct, results []reflect.Value) bool {
				return len(results) == 1 && results[0].Interface() == "Name: 张三, Ages: []" && t.Name == "张三" && t.Age == 0
			},
		},
		{
			name:    "测试可变长参数方法-一个可变参数",
			obj:     &TestStruct{},
			method:  "VariadicParam",
			jsonStr: `["张三", 25]`,
			wantErr: false,
			validate: func(t *TestStruct, results []reflect.Value) bool {
				return len(results) == 1 && results[0].Interface() == "Name: 张三, Ages: [25]" && t.Name == "张三" && t.Age == 25
			},
		},
		{
			name:    "测试可变长参数方法-多个可变参数",
			obj:     &TestStruct{},
			method:  "VariadicParam",
			jsonStr: `["张三", 25, 30, 35]`,
			wantErr: false,
			validate: func(t *TestStruct, results []reflect.Value) bool {
				return len(results) == 1 && results[0].Interface() == "Name: 张三, Ages: [25 30 35]" && t.Name == "张三" && t.Age == 25
			},
		},
		{
			name:    "测试可变长参数方法-参数不足",
			obj:     &TestStruct{},
			method:  "VariadicParam",
			jsonStr: `[]`,
			wantErr: true,
			validate: func(t *TestStruct, results []reflect.Value) bool {
				return len(results) == 0
			},
		},
		{
			name:    "测试固定参数和可变长参数组合-无可变参数",
			obj:     &TestStruct{},
			method:  "FixedAndVariadicParam",
			jsonStr: `["张三", 25]`,
			wantErr: false,
			validate: func(t *TestStruct, results []reflect.Value) bool {
				return len(results) == 1 && results[0].Interface() == "Name: 张三, Age: 25, Scores: []" && t.Name == "张三" && t.Age == 25
			},
		},
		{
			name:    "测试固定参数和可变长参数组合-一个可变参数",
			obj:     &TestStruct{},
			method:  "FixedAndVariadicParam",
			jsonStr: `["张三", 25, 90]`,
			wantErr: false,
			validate: func(t *TestStruct, results []reflect.Value) bool {
				return len(results) == 1 && results[0].Interface() == "Name: 张三, Age: 25, Scores: [90]" && t.Name == "张三" && t.Age == 25
			},
		},
		{
			name:    "测试固定参数和可变长参数组合-多个可变参数",
			obj:     &TestStruct{},
			method:  "FixedAndVariadicParam",
			jsonStr: `["张三", 25, 90, 85, 95]`,
			wantErr: false,
			validate: func(t *TestStruct, results []reflect.Value) bool {
				return len(results) == 1 && results[0].Interface() == "Name: 张三, Age: 25, Scores: [90 85 95]" && t.Name == "张三" && t.Age == 25
			},
		},
		{
			name:    "测试固定参数和可变长参数组合-参数不足",
			obj:     &TestStruct{},
			method:  "FixedAndVariadicParam",
			jsonStr: `["张三"]`,
			wantErr: true,
			validate: func(t *TestStruct, results []reflect.Value) bool {
				return len(results) == 0
			},
		},
		{
			name:    "测试接口参数-基本实现",
			obj:     &TestStruct{},
			method:  "InterfaceParamMethod",
			jsonStr: `[{"Name": "接口实现"}]`,
			wantErr: true,
			validate: func(t *TestStruct, results []reflect.Value) bool {
				return true
			},
		},
		{
			name:    "测试接口参数-空接口",
			obj:     &TestStruct{},
			method:  "InterfaceParamMethod",
			jsonStr: `[{}]`,
			wantErr: true,
			validate: func(t *TestStruct, results []reflect.Value) bool {
				return true
			},
		},
		{
			name:    "测试循环引用结构体-基本字段",
			obj:     &TestStruct{},
			method:  "CircularStructParam",
			jsonStr: `[{"Name": "A", "B": null}]`,
			wantErr: false,
			validate: func(t *TestStruct, results []reflect.Value) bool {
				return len(results) == 1 && results[0].Interface() == "Name: A, B: <nil>" && t.Name == "A"
			},
		},
		{
			name:    "测试循环引用结构体-简单循环",
			obj:     &TestStruct{},
			method:  "CircularStructParam",
			jsonStr: `[{"Name": "A", "B": {"Name": "B", "A": {"Name": "A", "B": null}}}]`,
			wantErr: false,
			validate: func(t *TestStruct, results []reflect.Value) bool {
				got := fmt.Sprintf("%v", results[0].Interface())
				pattern := `Name: A, B: &{Name:B A:0x[0-9a-f]+}`
				matched, _ := regexp.MatchString(pattern, got)
				return len(results) == 1 && matched && t.Name == "A"
			},
		},
		{
			name:    "测试循环引用结构体-多层循环",
			obj:     &TestStruct{},
			method:  "CircularStructParam",
			jsonStr: `[{"Name": "A1", "B": {"Name": "B1", "A": {"Name": "A2", "B": {"Name": "B2", "A": {"Name": "A1", "B": null}}}}}]`,
			wantErr: false,
			validate: func(t *TestStruct, results []reflect.Value) bool {
				got := fmt.Sprintf("%v", results[0].Interface())
				pattern := `Name: A1, B: &{Name:B1 A:0x[0-9a-f]+}`
				matched, _ := regexp.MatchString(pattern, got)
				return len(results) == 1 && matched && t.Name == "A1"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := InvokeByJson(tt.obj, tt.method, []byte(tt.jsonStr))

			// 验证错误
			if (err != nil) != tt.wantErr {
				t.Errorf("InvokeByJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 验证结果
			if tt.validate != nil {
				if obj, ok := tt.obj.(*TestStruct); ok {
					if !tt.validate(obj, results) {
						t.Error("validation failed")
					}
				} else if obj, ok := tt.obj.(TestStruct); ok {
					if !tt.validate(&obj, results) {
						t.Error("validation failed")
					}
				}
			}
		})
	}
}
