package transfer

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

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

func (r *Repo) Create(ctx context.Context, transfer *domain.Transfer) error {
	_, err := r.db(ctx).Exec(ctx, `
	INSERT INTO transfers (id, from_account_id, to_account_id, amount, currency, created_at)
	VALUES ($1, $2, $3, $4, $5, $6)
`,
		transfer.ID(),
		transfer.FromAccountID(),
		transfer.ToAccountID(),
		transfer.Amount().Amount(),
		string(transfer.Amount().Currency()),
		transfer.CreatedAt(),
	)
	if err != nil {
		return fmt.Errorf("insert transfer: %w", err)
	}

	// entries - часть того же агрегата
	batch := &pgx.Batch{}
	for _, entry := range transfer.Entries() {
		batch.Queue(`
	INSERT INTO entries (id, account_id, transfer_id, amount, currency, created_at, direction)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
`,
			entry.ID(),
			entry.AccountID(),
			transfer.ID(),
			entry.Amount().Amount(),
			string(entry.Amount().Currency()),
			entry.CreatedAt(),
			string(entry.Direction()),
		)
	}

	br := r.db(ctx).SendBatch(ctx, batch)
	defer br.Close()

	for range transfer.Entries() {
		if _, execErr := br.Exec(); execErr != nil {
			return fmt.Errorf("insert entries: %w", execErr)
		}
	}

	return nil
}
