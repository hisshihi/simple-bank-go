package app

import (
	"log/slog"
	"os"
	"time"

	"github.com/hisshihi/simple-bank/internal/closer"
	"github.com/hisshihi/simple-bank/internal/config"
	"github.com/hisshihi/simple-bank/internal/database"
	accounthandler "github.com/hisshihi/simple-bank/internal/handler/account"
	accountrepository "github.com/hisshihi/simple-bank/internal/repository/account"
	accountservice "github.com/hisshihi/simple-bank/internal/service/account"
)

type DiContainer struct {
	config *config.Config

	db *database.DB

	// repo
	accountRepo *accountrepository.Repo

	//	service
	accountService *accountservice.Service

	//	handler
	accountHandler *accounthandler.Handler
}

func NewDIContainer(config *config.Config) *DiContainer {
	return &DiContainer{
		config: config,
	}
}

func (c *DiContainer) DB() *database.DB {
	if c.db == nil {
		dsn := c.config.DSN()
		db, err := database.New(dsn)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		closer.Add("база данных", db.Close)

		c.db = db
	}
	return c.db
}

func (c *DiContainer) AccountRepo() *accountrepository.Repo {
	if c.accountRepo == nil {
		c.accountRepo = accountrepository.New(c.DB())
	}
	return c.accountRepo
}

func (c *DiContainer) AccountService() *accountservice.Service {
	if c.accountService == nil {
		c.accountService = accountservice.New(c.AccountRepo(), func() time.Time { return time.Now().UTC() })
	}
	return c.accountService
}

func (c *DiContainer) AccountHandler() *accounthandler.Handler {
	if c.accountHandler == nil {
		c.accountHandler = accounthandler.NewHandler(c.AccountService())
	}
	return c.accountHandler
}
