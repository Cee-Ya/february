package router

import (
	"ai-report/entity"
	"github.com/gin-gonic/gin"
)

// GetPage 获取分页参数
func GetPage(c *gin.Context) *entity.Page {
	page := &entity.Page{}
	if err := c.ShouldBind(page); err != nil {
		panic(err)
	}
	if page.PageNo == 0 {
		page.PageNo = 1
	}
	if page.PageSize == 0 {
		page.PageSize = 10
	}
	return page
}
