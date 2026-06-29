package domain

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
		currency:  currency,
		status:    AccountStatusActive,
	}, nil
}

func RestoreAccount(
	id uuid.UUID,
	owner string,
	balance Money,
	currency Currency,
	createdAt time.Time,
	status AccountStatus,
) *Account {
	return &Account{
		id:        id,
		owner:     owner,
		balance:   balance,
		currency:  currency,
		createdAt: createdAt,
		status:    status,
	}
}

func (a *Account) ID() uuid.UUID {
	return a.id
}
func (a *Account) Owner() string {
	return a.owner
}
func (a *Account) Currency() Currency {
	return a.currency
}
func (a *Account) Balance() Money {
	return a.balance
}
func (a *Account) CreatedAt() time.Time {
	return a.createdAt
}
func (a *Account) Status() AccountStatus {
	return a.status
}
