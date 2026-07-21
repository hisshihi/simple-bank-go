package transfer

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hisshihi/simple-bank/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo interface {
	Create(ctx context.Context, transfer *domain.Transfer) error
}

type EntryRepo interface {
	Create(ctx context.Context, entry *domain.Entry) error
}

type AccountRepo interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error)
	UpdateBalance(ctx context.Context, id uuid.UUID, balance int64) error
}

type Service struct {
	pool        *pgxpool.Pool
	repo        Repo
	entryRepo   EntryRepo
	clock       func() time.Time
	accountRepo AccountRepo
}

func NewService(
	pool *pgxpool.Pool,
	repo Repo,
	entryRepo EntryRepo,
	clock func() time.Time,
	accountRepo AccountRepo,
) *Service {
	return &Service{
		pool:        pool,
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

	from, err := s.getAccountID(ctx, fromID)
	if err != nil {
		return err
	}
	to, err := s.getAccountID(ctx, toID)
	if err != nil {
		return err
	}

	transferID, transferIDErr := newUUID()
	if transferIDErr != nil {
		return fmt.Errorf("create transfer id: %w", transferIDErr)
	}

	//	1. открываем транзакцию
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	err = from.Withdraw(money)
	if err != nil {
		return err
	}

	err = to.Deposit(money)
	if err != nil {
		return err
	}

	transfer, err := domain.NewTransfer(transferID, fromID, toID, money, s.clock)
	if err != nil {
		return fmt.Errorf("create transfer: %w", err)
	}

	err = s.accountRepo.UpdateBalance(ctx, fromID, amount)
	if err != nil {
		return fmt.Errorf("update from account balance: %w", err)
	}

	err = s.accountRepo.UpdateBalance(ctx, toID, amount)
	if err != nil {
		return fmt.Errorf("update to account balance: %w", err)
	}

	createTransferErr := s.repo.Create(ctx, transfer)
	if createTransferErr != nil {
		return fmt.Errorf("create transfer: %w", createTransferErr)
	}

	for _, entry := range transfer.Entries() {
		if err = s.entryRepo.Create(ctx, &entry); err != nil {
			return fmt.Errorf("create entry: %w", err)
		}
	}

	// 2. rollback по умолчанию после Commit это no-op
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func (s *Service) getAccountID(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	return s.accountRepo.GetByID(ctx, id)
}

func newUUID() (uuid.UUID, error) {
	return uuid.NewV7()
}
