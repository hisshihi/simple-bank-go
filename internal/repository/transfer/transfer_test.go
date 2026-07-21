package transfer

import (
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hisshihi/simple-bank/internal/domain"
	"github.com/hisshihi/simple-bank/internal/helper"
	"github.com/hisshihi/simple-bank/internal/repository/account"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	os.Exit(helper.TestSetupPool(m))
}

func TestRepo_Create(t *testing.T) {
	tx := helper.NewTx(t)
	repo := New(tx)
	accountRepo := account.New(tx)
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
