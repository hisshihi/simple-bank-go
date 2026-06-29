package account

import (
	"context"

	"github.com/hisshihi/simple-bank/internal/domain"
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
