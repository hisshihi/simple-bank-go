package model

import (
	"errors"
	"fmt"

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

var validCurrency = map[CurrencyType]struct{}{
	CurrencyTypeUSD: {},
	CurrencyTypeRUB: {},
	CurrencyTypeEUR: {},
	CurrencyTypeYEN: {},
}

type Account struct {
	ID       uuid.UUID    `db:"id"`
	Owner    string       `db:"owner"`
	Balance  float64      `db:"balance"`
	Currency CurrencyType `db:"currency"`
}

func NewAccount(id uuid.UUID, owner string, balance float64, currency CurrencyType) (*Account, error) {
	if id == uuid.Nil {
		return nil, errs.ErrInvalidAccountID
	}

	if balance < 0 {
		return nil, errs.ErrBalanceIsLessThanZero
	}

	validateCurrency, err := CheckCurrency(string(currency))
	if err != nil {
		return nil, errs.ErrCurrencyIsNotSupported
	}

	return &Account{
		ID:       id,
		Owner:    owner,
		Balance:  balance,
		Currency: *validateCurrency,
	}, nil
}

func CheckCurrency(currency string) (*CurrencyType, error) {
	if currency == "" {
		return nil, errors.New("currency is empty")
	}

	c := CurrencyType(currency)
	if _, ok := validCurrency[c]; !ok {
		return nil, fmt.Errorf("invalid currency %s", currency)
	}

	return &c, nil
}
