package main

import (
	"log"

	"github.com/hisshihi/simple-bank/core/app"
	"github.com/hisshihi/simple-bank/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	_, err = app.NewContainer(config)
	if err != nil {
		log.Fatalf("failed to create app container: %v", err)
	}
}
