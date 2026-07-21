package domain

import (
	"time"

	"github.com/google/uuid"
	errs "github.com/hisshihi/simple-bank/internal/errors"
)

type EntryDirection string

const (
	EntryDirectionDebit  EntryDirection = "debit"
	EntryDirectionCredit EntryDirection = "credit"
)

// Entry - все изменения баланса счёта.
// К примеру у аккаунта 5 списание -5, зачисление +15 и тд.
type Entry struct {
	id         uuid.UUID
	accountID  uuid.UUID
	transferID uuid.UUID
	amount     Money
	direction  EntryDirection
	createdAt  time.Time
}

func newEntry(id uuid.UUID, accountID, transferID uuid.UUID, amount Money, direction EntryDirection, createdAt Clock) (Entry, error) {
	if id == uuid.Nil {
		return Entry{}, errs.ErrInvalidEntryID
	}

	if accountID == uuid.Nil {
		return Entry{}, errs.ErrInvalidAccountID
	}

	if transferID == uuid.Nil {
		return Entry{}, errs.ErrInvalidTransferID
	}

	if amount.IsZero() {
		return Entry{}, errs.ErrZeroAmount
	}

	now := createdAt()

	return Entry{
		id:         id,
		accountID:  accountID,
		transferID: transferID,
		amount:     amount,
		direction:  direction,
		createdAt:  now,
	}, nil
}

func RestoreEntry(id, accountID, transferID uuid.UUID, amount Money, direction EntryDirection, createdAt time.Time) Entry {
	return Entry{
		id:         id,
		accountID:  accountID,
		transferID: transferID,
		amount:     amount,
		direction:  direction,
		createdAt:  createdAt,
	}
}

func (e Entry) ID() uuid.UUID {
	return e.id
}
func (e Entry) AccountID() uuid.UUID {
	return e.accountID
}
func (e Entry) TransferID() uuid.UUID {
	return e.transferID
}
func (e Entry) Amount() Money {
	return e.amount
}
func (e Entry) Direction() EntryDirection {
	return e.direction
}
func (e Entry) CreatedAt() time.Time {
	return e.createdAt
}
