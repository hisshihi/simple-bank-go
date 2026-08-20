package entry

import (
	"context"

	"github.com/hisshihi/simple-bank/internal/domain"
	"github.com/hisshihi/simple-bank/internal/repository"
)

type Repo struct {
	pool repository.Querier
}

func (r *Repo) db(ctx context.Context) repository.Querier {
	if tx := repository.ExtractTx(ctx); tx != nil {
		return tx
	}
	return r.pool
}

func New(db repository.Querier) *Repo {
	return &Repo{pool: db}
}

func (r *Repo) Create(ctx context.Context, entry *domain.Entry) error {
	_, err := r.db(ctx).Exec(ctx, `INSERT INTO entries (id, account_id, transfer_id, amount, currency, direction, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		entry.ID(), entry.AccountID(), entry.Amount().Amount(), entry.CreatedAt())
	if err != nil {
		return err
	}
	return nil
}
