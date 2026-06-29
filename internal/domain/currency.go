package domain

import (
	errs "github.com/hisshihi/simple-bank/internal/errors"
)

type Currency string

var (
	CurrencyTypeUSD Currency = "USD"
	CurrencyTypeRUB Currency = "RUB"
	CurrencyTypeEUR Currency = "EUR"
	CurrencyTypeYEN Currency = "YEN"
)

var validateCurrency = map[Currency]struct{}{
	CurrencyTypeUSD: {},
	CurrencyTypeRUB: {},
	CurrencyTypeEUR: {},
	CurrencyTypeYEN: {},
}

func NewCurrency(s string) (Currency, error) {
	c := Currency(s)
	if _, ok := validateCurrency[c]; !ok {
		return "", errs.ErrCurrencyIsNotSupported
	}
	return c, nil
}
