package transfer

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/hisshihi/simple-bank/internal/domain"
)

type Repo interface {
	Create(ctx context.Context, transfer *domain.Transfer) error
}

type AccountRepo interface {
	GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.Account, error)
	Save(ctx context.Context, a *domain.Account) error
}

type txManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type Service struct {
	txm         txManager
	repo        Repo
	clock       func() time.Time
	accountRepo AccountRepo
}

func NewService(
	txm txManager,
	repo Repo,
	clock func() time.Time,
	accountRepo AccountRepo,
) *Service {
	return &Service{
		txm:         txm,
		repo:        repo,
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
		from, to, lockErr := s.lockPair(ctx, fromID, toID)
		if lockErr != nil {
			return lockErr
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
		return nil
	})
}

// lockPair берёт FOR UPDATE в детерминированном порядке (по возрастанию ID),
// чтобы встречные переводы A→B и B→A не дедлочили друг друга.
func (s *Service) lockPair(ctx context.Context, fromID, toID uuid.UUID) (from, to *domain.Account, err error) {
	firstID, secondID := fromID, toID
	if bytes.Compare(fromID[:], toID[:]) > 0 {
		firstID, secondID = toID, fromID
	}

	first, err := s.accountRepo.GetByIDForUpdate(ctx, firstID)
	if err != nil {
		return nil, nil, fmt.Errorf("lock %s: %w", firstID, err)
	}
	second, err := s.accountRepo.GetByIDForUpdate(ctx, secondID)
	if err != nil {
		return nil, nil, fmt.Errorf("lock %s: %w", secondID, err)
	}

	if first.ID() == fromID {
		return first, second, nil
	}
	return second, first, nil
}

func newUUID() (uuid.UUID, error) {
	return uuid.NewV7()
}
