package render

import (
	"ai-report/common"
	"ai-report/common/consts"
	"ai-report/common/tools"
	"ai-report/entity"
	"ai-report/pkg/logx"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

type Response struct {
	Ctx       *gin.Context `json:"-"`
	Code      uint16       `json:"code"`
	Message   string       `json:"message"`
	Data      *interface{} `json:"data"`
	Timestamp int64        `json:"timestamp"`
	Sno       string       `json:"sno"`
}

func Result(ctx *gin.Context) *Response {
	return &Response{
		Timestamp: time.Now().UnixMilli(),
		Ctx:       ctx,
	}
}

// Dangers 用于处理错误
func (r *Response) Dangers(err error) *Response {
	if err != nil {
		ctx := common.GetTraceCtx(r.Ctx)
		logx.ErrorF(ctx, "Dangers:: ", zap.Error(err))
		r.Error(err)
	}
	return r
}

func (r *Response) Ok(data interface{}) {
	if r.Code > 0 {
		return
	}
	r.Code = consts.Success
	r.Data = &data
	r.render()
}

func (r *Response) Error(err error) {
	r.Code = consts.Error
	r.Message = err.Error()
	r.render()
}

func (r *Response) Fail(message string) {
	r.Code = consts.Failed
	r.Message = message
	r.render()
}

func (r *Response) Warn(message string) {
	r.Code = consts.Warn
	r.Message = message
	r.render()
}

func (r *Response) toString() string {
	data, err := tools.ToJson(r)
	if err != nil {
		logx.ErrorF(r.Ctx, "response to json err:: ", zap.Error(err))
		r.Error(err)
		return r.toString()
	}
	return data
}

func (r *Response) render() {
	if r.Code != consts.Success {
		r.Message = consts.ResponseMap[r.Code] + ": " + r.Message
	} else {
		r.Message = consts.ResponseMap[r.Code]
	}
	ctx := common.GetTraceCtx(r.Ctx)
	trace := ctx.Value(consts.TraceKey).(*entity.Trace)
	r.Sno = trace.TraceId
	r.Ctx.Set(consts.ResponseData, r.toString())
	r.Ctx.JSON(consts.Success, r)
}
