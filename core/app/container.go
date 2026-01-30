// Package app хранение данных бд и сервисов. Работа с ними
package app

import (
	"log"

	"github.com/hisshihi/simple-bank/core/db"
	"github.com/hisshihi/simple-bank/core/repo"
	"github.com/hisshihi/simple-bank/core/service"
	"github.com/hisshihi/simple-bank/util"
	"gorm.io/gorm"
)

type Container struct {
	// DB connections
	SimpleBankDB *gorm.DB

	// Services
	AccountService *service.AccountService
}

func newConn(host, port, user, password, dbName, dbType, sslMode string, models ...interface{}) (*gorm.DB, error) {
	dbConfig := db.DBConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbName,
		DBType:   dbType,
		SSLMode:  sslMode,
	}

	database, err := db.InitDB(dbConfig, models...)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
		return nil, err
	}

	return database, nil
}

func NewContainer(config util.Config) (*Container, error) {
	simpleBankDB, err := newConn(
		config.DBHost,
		config.DBPort,
		config.DBSimpleBankUser,
		config.DBSimpleBankPassword,
		config.DBSimpleBankDatabase,
		"postgres",
		"disable",
		&db.Account{},
		&db.Entry{},
		&db.Transaction{},
	)
	if err != nil {
		return nil, err
	}

	// Repositories
	accountRepo := repo.NewAccountRepo(simpleBankDB)

	// Services
	accountService := service.NewAccountService(accountRepo)

	return &Container{
		SimpleBankDB:   simpleBankDB,
		AccountService: accountService,
	}, nil
}
