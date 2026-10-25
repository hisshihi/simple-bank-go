package main

import (
	"log"
	"log/slog"

	"github.com/joho/godotenv"

	"github.com/hisshihi/simple-bank/internal/app"
	"github.com/hisshihi/simple-bank/internal/config"
	"github.com/hisshihi/simple-bank/internal/logger"
)

func main() {
	_ = godotenv.Load()
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("err to load config: %v", err)
	}
	initLog := logger.Init(parseLogLevel(cfg.LoggerLevel), cfg.LoggerPretty)
	slog.SetDefault(initLog)
	err = app.New(cfg, initLog).Run()
	if err != nil {
		log.Fatal(err)
	}
}

func parseLogLevel(logLevel string) slog.Level {
	switch logLevel {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
