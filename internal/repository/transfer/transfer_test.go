package transfer

import (
	"context"
	"os"
	"testing"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/hisshihi/simple-bank/internal/domain"
	"github.com/hisshihi/simple-bank/internal/helper"
	"github.com/hisshihi/simple-bank/internal/repository/account"
)

func TestMain(m *testing.M) {
	os.Exit(helper.TestSetupPool(m))
}

func TestRepo_Create(t *testing.T) {
	tx := helper.NewTx(t)
	repo := New(tx, trmpgx.DefaultCtxGetter)
	accountRepo := account.New(tx, trmpgx.DefaultCtxGetter)
	t.Run("success", func(t *testing.T) {
		newCurrency, err := domain.NewCurrency("RUB")
		require.NoError(t, err)

		newMoney, err := domain.NewMoney(100, newCurrency)
		require.NoError(t, err)

		from, err := domain.NewAccount(NewUUID(t), "Denis", newCurrency, fixedClock)
		require.NoError(t, err)
		_, err = accountRepo.Create(t.Context(), from)
		require.NoError(t, err)

		to, err := domain.NewAccount(NewUUID(t), "Arina", newCurrency, fixedClock)
		require.NoError(t, err)
		_, err = accountRepo.Create(t.Context(), to)
		require.NoError(t, err)

		transfer, err := domain.NewTransfer(NewUUID(t), from.ID(), to.ID(), newMoney, fixedClock)
		require.NoError(t, err)
		require.NoError(t, repo.Create(t.Context(), transfer))
	})
	t.Run("failure", func(t *testing.T) {
		newCurrency, err := domain.NewCurrency("RUB")
		require.NoError(t, err)

		newMoney, err := domain.NewMoney(100, newCurrency)
		require.NoError(t, err)

		from, err := domain.NewAccount(NewUUID(t), "Denis", newCurrency, fixedClock)
		require.NoError(t, err)
		_, err = accountRepo.Create(t.Context(), from)
		require.NoError(t, err)

		to, err := domain.NewAccount(NewUUID(t), "Arina", newCurrency, fixedClock)
		require.NoError(t, err)
		_, err = accountRepo.Create(t.Context(), to)
		require.NoError(t, err)

		transfer, err := domain.NewTransfer(NewUUID(t), from.ID(), to.ID(), newMoney, fixedClock)
		require.NoError(t, err)
		ctxWithCancel, cancel := context.WithCancel(t.Context())
		cancel()
		require.Error(t, repo.Create(ctxWithCancel, transfer))
	})
}

func NewUUID(t *testing.T) uuid.UUID {
	t.Helper()
	newUUID, err := uuid.NewV7()
	require.NoError(t, err)
	return newUUID
}

var fixedClock = func() time.Time {
	return time.Date(2026, 6, 3, 0, 0, 0, 0, time.UTC)
}
