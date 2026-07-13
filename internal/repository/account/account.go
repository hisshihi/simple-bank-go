package account

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hisshihi/simple-bank/internal/domain"
	errs "github.com/hisshihi/simple-bank/internal/errors"
	"github.com/hisshihi/simple-bank/internal/repository/account/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DB interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type Repo struct {
	db DB
}

func New(db DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Create(ctx context.Context, account *domain.Account) (*domain.Account, error) {
	_, err := r.db.Exec(ctx, `INSERT INTO accounts (id, owner, balance, currency, created_at, status)
	VALUES ($1, $2, $3, $4, $5, $6)
 	ON CONFLICT (id) DO NOTHING`,
		account.ID(),
		account.Owner(),
		account.Balance().Amount(),
		account.Currency(),
		account.CreatedAt(),
		account.Status())
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (r *Repo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	record := entity.GetAccountRecord{}

	err := r.db.QueryRow(ctx, `
	SELECT id, owner, currency, balance, created_at, status
	FROM accounts
	WHERE id = $1 AND deleted_at IS NULL
`, id).Scan(&record.AccountID, &record.Owner, &record.Currency, &record.Balance, &record.CreatedAt, &record.Status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrAccountNotFound
		}
		return nil, err
	}

	return record.ToDomain()
}

func (r *Repo) UpdateBalance(ctx context.Context, id uuid.UUID, balance int64) error {
	rowsUpdated, err := r.db.Exec(ctx, `
	UPDATE accounts
	SET balance = $2
	WHERE id = $1 AND deleted_at IS NULL
`, id, balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errs.ErrAccountNotFound
		}
		return err
	}

	if rowsUpdated.RowsAffected() == 0 {
		return errs.ErrAccountNotFound
	}

	return nil
}

func (r *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `UPDATE accounts SET status = $2, deleted_at = now() WHERE id = $1 AND status != $2`, id, domain.AccountStatusDeleted)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errs.ErrAccountNotFound
	}

	return nil
}
