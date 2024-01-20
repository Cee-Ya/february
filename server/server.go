package server

import (
	"february/common"
	"february/pkg/config"
	"february/pkg/httpx"
	"february/pkg/logx"
	"february/pkg/ormx"
	"february/pkg/redis"
	"february/server/router"
	"strings"
)

func Initialize(path, configName string) (func(), error) {
	if configName == "" {
		configName = "default.toml"
	}
	names := strings.Split(configName, ".")
	config.LoadConfigFile(names[0], names[1], path)
	logx.Init()

	// init database
	if err := ormx.Init(); err != nil {
		return nil, err
	}

	// init redis
	if err := redisx.InitRedis(common.GlobalConfig.Redis); err != nil {
		return nil, err
	}

	// init http server
	r := router.Init()
	httpClean := httpx.Init(common.GlobalConfig.Server, r)

	// release all the resources
	return func() {
		httpClean()
	}, nil
}
