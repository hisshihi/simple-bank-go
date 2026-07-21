package domain

import (
	"time"

	"github.com/google/uuid"
	errs "github.com/hisshihi/simple-bank/internal/errors"
)

// Transfer - перевод между аккаунтами.
type Transfer struct {
	id            uuid.UUID
	fromAccountID uuid.UUID
	toAccountID   uuid.UUID
	amount        Money
	entries       []Entry
	createdAt     time.Time
}

func NewTransfer(id, fromAccountID, toAccountID uuid.UUID, amount Money, clock Clock) (*Transfer, error) {
	if id == uuid.Nil {
		return nil, errs.ErrInvalidTransferID
	}

	if fromAccountID == uuid.Nil {
		return nil, errs.ErrInvalidFromAccountID
	}

	if toAccountID == uuid.Nil {
		return nil, errs.ErrInvalidToAccountID
	}

	if fromAccountID == toAccountID {
		return nil, errs.ErrSameAccountTransfer
	}

	if amount.IsZero() {
		return nil, errs.ErrZeroAmount
	}

	debitID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	creditID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	debit, err := newEntry(debitID, fromAccountID, id, amount, EntryDirectionDebit, clock)
	if err != nil {
		return nil, err
	}
	credit, err := newEntry(creditID, toAccountID, id, amount, EntryDirectionCredit, clock)
	if err != nil {
		return nil, err
	}

	return &Transfer{
		id:            id,
		fromAccountID: fromAccountID,
		toAccountID:   toAccountID,
		amount:        amount,
		entries:       []Entry{debit, credit},
		createdAt:     clock(),
	}, nil
}

func RestoreTransfer(id, fromAccountID, toAccountID uuid.UUID, amount Money, entries []Entry, createdAt time.Time) *Transfer {
	return &Transfer{
		id:            id,
		fromAccountID: fromAccountID,
		toAccountID:   toAccountID,
		amount:        amount,
		entries:       entries,
		createdAt:     createdAt,
	}
}

func (t *Transfer) ID() uuid.UUID {
	return t.id
}
func (t *Transfer) FromAccountID() uuid.UUID {
	return t.fromAccountID
}
func (t *Transfer) ToAccountID() uuid.UUID {
	return t.toAccountID
}
func (t *Transfer) Amount() Money {
	return t.amount
}

// Entries - возвращает копию, снаружи слайс не мутируют
func (t *Transfer) Entries() []Entry {
	out := make([]Entry, len(t.entries))
	copy(out, t.entries)
	return out
}
func (t *Transfer) CreatedAt() time.Time {
	return t.createdAt
}
