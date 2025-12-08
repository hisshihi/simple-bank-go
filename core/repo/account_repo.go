// Package repo содержит репозитории для работы с данными
package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/hisshihi/simple-bank/core/db"
	"gorm.io/gorm"
)

type AccountRepo struct {
	db *gorm.DB
}

func NewAccountRepo(db *gorm.DB) *AccountRepo {
	return &AccountRepo{db: db}
}

func (r *AccountRepo) Create(ctx context.Context, owner string, balance float64, currency string) error {
	account := &db.Account{
		Owner:    owner,
		Balance:  balance,
		Currency: db.Currency(currency),
	}
	return gorm.G[db.Account](r.db).Create(ctx, account)
}

func (r *AccountRepo) FindByID(ctx context.Context, id uint) (*db.Account, error) {
	result := gorm.WithResult()
	account, err := gorm.G[db.Account](r.db, result).Where("id = ?", id).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, db.ErrRecordNotFound
		}
		return nil, fmt.Errorf("%s: %w", db.ErrInQuery.Error(), err)
	}

	return &account, nil
}
