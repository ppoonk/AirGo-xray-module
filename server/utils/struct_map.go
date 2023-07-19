package utils

import (
	"reflect"
	"strings"
)

// struct转map
func StructToMap(data interface{}) map[string]interface{} {

	m := make(map[string]interface{})

	v := reflect.ValueOf(data)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return m
	}
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Name
		tag := t.Field(i).Tag.Get("json")
		if tag == "-" || name == "-" {
			continue
		}
		if tag != "" {
			index := strings.Index(tag, ",")
			if index == -1 {
				name = tag
			} else {
				name = tag[:index]
			}
		}
		m[name] = v.Field(i).Interface()
	}
	return m
}

// 利用反射将结构体转化为map
func StructToMap1(obj interface{}) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	data := make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		if obj1.Field(i).Tag.Get("mapstructure") != "" {
			data[obj1.Field(i).Tag.Get("mapstructure")] = obj2.Field(i).Interface()
		} else {
			data[obj1.Field(i).Name] = obj2.Field(i).Interface()
		}
	}
	return data
}

// 结构体拷贝，结构体成员类型一致
func CopyStruct(fromStruct, toStruct interface{}) {
	fromValue := reflect.ValueOf(fromStruct).Elem()
	toValue := reflect.ValueOf(toStruct).Elem()

	for i := 0; i < fromValue.NumField(); i++ {
		value := fromValue.Field(i)
		name := fromValue.Type().Field(i).Name

		dvalue := toValue.FieldByName(name)
		if dvalue.IsValid() == false { //查询有数据结构体中相同属性的字段，有则修改其值，否则跳过
			continue
		}
		dvalue.Set(value) //这里默认共同成员的类型一样，否则这个地方可能导致 panic，需要简单修改一下。
	}
}
