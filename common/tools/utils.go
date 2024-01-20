package tools

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
)

// ToJson 将结构体转换为json
func ToJson(v interface{}) (string, error) {
	// 结构体转json
	jsonData, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// Int2String 将int转换为string
func Int2String[T int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64](i T) string {
	return strconv.FormatInt(int64(i), 10)
}

// Str2Uint64 将字符串转换为uint64
func Str2Uint64(str string) (uint64, error) {
	return strconv.ParseUint(str, 10, 64)
}

// ToString 将结构体转换为json
func ToString(v interface{}) (string, error) {
	// 结构体转json
	jsonData, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// Str2Struct 将json字符串转换为结构体
func Str2Struct(str string, v interface{}) error {
	return json.Unmarshal([]byte(str), v)
}

// ContainsString 判断字符串是否在切片中
func ContainsString(slice []string, s string) bool {
	return slice != nil && len(slice) > 0 && IndexString(slice, s) != -1
}

// IndexString 获取字符串在切片中的索引
func IndexString(slice []string, s string) int {
	for i, v := range slice {
		if v == s {
			return i
		}
	}
	return -1
}

// Any2Uint64 将任意类型转换为uint64
func Any2Uint64(v interface{}) (uint64, error) {
	switch v.(type) {
	case int:
		return uint64(v.(int)), nil
	case int8:
		return uint64(v.(int8)), nil
	case int16:
		return uint64(v.(int16)), nil
	case int32:
		return uint64(v.(int32)), nil
	case int64:
		return uint64(v.(int64)), nil
	case uint:
		return uint64(v.(uint)), nil
	case uint8:
		return uint64(v.(uint8)), nil
	case uint16:
		return uint64(v.(uint16)), nil
	case uint32:
		return uint64(v.(uint32)), nil
	case uint64:
		return v.(uint64), nil
	case string:
		return strconv.ParseUint(v.(string), 10, 64)
	default:
		return 0, errors.New("unsupported type")
	}
}

// GetStructField 获取结构体的字段值
func GetStructField(t interface{}, field string) (interface{}, error) {
	// 获取输入参数的类型
	tType := reflect.TypeOf(t)

	// 确保输入参数是结构体类型
	if tType.Kind() != reflect.Struct {
		return nil, errors.New("input parameter is not a struct")
	}

	// 获取结构体的值
	tValue := reflect.ValueOf(t)

	// 遍历结构体的字段
	for i := 0; i < tType.NumField(); i++ {
		// 获取字段的名称
		fieldName := tType.Field(i).Name
		if tValue.Field(i).Kind() == reflect.Struct {
			if fieldValue, err := GetStructField(tValue.Field(i).Interface(), field); err == nil {
				return fieldValue, nil
			}
		}
		// 如果字段名称匹配，则获取字段的值
		if fieldName == field {
			fieldValue := tValue.Field(i).Interface()
			return fieldValue, nil
		}
	}

	// 如果没有找到匹配的字段，则返回错误
	return nil, errors.New("field not found")
}
