// Package service сервисный слой для работы с аккаунтами
package service

type AccountRepository interface {
}

type AccountService struct {
	repo AccountRepository
}

func NewAccountService(repo AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}
