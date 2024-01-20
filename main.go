package main

import (
	"ai-report/server"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cleanFunc, err := server.Initialize("./", "default.toml")
	if err != nil {
		log.Fatalln("failed to initialize:", err)
	}
	code := 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

EXIT:
	for {
		sig := <-sc
		fmt.Println("received signal:", sig.String())
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			code = 0
			break EXIT
		case syscall.SIGHUP:
			// reload configuration?
		default:
			break EXIT
		}
	}

	cleanFunc()
	fmt.Println("process exited")
	os.Exit(code)
}
