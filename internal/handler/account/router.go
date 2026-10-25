package account

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/hisshihi/simple-bank/internal/handler/helper"
	"github.com/hisshihi/simple-bank/internal/usecase/account/dto"
)

type Service interface {
	Create(ctx context.Context, name, currency string) (dto.CreateRecord, error)
	GetByID(ctx context.Context, id uuid.UUID) (dto.GetRecord, error)
}

func (h *Handler) RegisterRouters(mux *http.ServeMux) {
	mux.HandleFunc("POST /create", helper.Adapt(h.Create))
	mux.HandleFunc("GET /get/{id}", helper.Adapt(h.Get))
}
