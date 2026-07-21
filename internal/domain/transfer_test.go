package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	errs "github.com/hisshihi/simple-bank/internal/errors"
	"github.com/stretchr/testify/require"
)

func TestNewTransfer(t *testing.T) {
	transferID, fromID, toID := NewUUID(t), NewUUID(t), NewUUID(t)
	fixedTime := time.Date(2026, 7, 21, 0, 0, 0, 0, time.UTC)
	clock := func() time.Time { return fixedTime }
	rub := Money{amount: 100, currency: Currency("RUB")}

	t.Run("success", func(t *testing.T) {
		got, err := NewTransfer(transferID, fromID, toID, rub, clock)
		require.NoError(t, err)

		require.Equal(t, transferID, got.ID())
		require.Equal(t, fromID, got.FromAccountID())
		require.Equal(t, toID, got.ToAccountID())
		require.Equal(t, rub, got.Amount())
		require.Equal(t, fixedTime, got.CreatedAt())

		entries := got.Entries()
		require.Len(t, entries, 2)

		debit, credit := entries[0], entries[1]

		require.Equal(t, EntryDirectionDebit, debit.Direction())
		require.Equal(t, fromID, debit.AccountID())
		require.Equal(t, transferID, debit.TransferID())

		require.Equal(t, EntryDirectionCredit, credit.Direction())
		require.Equal(t, toID, credit.AccountID())
		require.Equal(t, transferID, credit.TransferID())

		require.NotEqual(t, debit.ID(), credit.ID())
	})

	tests := []struct {
		name    string
		id      uuid.UUID
		from    uuid.UUID
		to      uuid.UUID
		amount  Money
		wantErr error
	}{
		{"invalid transfer id", uuid.Nil, fromID, toID, rub, errs.ErrInvalidTransferID},
		{"invalid from account id", transferID, uuid.Nil, toID, rub, errs.ErrInvalidFromAccountID},
		{"invalid to account id", transferID, fromID, uuid.Nil, rub, errs.ErrInvalidToAccountID},
		{"amount is zero", transferID, fromID, toID, Money{amount: 0, currency: Currency("RUB")}, errs.ErrZeroAmount},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTransfer(tt.id, tt.from, tt.to, tt.amount, clock)
			require.ErrorIs(t, err, tt.wantErr)
			require.Nil(t, got)
		})
	}
}
