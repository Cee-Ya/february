package conf

import (
	"february/gen/pkg/ormx"
	"github.com/spf13/viper"
)

var C = new(Config)

type GenConf struct {
	EnableCache       bool   // enable cache
	AbsPath           string // absolute path
	TargetPath        string // target result
	DomainTargetPath  string // domain target result
	ServiceTargetPath string // service target result
	ApiTargetPath     string // api target result
}

type Config struct {
	RunMode  string
	Database ormx.DBConfig
	Gen      GenConf
}

func MustLoad(path, configName, ty string) {
	viper.SetConfigName(configName)
	viper.SetConfigType(ty)
	viper.AddConfigPath(path)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(C); err != nil {
		panic(err)
	}
}
