// Package repo содержит репозитории для работы с данными
package repo

import (
	"context"
	"fmt"

	"github.com/hisshihi/simple-bank/core/db"
	"gorm.io/gorm"
)

type accountRepo struct {
	db *gorm.DB
}

func NewAccountRepo(db *gorm.DB) *accountRepo {
	return &accountRepo{db: db}
}

func (r *accountRepo) Create(ctx context.Context, owner string, balance float64, currency string) error {
	account := &db.Account{
		Owner:    owner,
		Balance:  balance,
		Currency: db.Currency(currency),
	}
	return gorm.G[db.Account](r.db).Create(ctx, account)
}

func (r *accountRepo) FindByID(ctx context.Context, id uint) (*db.Account, error) {
	result := gorm.WithResult()
	account, err := gorm.G[db.Account](r.db, result).Where("id = ?", id).First(ctx)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, db.ErrRecordNotFound
		}
		return nil, fmt.Errorf("%s: %w", db.ErrInQuery.Error(), err)
	}

	return &account, nil
}
