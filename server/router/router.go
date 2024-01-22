package router

import (
	"february/pkg/ginx/mid"
	"github.com/gin-gonic/gin"
)

// Init router
func Init() *gin.Engine {
	r := gin.New()
	r.NoRoute(mid.NoRoute())
	r.Use(mid.Cors(), mid.XSSFilter([]string{}), mid.GinLogger(), mid.GinRecovery(true))
	UserRouter(r.Group("/user"))
	return r
}

func UserRouter(group *gin.RouterGroup) {
	group.POST("add", AddUser)
	group.POST("update", UpdateUser)
	group.GET("page", GetPageList)
	group.GET("info", GetUser)
}
