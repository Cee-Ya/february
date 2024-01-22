package router

import (
	"february/common/tools"
	"february/pkg/ginx/render"
	"february/server/service"
	"february/server/vo"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func GetUser(ctx *gin.Context) {
	id, err := tools.Str2Uint64(ctx.Query("id"))
	render.Result(ctx).
		Danger(err).
		DangerRender(service.NewUserService(ctx).FindById(id))
}

func AddUser(ctx *gin.Context) {
	var req vo.UserAddVo
	render.Result(ctx).
		Dangers(ctx.ShouldBindJSON(&req),
			errors.Wrap(service.NewUserService(ctx).Create(req), "add user err:")).
		Ok(nil)
}

func UpdateUser(ctx *gin.Context) {
	var req vo.UserUpdateVo
	render.Result(ctx).
		Dangers(ctx.ShouldBindJSON(&req),
			errors.Wrap(service.NewUserService(ctx).Update(req), "update user err:")).
		Ok(nil)
}

func GetPageList(ctx *gin.Context) {
	render.Result(ctx).DangerRender(service.NewUserService(ctx).PageList(GetPage(ctx)))
}
