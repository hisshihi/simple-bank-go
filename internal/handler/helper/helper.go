package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	errs "github.com/hisshihi/simple-bank/internal/errs"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

// Adapt - единственное место, где ошибка превращается в HTTP-ответ.
func Adapt(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			RespondError(w, err)
		}
	}
}

func RespondJSON(w http.ResponseWriter, status int, payload any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		return fmt.Errorf("encode payload: %w", err)
	}
	return nil
}

type errorInfo struct {
	status      int
	code        string
	description string
}

var errorMapping = map[error]errorInfo{
	errs.ErrInvalidAccountID: {
		status:      http.StatusBadRequest,
		code:        "invalid_account_id",
		description: "ID аккаунта передан не верно, проверьте передаваемые данные.",
	},
	errs.ErrCurrencyIsNotSupported: {
		status:      http.StatusBadRequest,
		code:        "currency_is_not_supported",
		description: "Данная валюта в настоящее время не поддерживается, укажите другую валюту.",
	},
	errs.ErrOwnerRequired: {
		status:      http.StatusBadRequest,
		code:        "owner_required",
		description: "Владелец аккаунта обязателен при создании.",
	},
	errs.ErrCurrencyMismatch: {
		status:      http.StatusBadRequest,
		code:        "currency_mismatch",
		description: "Валюты счетов не совпадают.",
	},
	errs.ErrInsufficientFunds: {
		status:      http.StatusConflict,
		code:        "insufficient_funds",
		description: "Недостаточно средств на счету.",
	},
	errs.ErrBalanceIsLessThanZero: {
		status:      http.StatusConflict,
		code:        "balance_is_less_than_zero",
		description: "Баланс не может быть меньше нуля.",
	},
	errs.ErrAccountNotFound: {
		status:      http.StatusNotFound,
		code:        "account_not_found",
		description: "Аккаунт не найден",
	},
}

type errorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func RespondError(w http.ResponseWriter, err error) {
	for sentinel, info := range errorMapping {
		if errors.Is(err, sentinel) {
			_ = RespondJSON(w, info.status, errorResponse{Error: info.code, ErrorDescription: info.description})
			return
		}
	}

	_ = RespondJSON(w, http.StatusInternalServerError, errorResponse{"server_error", "Внутренняя ошибка сервера. Пожалуйста, попробуйте позже или свяжитесь с поддержкой."})
}

type HTTPError struct {
	Status      int
	Code        string
	Description string
}

func (e *HTTPError) Error() string {
	return e.Description
}

func BadRequest(code, description string) *HTTPError {
	return &HTTPError{http.StatusBadRequest, code, description}
}
