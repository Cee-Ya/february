package entity

import (
	"database/sql/driver"
	"february/common/consts"
	"february/pkg/tls"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/plugin/optimisticlock"
	"time"
)

type LocalTime time.Time

// MarshalJSON override time.Time's MarshalJSON method
func (t *LocalTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(*t)
	return []byte(fmt.Sprintf("\"%v\"", tTime.Format(consts.DateFormatYmdhms))), nil
}

func (t LocalTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	tlt := time.Time(t)
	//判断给定时间是否和默认零时间的时间戳相同
	if tlt.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return tlt, nil
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
	Host             string
	Port             int
	ShutdownTimeout  int
	MaxContentLength int64
	ReadTimeout      int
	WriteTimeout     int
	IdleTimeout      int
}

// DB 数据库配置
type DB struct {
	Dsn          string // 连接信息
	Debug        bool   // 是否开启调试模式
	MaxLifetime  int    // 最大连接周期，超过时间的连接就close
	MaxOpenConns int    // 设置最大连接数
	MaxIdleConns int    // 设置闲置连接数
}

// RedisConfig Redis配置
type RedisConfig struct {
	Addr     string
	Username string
	Password string
	DB       int
	UseTLS   bool
	tls.ClientConfig
	RedisType        string
	MasterName       string
	SentinelUsername string
	SentinelPassword string
}

// BaseEntity 基础业务实体
type BaseEntity struct {
	ID         uint64                 `gorm:"id"`
	CreateTime *LocalTime             `gorm:"create_time"`
	ModifyTime *LocalTime             `gorm:"modify_time"`
	Version    optimisticlock.Version `gorm:"version" json:"-"`
}

func (b *BaseEntity) BeforeSave(tx *gorm.DB) error {
	localTime := LocalTime(time.Now())
	b.CreateTime = &localTime
	b.ModifyTime = &localTime
	return nil
}

func (b *BaseEntity) BeforeUpdate(tx *gorm.DB) error {
	localTime := LocalTime(time.Now())
	b.ModifyTime = &localTime
	return nil
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
