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
	render.Result(ctx).
		Dangers(err).
		DangersRender(service.NewUserService(common.GetTraceCtx(ctx)).FindById(id))
}

func AddUser(ctx *gin.Context) {
	var req vo.UserAddVo
	render.Result(ctx).
		Dangers(ctx.ShouldBindJSON(&req)).
		Dangers(errors.Wrap(service.NewUserService(common.GetTraceCtx(ctx)).Create(req), "add user err:")).
		Ok(nil)
}

func UpdateUser(ctx *gin.Context) {
	var req vo.UserUpdateVo
	render.Result(ctx).
		Dangers(ctx.ShouldBindJSON(&req)).
		Dangers(errors.Wrap(service.NewUserService(common.GetTraceCtx(ctx)).Update(req), "update user err:")).
		Ok(nil)
}

func GetPageList(ctx *gin.Context) {
	render.Result(ctx).DangersRender(service.NewUserService(common.GetTraceCtx(ctx)).PageList(GetPage(ctx)))
}
