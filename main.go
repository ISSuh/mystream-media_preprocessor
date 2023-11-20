package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/ISSuh/mystream-media_preprocessor/internal/configure"
	"github.com/ISSuh/mystream-media_preprocessor/internal/service"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatal("need configure file path.")
		return
	}

	configureFilePath := args[0]
	configure, err := configure.LoadConfigure(configureFilePath)
	if err != nil {
		log.Fatal("configure parse error. ", err)
		return
	}

	service := service.NewService(configure)
	service.Run()
}
