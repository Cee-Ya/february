package mid

import (
	"ai-report/common"
	"ai-report/common/consts"
	"ai-report/config/log"
	"ai-report/entity"
	"ai-report/pkg/ginx/render"
	"bytes"
	"context"
	"errors"
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

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		uuidStr := strings.ReplaceAll(uuid.New().String(), "-", "")
		path := c.Request.URL.Path
		userId := 0
		ctx := context.WithValue(context.Background(), consts.TraceKey, &entity.Trace{TraceId: uuidStr, Caller: path, UserId: userId})
		dataByte, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewReader(dataByte))
		c.Set(consts.TraceCtx, ctx)
		c.Next()
		cost := time.Since(start)
		zapFields := make([]zap.Field, 0)
		zapFields = append(zapFields, zap.Int("Status", c.Writer.Status()))
		zapFields = append(zapFields, zap.String("Method", c.Request.Method))
		zapFields = append(zapFields, zap.String("IP", c.ClientIP()))
		zapFields = append(zapFields, zap.String("Path", path))
		zapFields = append(zapFields, zap.String("query", c.Request.URL.RawQuery))
		zapFields = append(zapFields, zap.String("body", string(dataByte)))
		if result, ok := c.Get(consts.ResponseData); ok {
			zapFields = append(zapFields, zap.String("result", result.(string)))
		}
		zapFields = append(zapFields, zap.String("UserAgent", c.Request.UserAgent()))
		zapFields = append(zapFields, zap.Duration("Cost", cost))
		log.InfoF(ctx, "[rest]", zapFields...)
	}
}

func NoRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		render.Result(c).Error(errors.New("router not found"))
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
					log.ErrorF(ctx, c.Request.URL.Path,
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
					log.ErrorF(ctx, "[panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					log.ErrorF(ctx, "[panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				render.Result(c).Error(err.(error))
				c.Abort()
			}
		}()
		c.Next()
	}
}
