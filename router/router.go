package router

import "github.com/gin-gonic/gin"

func UserRouter(group *gin.RouterGroup) {
	group.GET("/list", GetUserList)
	group.GET("/page", GetPageList)
}
