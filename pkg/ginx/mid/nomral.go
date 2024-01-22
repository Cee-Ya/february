package mid

import (
	"bytes"
	"context"
	"errors"
	"february/common"
	"february/common/consts"
	"february/entity"
	"february/pkg/ginx/render"
	"february/pkg/logx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

// Cors 跨域中间件
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			allowMethod  = "GET,HEAD,PUT,POST,DELETE,OPTIONS"
			allowOrigin  = "*"
			allowHeaders = "Content-Type,AccessToken,X-CSRF-Token, Authorization, token, sign, useragent"
		)

		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", allowOrigin)
		c.Header("Access-Control-Allow-Headers", allowHeaders)
		c.Header("Access-Control-Allow-Methods", allowMethod)
		// c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

// 自定义 ResponseWriter
type noLoggingResponseWriter struct {
	gin.ResponseWriter
	bodyBuffer *bytes.Buffer
}

func (w *noLoggingResponseWriter) Write(data []byte) (int, error) {
	// 将响应写入缓冲区
	w.bodyBuffer.Write(data)
	// 调用原始 ResponseWriter 的 Write 方法
	return w.ResponseWriter.Write(data)
}

func (w *noLoggingResponseWriter) isStream() bool {
	// 判断响应体是否是流
	// 这里可以根据实际情况来判断，例如判断 Content-Type 是否为流媒体类型
	// 这里简单示范了通过缓冲区大小来判断是否是流
	return w.bodyBuffer.Len() > 1024 // 通过缓冲区大小来判断是否是流
}

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		uuidStr := strings.ReplaceAll(uuid.New().String(), "-", "")
		path := c.Request.URL.Path
		contentType := c.GetHeader("Content-Type")
		userId := 0
		ctx := context.WithValue(context.Background(), consts.TraceKey, &entity.Trace{TraceId: uuidStr, UserId: userId})
		var (
			dataByte []byte
			err      error
		)
		switch c.Request.Method {
		// 只有在以下请求方法中才读取body
		case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
			if dataByte, err = io.ReadAll(c.Request.Body); err != nil {
				logx.ErrorF(ctx, "read body err:: ", zap.Error(err))
				render.Result(c).Error(err)
				c.Abort()
				return
			}
			c.Request.Body = io.NopCloser(bytes.NewReader(dataByte))
		}
		// 提前注入traceId
		//common.Ormx = common.Ormx.WithContext(ctx)
		c.Set(consts.TraceCtx, ctx)
		// 使用自定义 ResponseWriter
		writer := &noLoggingResponseWriter{c.Writer, bytes.NewBuffer(nil)}
		// 替换 gin 的 Writer
		c.Writer = writer
		c.Next()
		cost := time.Since(start)
		zapFields := make([]zap.Field, 0)
		zapFields = append(zapFields, zap.Int("status", c.Writer.Status()))
		zapFields = append(zapFields, zap.String("method", c.Request.Method))
		zapFields = append(zapFields, zap.String("ip", c.ClientIP()))
		zapFields = append(zapFields, zap.String("path", path))
		zapFields = append(zapFields, zap.String("query", c.Request.URL.RawQuery))
		// 只有在以下请求方法中才打印body
		switch c.Request.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
			// 如果是文件上传，不打印body
			if contentType != "multipart/form-data" {
				zapFields = append(zapFields, zap.String("body", string(dataByte)))
			}
		}
		if result, ok := c.Get(consts.ResponseData); ok {
			// 检查是否是流，如果是则进行相应的处理
			if !writer.isStream() {
				// 在这里进行流相关的处理，不记录返回值日志
				zapFields = append(zapFields, zap.String("result", result.(string)))
			}
		}
		zapFields = append(zapFields, zap.String("userAgent", c.Request.UserAgent()))
		zapFields = append(zapFields, zap.Duration("cost", cost))
		logx.InfoF(ctx, "REST:: ", zapFields...)
	}
}

// NoRoute 404处理
func NoRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    consts.NotFound,
			"message": "404, page not exists!",
		})
		c.AbortWithStatus(http.StatusNotFound)
	}
}

// GinRecovery recover掉项目可能出现的panic
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				ctx := common.GetTraceCtx(c)
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					var se *os.SyscallError
					if errors.As(ne.Err, &se) {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logx.ErrorF(ctx, "PANIC::",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}
				if stack {
					debug.PrintStack()
					logx.ErrorF(ctx, "PANIC::",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.Stack(common.GlobalConfig.Log.StacktraceKey),
					)
				} else {
					logx.ErrorF(ctx, "PANIC::",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				render.Result(c).Fail(err.(error))
				c.Abort()
			}
		}()
		c.Next()
	}
}
