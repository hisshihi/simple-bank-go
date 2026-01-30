package db

import (
	"errors"

	"gorm.io/gorm"
)

type Currency string

var (
	ErrRecordNotFound = errors.New("аккаунт не найден")
	ErrInQuery        = errors.New("ошибка при выполнении запроса к базе данных")
)

const (
	USD Currency = "USD"
	EUR Currency = "EUR"
	GBP Currency = "GBP"
	RUB Currency = "RUB"
)

func (c Currency) String() string {
	return string(c)
}

// Account аккаунт пользователя в системе.
type Account struct {
	gorm.Model
	Owner    string   `gorm:"type:text;not null"`                      // user`s username
	Balance  float64  `gorm:"precision:20;scale:2;not null;default:0"` // account balance
	Currency Currency `gorm:"type:text;not null"`                      // account currency
}

func (Account) TableName() string {
	return "accounts"
}
