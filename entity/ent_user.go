package entity

import "gorm.io/plugin/soft_delete"

// User 用户
type User struct {
	BaseEntity
	Mobile        string                `gorm:"mobile"`
	DeleteFlag    soft_delete.DeletedAt `gorm:"column:delete_status" json:"-"`
	LastLoginTime LocalTime             `gorm:"last_login_time"`
	Backup        string                `gorm:"backup"`
}

// TableName 表名
func (User) TableName() string {
	return "t_sys_user"
}
