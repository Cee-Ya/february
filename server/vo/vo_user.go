package vo

// UserAddVo 用户添加vo
type UserAddVo struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Mobile   string `json:"mobile"`
	Backup   string `json:"backup"`
}
