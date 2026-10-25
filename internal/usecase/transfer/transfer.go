package transfer

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/hisshihi/simple-bank/internal/domain"
)

type Repo interface {
	Create(ctx context.Context, transfer *domain.Transfer) error
}

type EntryRepo interface {
	Create(ctx context.Context, entry *domain.Entry) error
}

type AccountRepo interface {
	GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.Account, error)
	UpdateBalance(ctx context.Context, id uuid.UUID, balance int64) error
	Save(ctx context.Context, a *domain.Account) error
}

type txManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type Service struct {
	txm         txManager
	repo        Repo
	entryRepo   EntryRepo
	clock       func() time.Time
	accountRepo AccountRepo
}

func NewService(
	txm txManager,
	repo Repo,
	entryRepo EntryRepo,
	clock func() time.Time,
	accountRepo AccountRepo,
) *Service {
	return &Service{
		txm:         txm,
		repo:        repo,
		entryRepo:   entryRepo,
		clock:       clock,
		accountRepo: accountRepo,
	}
}

func (s *Service) Create(ctx context.Context, fromID, toID uuid.UUID, amount int64, currency string) error {
	curr, err := domain.NewCurrency(currency)
	if err != nil {
		return err
	}
	money, err := domain.NewMoney(amount, curr)
	if err != nil {
		return err
	}

	transferID, err := newUUID()
	if err != nil {
		return fmt.Errorf("create transfer id: %w", err)
	}

	return s.txm.Do(ctx, func(ctx context.Context) error {
		from, fromErr := s.accountRepo.GetByIDForUpdate(ctx, fromID)
		if fromErr != nil {
			return fromErr
		}
		to, toErr := s.accountRepo.GetByIDForUpdate(ctx, toID)
		if toErr != nil {
			return toErr
		}

		if errWithdraw := from.Withdraw(money); errWithdraw != nil {
			return errWithdraw
		}
		if errDeposit := to.Deposit(money); errDeposit != nil {
			return errDeposit
		}

		transfer, transferErr := domain.NewTransfer(transferID, fromID, toID, money, s.clock)
		if transferErr != nil {
			return fmt.Errorf("create transfer: %w", transferErr)
		}

		// сохраняем посчитанный баланс
		if err = s.accountRepo.Save(ctx, from); err != nil {
			return fmt.Errorf("save from: %w", err)
		}
		if err = s.accountRepo.Save(ctx, to); err != nil {
			return fmt.Errorf("save to: %w", err)
		}
		if err = s.repo.Create(ctx, transfer); err != nil {
			return fmt.Errorf("create transfer: %w", err)
		}
		for _, entry := range transfer.Entries() {
			if err = s.entryRepo.Create(ctx, &entry); err != nil {
				return fmt.Errorf("create entry: %w", err)
			}
		}
		return nil
	})
}

func newUUID() (uuid.UUID, error) {
	return uuid.NewV7()
}
