package main

import (
	"log"

	"github.com/hisshihi/simple-bank/db"
)

func main() {
	config := db.DBConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "root",
		Password: "secret",
		DBName:   "simple_bank",
		DBType:   "postgres",
		SSLMode:  "disable",
	}

	_, err := db.InitDB(config, &db.Account{}, &db.Entry{}, &db.Transaction{})
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
}
