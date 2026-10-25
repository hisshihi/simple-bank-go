package account

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/hisshihi/simple-bank/internal/domain"
	errs "github.com/hisshihi/simple-bank/internal/errs"
	"github.com/hisshihi/simple-bank/internal/usecase/account/dto"
)

type Repository interface {
	Create(ctx context.Context, account *domain.Account) (*domain.Account, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error)
}

type Service struct {
	accountRepo Repository
	clock       func() time.Time
	log         *slog.Logger
}

func New(accountRepo Repository, clock func() time.Time, log *slog.Logger) *Service {
	return &Service{accountRepo: accountRepo, clock: clock, log: log}
}

func (s *Service) Create(ctx context.Context, name, currency string) (dto.CreateRecord, error) {
	id, err := uuid.NewV7()
	if err != nil {
		s.log.ErrorContext(ctx, "failed to generate uuid", "err", err)
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

	s.log.InfoContext(ctx, "record_create", "record", record)

	return record, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (dto.GetRecord, error) {
	account, err := s.accountRepo.GetByID(ctx, id)
	if err != nil {
		return dto.GetRecord{}, err
	}

	record := dto.RestoreGetRecord(account.ID(), account.Owner(), string(account.Currency()), string(account.Status()), account.CreatedAt().Format("02.01.2006 15:01"), account.Balance().Amount())
	return record, nil
}
