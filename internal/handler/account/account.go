package account

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

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
		slog.Error("ошибка при создании аккаунта", "name", reqBody.Name, "currency", reqBody.Currency, "error", err)
		return err
	}

	return helper.RespondJSON(w, http.StatusCreated, record)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")

	uuidID, err := uuid.Parse(id)
	if err != nil {
		return helper.BadRequest("invalid_argument", "ID аккаунта не валиден")
	}

	record, err := h.accountSvc.GetByID(r.Context(), uuidID)
	if err != nil {
		return err
	}

	return helper.RespondJSON(w, http.StatusOK, record)
}
