package entity

import "gorm.io/plugin/soft_delete"

// User 用户
type User struct {
	BaseEntity
	Username      string                `gorm:"username" json:"username"`
	Password      string                `gorm:"password" json:"-"`
	Mobile        string                `gorm:"mobile" json:"mobile"`
	DeleteFlag    soft_delete.DeletedAt `gorm:"column:delete_status" json:"-"`
	LastLoginTime *LocalTime            `gorm:"last_login_time" json:"lastLoginTime"`
	Backup        string                `gorm:"backup" json:"backup"`
}

// TableName 表名
func (User) TableName() string {
	return "t_sys_user"
}

func (User) EnableRedis() bool {
	return true
}

func (User) CacheKey() string {
	return "cache::user::"
}
