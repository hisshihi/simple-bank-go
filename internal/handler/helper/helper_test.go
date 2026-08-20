package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	errs "github.com/hisshihi/simple-bank/internal/errors"
)

func TestRespondError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantCode   string
	}{
		{
			name:       "known sentinel",
			err:        errs.ErrInvalidAccountID,
			wantStatus: http.StatusBadRequest,
			wantCode:   "invalid_account_id",
		},
		{
			name:       "wrapped sentinel",
			err:        fmt.Errorf("service: %w", errs.ErrAccountNotFound),
			wantStatus: http.StatusNotFound,
			wantCode:   "account_not_found",
		},
		{
			name:       "unknown error fails back to 500",
			err:        errors.New("unknown error"),
			wantStatus: http.StatusInternalServerError,
			wantCode:   "server_error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			RespondError(rec, tt.err)

			require.Equal(t, tt.wantStatus, rec.Code)

			var body map[string]string
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
			require.Equal(t, tt.wantCode, body["error"])
		})
	}
}
