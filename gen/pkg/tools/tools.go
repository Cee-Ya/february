package tools

import (
	"bytes"
	"february/gen/pkg/conf"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log"
	"os"
	"strings"
	"unicode"
)

func PathCreate(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

// CreateCacheName 返回缓存名称
// 例如输入 User 返回cache::user::
// 输入 AiLog 返回cache::ai::log::
func CreateCacheName(clazzName string) string {
	cacheName := "cache"
	for _, v := range clazzName {
		if unicode.IsUpper(v) {
			cacheName += "::"
		}
		cacheName += strings.ToLower(string(v))
	}
	cacheName += "::"
	return cacheName
}

// FormatCamelToSnake 驼峰转蛇形
func FormatCamelToSnake(text string) string {
	var result string
	for i, v := range text {
		if unicode.IsUpper(v) {
			if i != 0 {
				result += "_"
			}
			result += strings.ToLower(string(v))
		} else {
			result += string(v)
		}
	}
	return result
}

// FormatStructName 格式化结构体名称和字段名称  首字母大写
func FormatStructName(prefix, tableName string) string {
	if prefix != "" {
		tableName = strings.TrimPrefix(tableName, prefix)
	}
	caser := cases.Title(language.Und)
	parts := strings.Split(tableName, "_")
	for i := range parts {
		parts[i] = caser.String(strings.ToLower(strings.TrimSpace(parts[i])))
	}
	return strings.Join(parts, "")
}

// FormatJsonColumn 格式化结构体名称和字段名称  首字母小写
func FormatJsonColumn(prefix, tableName string) string {
	//只针对HH，kingbase数据库表未用下划线分割，首字母小写返回
	if conf.C.Database.DBType == "mysql" {
		if prefix != "" {
			tableName = strings.TrimPrefix(tableName, prefix)
		}
		caser := cases.Title(language.Und)
		parts := strings.Split(tableName, "_")
		for i, part := range parts {
			if i == 0 {
				parts[i] = strings.ToLower(strings.TrimSpace(part))
				continue
			}
			parts[i] = caser.String(strings.ToLower(strings.TrimSpace(parts[i])))
		}
		return strings.Join(parts, "")
	}
	if len(tableName) < 1 {
		return tableName
	}
	runes := []rune(tableName)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// PathExist 判断目录是否存在
func PathExist(addr string) bool {
	s, err := os.Stat(addr)
	if err != nil {
		log.Println(err)
		return false
	}
	return s.IsDir()
}

func FileCreate(content bytes.Buffer, name string) {
	file, err := os.Create(name)
	if err != nil {
		log.Println(err)
	}
	_, err = file.WriteString(content.String())
	if err != nil {
		log.Println(err)
	}
	file.Close()
}
