package app

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"github.com/hisshihi/simple-bank/config"
	"github.com/hisshihi/simple-bank/internal/closer"
)

type App struct {
	diContainer *DiContainer
}

func New(config *config.Config) *App {
	a := &App{
		diContainer: NewDIContainer(config),
	}

	a.initDeps()
	return a
}

func (a *App) initDeps() {
	var funks []func()
	for _, fn := range funks {
		fn()
	}
}

func (a *App) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	slog.Info("получен сигнал, завершаем...")

	stop()

	_, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	slog.Info("сервер остановлен")

	closerCtx, closerCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer closerCancel()

	if err := closer.CloseAll(closerCtx); err != nil {
		slog.Error("ошибки при закрытии ресурсов", "err", err)
	}

	return nil
}
