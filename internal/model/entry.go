package model

import (
	"time"

	"github.com/google/uuid"
	errs "github.com/hisshihi/simple-bank/internal/errors"
)

type Entry struct {
	ID        uuid.UUID `db:"id"`
	AccountID uuid.UUID `db:"account_id"`
	Amount    int64     `db:"amount"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewEntry(id uuid.UUID, accountID uuid.UUID, amount int64) (*Entry, error) {
	if id == uuid.Nil {
		return nil, errs.ErrInvalidEntryID
	}

	if accountID == uuid.Nil {
		return nil, errs.ErrInvalidAccountID
	}

	return &Entry{
		ID:        id,
		AccountID: accountID,
		Amount:    amount,
	}, nil
}
