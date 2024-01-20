package router

import (
	"ai-report/common"
	"ai-report/pkg/ginx/render"
	"ai-report/server/service"
	"ai-report/server/vo"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func AddUser(ctx *gin.Context) {
	var req vo.UserAddVo
	if err := ctx.ShouldBindJSON(&req); err != nil {
		render.Result(ctx).Dangers(err)
		return
	}
	if err := service.NewUserService(common.GetTraceCtx(ctx)).Create(req); err != nil {
		render.Result(ctx).Dangers(err)
		return
	}
	render.Result(ctx).Ok(nil)
}

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
