package model

import (
	"time"

	"github.com/google/uuid"
	errs "github.com/hisshihi/simple-bank/internal/errors"
)

type AccountStatus string

type Clock func() time.Time

const (
	AccountStatusActive  AccountStatus = "active"
	AccountStatusBlocked AccountStatus = "blocked"
)

type Account struct {
	id        uuid.UUID     `db:"id"`
	owner     string        `db:"owner"`
	balance   Money         `db:"balance"`
	currency  Currency      `db:"currency"`
	createdAt time.Time     `db:"created_at"`
	status    AccountStatus `db:"status"`
}

func NewAccount(id uuid.UUID, owner string, currency Currency, clock Clock) (*Account, error) {
	if id == uuid.Nil {
		return nil, errs.ErrInvalidAccountID
	}

	if owner == "" {
		return nil, errs.ErrOwnerRequired
	}

	return &Account{
		id:        id,
		owner:     owner,
		balance:   Money{amount: 0, currency: currency},
		createdAt: clock(),
		status:    AccountStatusActive,
	}, nil
}
