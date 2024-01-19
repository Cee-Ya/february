package entity

// User 用户
type User struct {
	BaseEntity
	Mobile        string    `gorm:"mobile"`
	DeleteFlag    bool      `gorm:"delete_status"`
	LastLoginTime LocalTime `gorm:"last_login_time"`
	Backup        string    `gorm:"backup"`
}

// TableName 表名
func (User) TableName() string {
	return "t_sys_user"
}
