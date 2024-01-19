package router

import (
	"ai-report/common"
	"ai-report/config/log"
	"ai-report/pkg/ginx/render"
	"ai-report/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func List(ctx *gin.Context) {
	com := common.GetTraceCtx(ctx)
	res, err := service.NewUserService(com).FindList(nil)
	if err != nil {
		log.ErrorF(com, "find list err: %v", zap.Error(err))
		return
	}
	render.Result(ctx).Ok(res)
}
