package common

import (
	"ai-report/common/consts"
	"context"
	"github.com/gin-gonic/gin"
)

// Trace 定义trace结构体
type Trace struct {
	TraceId   string  `json:"trace_id"`
	SpanId    string  `json:"span_id"`
	Caller    string  `json:"caller"`
	SrcMethod *string `json:"srcMethod,omitempty"`
	UserId    int     `json:"user_id"`
}

// GetTraceCtx 根据gin的context获取context，使log trace更加通用
func GetTraceCtx(c *gin.Context) context.Context {
	return c.MustGet(consts.TraceCtx).(context.Context)
}
