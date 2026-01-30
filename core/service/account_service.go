// Package service сервисный слой для работы с аккаунтами
package service

import (
	"context"

	"github.com/hisshihi/simple-bank/core/db"
)

type AccountRepository interface {
	Create(ctx context.Context, owner string, balance float64, currency string) error
	FindByID(ctx context.Context, id uint) (*db.Account, error)
	FindAllAccounts(ctx context.Context) ([]db.Account, error)
	UpdateBalance(ctx context.Context, id uint, balance float64) error
	Delete(ctx context.Context, id uint) error
}

type AccountService struct {
	repo AccountRepository
}

func NewAccountService(repo AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}
