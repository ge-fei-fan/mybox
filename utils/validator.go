package utils

import (
	"errors"
	"reflect"
)

//判断参数是否为结构体
func verifyStruct(st interface{}) (reflect.Type, *reflect.Value, error) {

	typ := reflect.TypeOf(st)
	val := reflect.ValueOf(st) // 获取reflect.Type类型

	kd := val.Kind()
	if kd != reflect.Struct {
		return nil, nil, errors.New("expect struct")
	}
	return typ, &val, nil
}

func RegisterVerify(st interface{}) error {
	typ, val, err := verifyStruct(st)
	if err != nil {
		return err
	}
	num := typ.NumField()
	for i := 0; i < num; i++ {
		tagVal := typ.Field(i)
		value := val.Field(i)
		if tagVal.Name == "HeaderImg" {
			continue
		}
		if isBlank(value) {
			return errors.New(tagVal.Name + "值不能为空")
		}
	}
	return nil
}

//非空校验
func isBlank(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String, reflect.Slice:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}
