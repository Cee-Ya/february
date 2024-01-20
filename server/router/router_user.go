package router

import (
	"february/common"
	"february/common/tools"
	"february/pkg/ginx/render"
	"february/server/service"
	"february/server/vo"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func GetUser(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := tools.Str2Uint64(idStr)
	if err != nil {
		render.Result(ctx).Dangers(err)
		return
	}
	res, err := service.NewUserService(common.GetTraceCtx(ctx)).FindById(id)
	render.Result(ctx).Dangers(errors.Wrap(err, "get user err::")).Ok(res)
}

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

func UpdateUser(ctx *gin.Context) {
	var req vo.UserUpdateVo
	if err := ctx.ShouldBindJSON(&req); err != nil {
		render.Result(ctx).Dangers(err)
		return
	}
	if err := service.NewUserService(common.GetTraceCtx(ctx)).Update(req); err != nil {
		render.Result(ctx).Dangers(err)
		return
	}
	render.Result(ctx).Ok(nil)
}

func GetPageList(ctx *gin.Context) {
	res, err := service.NewUserService(common.GetTraceCtx(ctx)).PageList(GetPage(ctx))
	render.Result(ctx).Dangers(errors.Wrap(err, "user page list err::")).Ok(res)
}
