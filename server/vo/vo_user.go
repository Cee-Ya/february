package vo

// UserAddVo 用户添加vo
type UserAddVo struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Mobile   string `json:"mobile" binding:"required"`
	Backup   string `json:"backup"`
}

// UserUpdateVo 用户更新vo
type UserUpdateVo struct {
	Id       uint64 `json:"id" binding:"required"`
	Username string `json:"username"`
	Password string `json:"password"`
	Mobile   string `json:"mobile"`
	Backup   string `json:"backup"`
}

// UserPageVo 用户分页vo
type UserPageVo struct {
	Id         uint64 `json:"id"`
	Username   string `json:"username"`
	Mobile     string `json:"mobile"`
	Backup     string `json:"backup"`
	CreateTime string `json:"createTime"`
	ModifyTime string `json:"modifyTime"`
}
