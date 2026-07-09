package account

import (
	"context"
	"net/http"

	"github.com/hisshihi/simple-bank/internal/handler/helper"
	"github.com/hisshihi/simple-bank/internal/service/account/dto"
)

type Service interface {
	Create(ctx context.Context, name, currency string) (dto.CreateRecord, error)
}

func (h *Handler) RegisterRouters(mux *http.ServeMux) {
	mux.HandleFunc("POST /create", helper.Adapt(h.Create))
}
