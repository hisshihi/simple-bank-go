package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/hisshihi/simple-bank/internal/closer"
	"github.com/hisshihi/simple-bank/internal/config"
)

const (
	readHeaderTimeout  = 5 * time.Second
	shutdownCtxTimeout = 15 * time.Second
	closerCtxTimeout   = 10 * time.Second
)

type App struct {
	diContainer *DiContainer
	httpServer  *http.Server
	log         *slog.Logger
}

func New(cfg *config.Config, log *slog.Logger) *App {
	a := &App{
		diContainer: NewDIContainer(cfg, log),
		log:         log,
	}

	a.initDeps()
	return a
}

func (a *App) initDeps() {
	funks := []func(){
		a.initHTTPServer,
	}
	for _, fn := range funks {
		fn()
	}
}

func (a *App) initHTTPServer() {
	mux := http.NewServeMux()

	a.diContainer.AccountHandler().RegisterRouters(mux)

	a.httpServer = &http.Server{
		Addr:              a.diContainer.cfg.HTTPAddr,
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
	}
}

func (a *App) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.log.Error("ошибка старта сервера", "err", err)
		}
	}()

	if a.diContainer.DB() == nil {
		return errors.New("db is nil")
	}

	a.log.Info("сервис запущен", "http_addr", a.diContainer.cfg.HTTPAddr)

	<-ctx.Done()
	a.log.Info("получен сигнал, завершаем...")

	stop()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownCtxTimeout)
	defer shutdownCancel()

	if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
		a.log.Error("ошибка при остановке сервера", "err", err)
	}

	closerCtx, closerCancel := context.WithTimeout(context.Background(), closerCtxTimeout)
	defer closerCancel()

	if err := closer.CloseAll(closerCtx, a.log); err != nil {
		a.log.Error("ошибки при закрытии ресурсов", "err", err)
	}

	a.log.Info("сервер остановлен")

	return nil
}
