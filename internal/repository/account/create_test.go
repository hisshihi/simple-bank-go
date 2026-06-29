package account

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hisshihi/simple-bank/internal/domain"
	"github.com/hisshihi/simple-bank/internal/helper"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	os.Exit(helper.TestSetupPool(m))
}

func TestRepo_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tx := helper.NewTx(t)
		repo := New(tx)

		newCurrency, err := domain.NewCurrency("RUB")
		require.NoError(t, err)
		id, err := uuid.NewV7()
		require.NoError(t, err)
		newAccount, err := domain.NewAccount(id, "hiss", newCurrency, fixedClock)
		require.NoError(t, err)

		account, err := repo.Create(context.Background(), newAccount)
		require.NoError(t, err)
		require.NotEmpty(t, account.ID())

		findAccount, err := repo.GetByID(context.Background(), account.ID())
		require.NoError(t, err)
		require.Equal(t, findAccount.ID(), account.ID())
		require.Equal(t, findAccount.Owner(), newAccount.Owner())
		require.Equal(t, findAccount.Currency(), newAccount.Currency())
		require.Equal(t, findAccount.Balance(), newAccount.Balance())
		require.WithinDuration(t, findAccount.CreatedAt(), newAccount.CreatedAt(), time.Second)
	})
	t.Run("create error", func(t *testing.T) {
		tx := helper.NewTx(t)
		repo := New(tx)

		newCurrency, err := domain.NewCurrency("RUB")
		require.NoError(t, err)

		id, err := uuid.NewV7()
		require.NoError(t, err)

		newAccount, err := domain.NewAccount(id, "hiss", newCurrency, fixedClock)
		require.NoError(t, err)

		ctxWithTimeout, cancel := context.WithCancel(context.Background())
		cancel()

		_, err = repo.Create(ctxWithTimeout, newAccount)
		require.Error(t, err)
	})
}

var fixedClock = func() time.Time {
	return time.Date(2026, 6, 3, 0, 0, 0, 0, time.UTC)
}
