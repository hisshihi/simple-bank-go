package account

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hisshihi/simple-bank/internal/domain"
	"github.com/hisshihi/simple-bank/internal/helper"
	"github.com/stretchr/testify/require"
)

func TestRepo_GetByID(t *testing.T) {
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

		findAccount, err := repo.GetByID(context.Background(), account.ID())
		require.NoError(t, err)
		require.Equal(t, account.ID(), findAccount.ID())
		require.Equal(t, account.Owner(), findAccount.Owner())
		require.Equal(t, account.Balance(), findAccount.Balance())
		require.Equal(t, account.Currency(), findAccount.Currency())
		require.Equal(t, account.Status(), findAccount.Status())
		require.WithinDuration(t, account.CreatedAt(), findAccount.CreatedAt(), time.Second)
	})
	t.Run("account not found", func(t *testing.T) {
		tx := helper.NewTx(t)
		repo := New(tx)

		fakeID, err := uuid.NewV7()
		require.NoError(t, err)

		findAccount, err := repo.GetByID(context.Background(), fakeID)
		require.NoError(t, err)
		require.Nil(t, findAccount)
	})
	t.Run("other error", func(t *testing.T) {
		_, err := domain.NewCurrency("YAN")
		require.Error(t, err)

	})
}
