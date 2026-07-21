package entry

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

func (r *Repo) Create(ctx context.Context, entry *domain.Entry) error {
	_, err := r.db.Exec(ctx, `INSERT INTO entries (id, account_id, transfer_id, amount, currency, direction, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		entry.ID(), entry.AccountID(), entry.Amount().Amount(), entry.CreatedAt())
	if err != nil {
		return err
	}
	return nil
}
