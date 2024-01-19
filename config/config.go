package config

import (
	"ai-report/common"
	"ai-report/config/log"
	"ai-report/pkg/ormx"
	"fmt"
	"github.com/spf13/viper"
)

// InitConfig Config init
func InitConfig(name, suffix, path string) {
	LoadConfigFile(name, suffix, path)
	log.InitZap()
	ormx.InitOrmx()
}

func LoadConfigFile(name, suffix, path string) {
	viper.SetConfigName(name)
	viper.SetConfigType(suffix)
	viper.AddConfigPath(path)
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	if err := viper.Unmarshal(&common.GlobalConfig); err != nil {
		panic(fmt.Errorf("Fatal error format config to json: %s \n", err))
	}
}
