package account

import (
	"encoding/json"
	"net/http"

	"github.com/hisshihi/simple-bank/internal/handler/helper"
)

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) error {
	var reqBody struct {
		Name     string `json:"name"`
		Currency string `json:"currency"`
	}

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		return helper.BadRequest("invalid_request", "Некорректное тело запроса.")
	}

	record, err := h.accountSvc.Create(r.Context(), reqBody.Name, reqBody.Currency)
	if err != nil {
		return err
	}

	return helper.RespondJSON(w, http.StatusCreated, record)
}
