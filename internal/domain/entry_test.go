package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	errs "github.com/hisshihi/simple-bank/internal/errors"
	"github.com/stretchr/testify/require"
)

func TestNewEntry(t *testing.T) {
	entryID, accountID, transferID := NewUUID(t), NewUUID(t), NewUUID(t)
	fixedTime := time.Date(2026, 7, 21, 0, 0, 0, 0, time.UTC)
	clock := func() time.Time { return fixedTime }
	rub := Money{amount: 100, currency: Currency("RUB")}

	t.Run("success", func(t *testing.T) {
		got, err := newEntry(entryID, accountID, transferID, rub, EntryDirectionDebit, clock)
		require.NoError(t, err)
		require.Equal(t, Entry{
			id:         entryID,
			accountID:  accountID,
			transferID: transferID,
			amount:     rub,
			direction:  EntryDirectionDebit,
			createdAt:  clock(),
		}, got)
	})

	tests := []struct {
		name       string
		id         uuid.UUID
		accountID  uuid.UUID
		transferID uuid.UUID
		amount     Money
		wantErr    error
	}{
		{"invalid entry id", uuid.Nil, accountID, transferID, rub, errs.ErrInvalidEntryID},
		{"invalid account id", entryID, uuid.Nil, transferID, rub, errs.ErrInvalidAccountID},
		{"invalid transfer id", entryID, accountID, uuid.Nil, rub, errs.ErrInvalidTransferID},
		{"zero amount", entryID, accountID, transferID, Money{}, errs.ErrZeroAmount},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newEntry(tt.id, tt.accountID, tt.transferID, tt.amount, EntryDirectionDebit, clock)
			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, Entry{}, got)
		})
	}
}
