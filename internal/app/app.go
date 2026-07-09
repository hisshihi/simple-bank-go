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

type App struct {
	diContainer *DiContainer
	httpServer  *http.Server
}

func New(config *config.Config) *App {
	a := &App{
		diContainer: NewDIContainer(config),
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
		Addr:    a.diContainer.config.HTTPAddr,
		Handler: mux,
	}
}

func (a *App) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("ошибка старта сервера", "error", err)
		}
	}()

	slog.Info("сервис запущен", "httpAddr", a.diContainer.config.HTTPAddr)

	<-ctx.Done()
	slog.Info("получен сигнал, завершаем...")

	stop()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("ошибка при остановке сервера", "error", err)
	}

	closerCtx, closerCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer closerCancel()

	if err := closer.CloseAll(closerCtx); err != nil {
		slog.Error("ошибки при закрытии ресурсов", "err", err)
	}

	slog.Info("сервер остановлен")

	return nil
}
