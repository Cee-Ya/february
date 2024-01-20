package tools

import (
	"encoding/json"
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
