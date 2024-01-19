package main

import (
	"ai-report/common"
	"ai-report/config"
	"ai-report/pkg/ginx/mid"
	"ai-report/router"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	// 初始化配置
	config.InitConfig("default", "toml", "./")

	r := gin.New()
	r.NoRoute(mid.NoRoute())
	r.Use(mid.Cors(), mid.GinLogger(), mid.GinRecovery(true))
	router.UserRouter(r.Group("/user"))
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", common.GlobalConfig.Server.Port),
		Handler: r,
	}
	if err := srv.ListenAndServe(); err != nil {
		panic(fmt.Errorf("start server error: %s \n", err))
	}
}
