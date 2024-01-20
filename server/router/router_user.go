package router

import (
	"ai-report/common"
	"ai-report/pkg/ginx/render"
	"ai-report/server/service"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func GetUserList(ctx *gin.Context) {
	res, err := service.NewUserService(common.GetTraceCtx(ctx)).FindList(func(where *gorm.DB) {
		where.Order("id desc")
	})
	render.Result(ctx).Dangers(errors.Wrap(err, "user list err::")).Ok(res)
}

func GetPageList(ctx *gin.Context) {
	res, err := service.NewUserService(common.GetTraceCtx(ctx)).FindPageList(nil, GetPage(ctx))
	render.Result(ctx).Dangers(errors.Wrap(err, "user page list err::")).Ok(res)
}
