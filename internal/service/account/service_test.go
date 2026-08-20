package account

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/hisshihi/simple-bank/internal/domain"
	errs "github.com/hisshihi/simple-bank/internal/errors"
	"github.com/hisshihi/simple-bank/internal/service/account/mocks"
)

func TestService_Create(t *testing.T) {
	fixedTime := time.Date(2026, 7, 9, 0, 0, 0, 0, time.UTC)
	clock := func() time.Time { return fixedTime }
	owner := "owner"
	currency := "RUB"

	t.Run("success", func(t *testing.T) {
		repo := mocks.NewMockAccountRepository(t)

		repo.EXPECT().
			Create(mock.Anything, mock.MatchedBy(func(account *domain.Account) bool {
				return account.Owner() == owner &&
					string(account.Currency()) == currency &&
					account.CreatedAt().Equal(fixedTime)
			})).
			RunAndReturn(func(_ context.Context, a *domain.Account) (*domain.Account, error) {
				return a, nil
			}).
			Once()

		svc := New(repo, clock)

		record, err := svc.Create(t.Context(), owner, currency)
		require.NoError(t, err)
		require.Equal(t, owner, record.Name)
		require.NotEqual(t, uuid.Nil, record.ID)
		require.Equal(t, currency, record.Currency)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := mocks.NewMockAccountRepository(t)
		wantErr := errors.New("db is down")

		repo.EXPECT().
			Create(mock.Anything, mock.Anything).
			Return(nil, wantErr).
			Once()

		svc := New(repo, clock)
		_, err := svc.Create(t.Context(), owner, currency)
		require.ErrorIs(t, err, wantErr)
	})

	t.Run("invalid currency - repo не должен вызываться", func(t *testing.T) {
		repo := mocks.NewMockAccountRepository(t)

		svc := New(repo, clock)
		_, err := svc.Create(t.Context(), owner, "INVALID_CURRENCY")
		require.ErrorIs(t, err, errs.ErrCurrencyIsNotSupported)
	})

	t.Run("domain validation error stops before repo call", func(t *testing.T) {
		repo := mocks.NewMockAccountRepository(t)

		svc := New(repo, clock)
		_, err := svc.Create(t.Context(), "", currency)
		require.ErrorIs(t, err, errs.ErrOwnerRequired)
	})
}
