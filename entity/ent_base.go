package entity

import (
	"ai-report/common/consts"
	"fmt"
	"time"
)

type LocalTime time.Time

// MarshalJSON override time.Time's MarshalJSON method
func (t *LocalTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(*t)
	return []byte(fmt.Sprintf("\"%v\"", tTime.Format(consts.DateFormatYmdhms))), nil
}

// Config 基础配置
type Config struct {
	Server Server
	Log    LogConfig
	DB     DB
	Redis  RedisConfig
}

// LogConfig 日志配置
type LogConfig struct {
	Director      string
	StacktraceKey string
	EncodeLevel   string
	Format        string
	LogInConsole  bool
}

// Server 服务配置
type Server struct {
	Port int
}

// DB 数据库配置
type DB struct {
	Dsn   string // 连接信息
	Debug bool   // 是否开启调试模式
}

// RedisConfig Redis配置
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// BaseEntity 基础业务实体
type BaseEntity struct {
	ID         uint64    `gorm:"id"`
	CreateTime LocalTime `gorm:"create_time"`
	UpdateTime LocalTime `gorm:"update_time"`
	Version    int       `gorm:"version" json:"-"`
}

// Page 分页查询
type Page struct {
	PageNo   int `form:"pageNo"`
	PageSize int `form:"pageSize"`
}

// PageResult 分页查询结果
type PageResult[T any] struct {
	Row   []T   `json:"row"`
	Total int64 `json:"total"`
}

// Trace 定义trace结构体
type Trace struct {
	TraceId   string  `json:"trace_id"`
	SrcMethod *string `json:"srcMethod,omitempty"`
	UserId    int     `json:"user_id"`
}
