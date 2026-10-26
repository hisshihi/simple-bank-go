package account

import "log/slog"

type Handler struct {
	accountSvc Service
	log        *slog.Logger
}

func NewHandler(svc Service, log *slog.Logger) *Handler {
	return &Handler{accountSvc: svc, log: log}
}
