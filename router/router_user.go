package router

import (
	"ai-report/common"
	"ai-report/pkg/ginx/render"
	"ai-report/service"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func GetUserList(ctx *gin.Context) {
	res, err := service.NewUserService(common.GetTraceCtx(ctx)).FindList(nil)
	render.Result(ctx).Dangers(errors.Wrap(err, "user list err::")).Ok(res)
}
