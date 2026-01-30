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

func (r *AccountRepo) FindAllAccounts(ctx context.Context) ([]db.Account, error) {
	result := gorm.WithResult()
	accounts, err := gorm.G[db.Account](r.db, result).Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", db.ErrInQuery.Error(), err)
	}

	if len(accounts) == 0 {
		return []db.Account{}, nil
	}

	return accounts, nil
}

func (r *AccountRepo) UpdateBalance(ctx context.Context, id uint, balance float64) error {
	rowsUpdated, err := gorm.G[db.Account](r.db).Where("id = ?", id).Update(ctx, "balance", balance)
	if err != nil {
		return fmt.Errorf("%s: %w", db.ErrInQuery.Error(), err)
	}

	if rowsUpdated == 0 {
		return db.ErrRecordNotFound
	}

	return nil
}

func (r *AccountRepo) Delete(ctx context.Context, id uint) error {
	rowsDeleted, err := gorm.G[db.Account](r.db).Where("id = ?", id).Delete(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", db.ErrInQuery.Error(), err)
	}

	if rowsDeleted == 0 {
		return db.ErrRecordNotFound
	}

	return nil
}
