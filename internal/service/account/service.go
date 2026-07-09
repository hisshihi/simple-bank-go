package account

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hisshihi/simple-bank/internal/domain"
	errs "github.com/hisshihi/simple-bank/internal/errors"
	"github.com/hisshihi/simple-bank/internal/logger"
	"github.com/hisshihi/simple-bank/internal/service/account/dto"
)

type AccountRepository interface {
	Create(ctx context.Context, account *domain.Account) (*domain.Account, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error)
}

type Service struct {
	accountRepo AccountRepository
	clock       func() time.Time
}

func New(accountRepo AccountRepository, clock func() time.Time) *Service {
	return &Service{accountRepo: accountRepo, clock: clock}
}

func (s *Service) Create(ctx context.Context, name, currency string) (dto.CreateRecord, error) {
	log := logger.FromContext(ctx)

	id, err := uuid.NewV7()
	if err != nil {
		return dto.CreateRecord{}, errs.ErrGenerateUUID
	}

	newCurrency, err := domain.NewCurrency(currency)
	if err != nil {
		return dto.CreateRecord{}, err
	}

	newAccount, err := domain.NewAccount(id, name, newCurrency, s.clock)
	if err != nil {
		return dto.CreateRecord{}, err
	}

	account, err := s.accountRepo.Create(ctx, newAccount)
	if err != nil {
		return dto.CreateRecord{}, err
	}

	record := dto.Restore(account.ID(), account.Owner(), string(account.Currency()))

	log.Info("record_create", "record", record)

	return record, nil
}
