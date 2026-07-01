package account

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hisshihi/simple-bank/internal/domain"
	errs "github.com/hisshihi/simple-bank/internal/errors"
	"github.com/hisshihi/simple-bank/internal/repository/account/entity"
	"github.com/jackc/pgx/v5"
)

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
