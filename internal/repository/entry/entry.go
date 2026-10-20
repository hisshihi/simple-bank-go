package entry

import (
	"context"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"

	"github.com/hisshihi/simple-bank/internal/domain"
)

type Repo struct {
	db     trmpgx.Tr
	getter *trmpgx.CtxGetter
}

func New(db trmpgx.Tr, getter *trmpgx.CtxGetter) *Repo {
	return &Repo{db: db, getter: getter}
}

func (r *Repo) Create(ctx context.Context, entry *domain.Entry) error {
	_, err := r.getter.DefaultTrOrDB(ctx, r.db).Exec(ctx, `INSERT INTO entries (id, account_id, transfer_id, amount, currency, direction, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		entry.ID(), entry.AccountID(), entry.TransferID(), entry.Amount().Amount(), entry.Amount().Currency(), entry.Direction(), entry.CreatedAt())
	if err != nil {
		return err
	}
	return nil
}
