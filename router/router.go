package router

import "github.com/gin-gonic/gin"

func AuthRouter(group *gin.RouterGroup) {
	group.GET("/list", List)
}