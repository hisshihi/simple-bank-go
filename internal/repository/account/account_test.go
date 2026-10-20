package account

import (
	"context"
	"os"
	"testing"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/hisshihi/simple-bank/internal/domain"
	errs "github.com/hisshihi/simple-bank/internal/errors"
	"github.com/hisshihi/simple-bank/internal/helper"
)

func TestMain(m *testing.M) {
	os.Exit(helper.TestSetupPool(m))
}

func TestRepo_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tx := helper.NewTx(t)
		repo := New(tx, trmpgx.DefaultCtxGetter)

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
		repo := New(tx, trmpgx.DefaultCtxGetter)

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

func TestRepo_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tx := helper.NewTx(t)
		repo := New(tx, trmpgx.DefaultCtxGetter)

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
		repo := New(tx, trmpgx.DefaultCtxGetter)

		fakeID, err := uuid.NewV7()
		require.NoError(t, err)

		findAccount, err := repo.GetByID(context.Background(), fakeID)
		require.ErrorIs(t, err, errs.ErrAccountNotFound)
		require.Nil(t, findAccount)
	})
	t.Run("other error", func(t *testing.T) {
		tx := helper.NewTx(t)
		repo := New(tx, trmpgx.DefaultCtxGetter)

		newCurrency, err := domain.NewCurrency("RUB")
		require.NoError(t, err)

		id, err := uuid.NewV7()
		require.NoError(t, err)
		newAccount, err := domain.NewAccount(id, "hiss", newCurrency, fixedClock)
		require.NoError(t, err)

		account, err := repo.Create(context.Background(), newAccount)
		require.NoError(t, err)

		ctxWithCancel, cancel := context.WithCancel(context.Background())
		cancel()

		_, err = repo.GetByID(ctxWithCancel, account.ID())
		require.Error(t, err)
	})
}

func TestRepo_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tx := helper.NewTx(t)
		repo := New(tx, trmpgx.DefaultCtxGetter)

		account := CreateAccount(t, tx)

		err := repo.Delete(t.Context(), account.ID())
		require.NoError(t, err)

		findAccount, err := repo.GetByID(t.Context(), account.ID())
		require.ErrorIs(t, err, errs.ErrAccountNotFound)
		require.Nil(t, findAccount)
	})
	t.Run("db error", func(t *testing.T) {
		tx := helper.NewTx(t)
		repo := New(tx, trmpgx.DefaultCtxGetter)
		ctxWithCancel, cancel := context.WithCancel(t.Context())
		cancel()
		err := repo.Delete(ctxWithCancel, uuid.New())
		require.Error(t, err)
	})
}

func TestRepo_UpdateBalance(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tx := helper.NewTx(t)
		repo := New(tx, trmpgx.DefaultCtxGetter)

		account := CreateAccount(t, tx)

		newBalance := account.Balance().Amount() + 10

		err := repo.UpdateBalance(t.Context(), account.ID(), newBalance)
		require.NoError(t, err)

		findAccount, err := repo.GetByID(t.Context(), account.ID())
		require.NoError(t, err)
		require.Equal(t, newBalance, findAccount.Balance().Amount())
	})
}

func TestRepo_Save(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tx := helper.NewTx(t)
		repo := New(tx, trmpgx.DefaultCtxGetter)

		account := CreateAccount(t, tx)

		amount, err := domain.NewMoney(20, domain.CurrencyTypeRUB)
		require.NoError(t, err)
		err = account.Deposit(amount)
		require.NoError(t, err)

		err = repo.Save(t.Context(), &account)
		require.NoError(t, err)

		findAccount, err := repo.GetByID(t.Context(), account.ID())
		require.NoError(t, err)
		require.Equal(t, amount.Amount(), findAccount.Balance().Amount())
	})
}

func CreateAccount(t *testing.T, tx trmpgx.Tr) domain.Account {
	t.Helper()
	repo := New(tx, trmpgx.DefaultCtxGetter)

	newCurrency, err := domain.NewCurrency("RUB")
	require.NoError(t, err)

	id, err := uuid.NewV7()
	require.NoError(t, err)
	newAccount, err := domain.NewAccount(id, "hiss", newCurrency, fixedClock)
	require.NoError(t, err)

	account, err := repo.Create(context.Background(), newAccount)
	require.NoError(t, err)

	return *account
}
