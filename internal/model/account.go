package model

import (
	"time"

	"github.com/google/uuid"
	errs "github.com/hisshihi/simple-bank/internal/errors"
)

type CurrencyType string

var (
	CurrencyTypeUSD CurrencyType = "USD"
	CurrencyTypeRUB CurrencyType = "RUB"
	CurrencyTypeEUR CurrencyType = "EUR"
	CurrencyTypeYEN CurrencyType = "YEN"
)

type Account struct {
	id        uuid.UUID    `db:"id"`
	owner     string       `db:"owner"`
	balance   float64      `db:"balance"`
	currency  CurrencyType `db:"currency"`
	createdAt time.Time    `db:"created_at"`
	updatedAt time.Time    `db:"updated_at"`
}

func NewAccount(id uuid.UUID, owner string, balance float64, accCurrency string) (*Account, error) {
	if id == uuid.Nil {
		return nil, errs.ErrInvalidAccountID
	}

	if balance < 0 {
		return nil, errs.ErrBalanceIsLessThanZero
	}

	currency := CurrencyType(accCurrency)
	switch currency {
	case CurrencyTypeEUR, CurrencyTypeUSD, CurrencyTypeYEN, CurrencyTypeRUB:
		account := &Account{
			id:       id,
			owner:    owner,
			balance:  balance,
			currency: currency,
		}
		return account, nil
	default:
		return nil, errs.ErrCurrencyIsNotSupported
	}
}
