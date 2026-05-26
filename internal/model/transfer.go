package model

import (
	"github.com/google/uuid"
	errs "github.com/hisshihi/simple-bank/internal/errors"
)

type Transfer struct {
	ID            uuid.UUID `db:"id"`
	FromAccountID uuid.UUID `db:"from_account_id"`
	ToAccountID   uuid.UUID `db:"to_account_id"`
	Amount        int64     `db:"amount"`
}

func NewTransfer(id, fromAccountID, toAccountID uuid.UUID, amount int64) (*Transfer, error) {
	if id == uuid.Nil {
		return nil, errs.ErrInvalidTransferID
	}

	if fromAccountID == uuid.Nil {
		return nil, errs.ErrInvalidFromAccountID
	}

	if toAccountID == uuid.Nil {
		return nil, errs.ErrInvalidToAccountID
	}

	return &Transfer{
		ID:            id,
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Amount:        amount,
	}, nil
}
