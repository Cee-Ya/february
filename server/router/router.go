package router

import (
	"ai-report/pkg/ginx/mid"
	"github.com/gin-gonic/gin"
)

// Init router
func Init() *gin.Engine {
	r := gin.New()
	r.NoRoute(mid.NoRoute())
	r.Use(mid.Cors(), mid.GinLogger(), mid.GinRecovery(true))
	UserRouter(r.Group("/user"))
	return r
}

func UserRouter(group *gin.RouterGroup) {
	group.POST("/add", AddUser)
	group.GET("/list", GetUserList)
	group.GET("/page", GetPageList)
}
