package server

import (
	"february/common"
	"february/common/consts"
	"february/pkg/cache/redisx"
	"february/pkg/config"
	"february/pkg/httpx"
	"february/pkg/logx"
	"february/pkg/ormx"
	"february/server/router"
	"strings"
)

func Initialize(path, configName string) (func(), error) {
	if configName == "" {
		configName = "default.toml"
	}
	names := strings.Split(configName, consts.DOT)
	config.LoadConfigFile(names[0], names[1], path)
	logx.Init()

	// init database
	dbClean, err := ormx.Init()
	if err != nil {
		return nil, err
	}

	// init redis
	// todo 选择缓存模式
	//memory.InitMemoryCache()
	if err = redisx.InitRedis(common.GlobalConfig.Cache); err != nil {
		return nil, err
	}

	// init http server
	r := router.Init(common.GlobalConfig.Server)
	httpClean := httpx.Init(common.GlobalConfig.Server, r)

	// release all the resources
	return func() {
		dbClean()
		httpClean()
	}, nil
}
