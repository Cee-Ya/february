package utils

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
