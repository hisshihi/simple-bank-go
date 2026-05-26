package main

import (
	"log"

	"github.com/hisshihi/simple-bank/config"
	"github.com/hisshihi/simple-bank/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	err = app.New(&cfg).Run()
	if err != nil {
		log.Fatal(err)
	}
}
