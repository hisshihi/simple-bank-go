package app

import (
	"log/slog"
	"os"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"

	"github.com/hisshihi/simple-bank/internal/closer"
	"github.com/hisshihi/simple-bank/internal/config"
	"github.com/hisshihi/simple-bank/internal/database"
	accounthandler "github.com/hisshihi/simple-bank/internal/handler/account"
	accountrepository "github.com/hisshihi/simple-bank/internal/repository/account"
	entryrepository "github.com/hisshihi/simple-bank/internal/repository/entry"
	transferrepository "github.com/hisshihi/simple-bank/internal/repository/transfer"
	accountservice "github.com/hisshihi/simple-bank/internal/usecase/account"
	transferservice "github.com/hisshihi/simple-bank/internal/usecase/transfer"
)

type DiContainer struct {
	cfg *config.Config

	db *database.DB

	log *slog.Logger

	// repo
	accountRepo  *accountrepository.Repo
	transferRepo *transferrepository.Repo
	entryRepo    *entryrepository.Repo

	//	service
	accountService  *accountservice.Service
	transferService *transferservice.Service

	//	handler
	accountHandler *accounthandler.Handler

	// other
	txManager *manager.Manager
}

func NewDIContainer(cfg *config.Config, log *slog.Logger) *DiContainer {
	return &DiContainer{
		cfg: cfg,
		log: log,
	}
}

func (c *DiContainer) DB() *database.DB {
	if c.db == nil {
		dsn := c.cfg.DSN()
		db, err := database.New(dsn)
		if err != nil {
			c.log.Error("ошибка подключения к бд", "err", err)
			os.Exit(1)
		}

		closer.Add("база данных", db.Close)

		c.db = db
	}
	return c.db
}

func (c *DiContainer) AccountRepo() *accountrepository.Repo {
	if c.accountRepo == nil {
		c.accountRepo = accountrepository.New(c.DB(), trmpgx.DefaultCtxGetter)
	}
	return c.accountRepo
}

func (c *DiContainer) TransferRepo() *transferrepository.Repo {
	if c.transferRepo == nil {
		c.transferRepo = transferrepository.New(c.DB(), trmpgx.DefaultCtxGetter)
	}
	return c.transferRepo
}

func (c *DiContainer) EntryRepo() *entryrepository.Repo {
	if c.entryRepo == nil {
		c.entryRepo = entryrepository.New(c.DB(), trmpgx.DefaultCtxGetter)
	}
	return c.entryRepo
}

func (c *DiContainer) AccountService() *accountservice.Service {
	if c.accountService == nil {
		c.accountService = accountservice.New(c.AccountRepo(), func() time.Time { return time.Now().UTC() }, c.log)
	}
	return c.accountService
}

func (c *DiContainer) AccountHandler() *accounthandler.Handler {
	if c.accountHandler == nil {
		c.accountHandler = accounthandler.NewHandler(c.AccountService(), c.log)
	}
	return c.accountHandler
}

func (c *DiContainer) TransferService() *transferservice.Service {
	if c.transferService == nil {
		c.transferService = transferservice.NewService(c.TxManager(), c.TransferRepo(), func() time.Time { return time.Now().UTC() }, c.AccountRepo())
	}
	return c.transferService
}

func (c *DiContainer) TxManager() *manager.Manager {
	if c.txManager == nil {
		c.txManager = manager.Must(trmpgx.NewDefaultFactory(c.DB().Pool()))
	}
	return c.txManager
}
