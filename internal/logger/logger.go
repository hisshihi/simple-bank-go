package logger

import (
	"context"
	"log/slog"
	"os"
)

type contextKey struct{}

func Init(level slog.Level, pretty bool) {
	var handler slog.Handler
	opts := &slog.HandlerOptions{Level: level, AddSource: true}
	if pretty {
		handler = slog.NewTextHandler(os.Stderr, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stderr, opts)
	}
	slog.SetDefault(slog.New(handler))
}

// WithContext кладёт логгер в контекст.
func WithContext(ctx context.Context, log *slog.Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, log)
}

// FromContext достаёт логгер из контекста.
// Если его там нет, то возвращает дефолный логгер.
func FromContext(ctx context.Context) *slog.Logger {
	if log, ok := ctx.Value(contextKey{}).(*slog.Logger); ok {
		return log
	}
	return slog.Default()
}

// With помещает аргументы в логгер и кладёт логгер в контекст.
func With(ctx context.Context, args ...any) context.Context {
	log := FromContext(ctx).With(args...)
	return WithContext(ctx, log)
}
