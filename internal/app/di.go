package app

import (
	"log/slog"
	"os"

	"github.com/hisshihi/simple-bank/config"
	"github.com/hisshihi/simple-bank/internal/closer"
	"github.com/hisshihi/simple-bank/internal/database"
)

type DiContainer struct {
	config *config.Config

	db *database.DB
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
