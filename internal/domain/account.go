package domain

import (
	"time"

	"github.com/google/uuid"

	errs "github.com/hisshihi/simple-bank/internal/errors"
)

type AccountStatus string

func (s AccountStatus) String() string {
	return string(s)
}

type Clock func() time.Time

const (
	AccountStatusActive  AccountStatus = "active"
	AccountStatusBlocked AccountStatus = "blocked"
	AccountStatusDeleted AccountStatus = "deleted"
)

type Account struct {
	id        uuid.UUID
	owner     string
	balance   Money
	currency  Currency
	createdAt time.Time
	status    AccountStatus
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

// Withdraw - снимает сумму с баланса
func (a *Account) Withdraw(amount Money) error {
	if a.Status() != AccountStatusActive {
		return errs.ErrAccountNotActive
	}
	if a.Balance().IsZero() || amount.IsZero() {
		return errs.ErrBalanceIsLessThanZero
	}
	newBalance, err := a.Balance().Subtrack(amount)
	if err != nil {
		return err
	}
	a.balance = newBalance
	return nil
}

// Deposit - начисляет сумму на баланс
func (a *Account) Deposit(amount Money) error {
	if a.Status() != AccountStatusActive {
		return errs.ErrAccountNotActive
	}
	if amount.IsZero() {
		return errs.ErrZeroAmount
	}
	newBalance, err := a.balance.Add(amount)
	if err != nil {
		return err
	}
	a.balance = newBalance
	return nil
}
