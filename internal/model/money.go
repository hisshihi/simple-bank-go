package model

import (
	errs "github.com/hisshihi/simple-bank/internal/errors"
)

type Money struct {
	amount   int64
	currency Currency
}

func NewMoney(amount int64, currency Currency) (Money, error) {
	if amount < 0 {
		return Money{}, errs.ErrBalanceIsLessThanZero
	}
	return Money{amount: amount, currency: currency}, nil
}

// Add возвращает новый Money — сумму m и other.
// Валюты должны совпадать, иначе вернёт ErrCurrencyMismatch.
func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, errs.ErrCurrencyMismatch
	}

	return Money{amount: m.amount + other.amount, currency: m.currency}, nil
}

// Subtrack возвращает новый Money - вычитание m и other.
// Валюты должны совпадать, amount Money должен быть больше other Money
func (m Money) Subtrack(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, errs.ErrCurrencyMismatch
	}
	if m.amount < other.amount {
		return Money{}, errs.ErrInsufficientFunds
	}
	return Money{amount: m.amount - other.amount, currency: m.currency}, nil
}

func (m Money) Equals(other Money) bool {
	return m.amount == other.amount && m.currency == other.currency
}
func (m Money) IsZero() bool {
	return m.amount == 0
}
func (m Money) Amount() int64 {
	return m.amount
}
func (m Money) Currency() Currency {
	return m.currency
}
