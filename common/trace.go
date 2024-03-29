package common

import (
	"context"
	"february/common/consts"
	"github.com/gin-gonic/gin"
)

// GetTraceCtx 根据gin的context获取context，使log trace更加通用
func GetTraceCtx(c *gin.Context) context.Context {
	return c.MustGet(consts.TraceCtx).(context.Context)
}
