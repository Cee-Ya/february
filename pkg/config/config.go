package config

import (
	"ai-report/common"
	"fmt"
	"github.com/spf13/viper"
)

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
