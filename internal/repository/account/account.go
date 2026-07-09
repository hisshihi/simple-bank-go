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
	WHERE id = $1
`, id).Scan(&record.AccountID, &record.Owner, &record.Currency, &record.Balance, &record.CreatedAt, &record.Status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrAccountNotFound
		}
		return nil, err
	}

	return record.ToDomain()
}
