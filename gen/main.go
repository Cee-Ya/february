package main

import (
	"february/gen/core"
	"february/gen/pkg/conf"
)

func main() {
	// init conf and db
	conf.MustLoad("./gen/", "default", "toml")
	// generate code
	core.GenCode([]string{"t_sys_user"})
}
