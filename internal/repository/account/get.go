package account

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/hisshihi/simple-bank/internal/domain"
	"github.com/jackc/pgx/v5"
)

func (r *Repo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	var (
		accountID uuid.UUID
		owner     string
		balance   int64
		currency  string
		createdAt time.Time
		status    string
	)

	err := r.db.QueryRow(ctx, `
	SELECT id, owner, currency, balance, created_at, status
	FROM accounts
	WHERE id = $1
`, id).Scan(&accountID, &owner, &currency, &balance, &createdAt, &status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	curr, err := domain.NewCurrency(currency)
	if err != nil {
		return nil, err
	}

	bal, err := domain.NewMoney(balance, curr)
	if err != nil {
		return nil, err
	}

	return domain.RestoreAccount(accountID, owner, bal, curr, createdAt, domain.AccountStatus(status)), nil
}
