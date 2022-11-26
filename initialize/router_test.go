package initialize

import (
	"fmt"
	"reflect"
	"testing"
)

type Person struct {
	Name string `mapstructure:"username"`
	Age  int
	Job  string
}

type Cat struct {
	Name  string
	Age   int
	Breed string
}

func TestRouter(t *testing.T) {
	p := Person{
		Name: "测试",
		Age:  11,
		Job:  "学生",
	}

	typ := reflect.TypeOf(p)
	val := reflect.ValueOf(p)
	//kd := val.Kind()
	//num := val.NumField()
	fmt.Printf("type:%v kind:%v\n", typ.Name(), typ.Kind())
	fmt.Println(val.Kind())
	fmt.Printf("%v", val.Field(1))
	//fmt.Println(typ)
	//fmt.Println(val)
	//fmt.Println(kd)
	//fmt.Println(num)

}
