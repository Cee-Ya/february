package conf

import (
	"february/gen/pkg/ormx"
	"github.com/spf13/viper"
)

var C = new(Config)

type GenConf struct {
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

func MustLoad(path string) {
	viper.SetConfigName("default")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./gen/")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(C); err != nil {
		panic(err)
	}
}
