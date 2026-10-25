package account

import (
	"context"
	"errors"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/hisshihi/simple-bank/internal/domain"
	errs "github.com/hisshihi/simple-bank/internal/errs"
	"github.com/hisshihi/simple-bank/internal/repository/account/entity"
)

type Repo struct {
	db     trmpgx.Tr
	getter *trmpgx.CtxGetter
}

func New(db trmpgx.Tr, getter *trmpgx.CtxGetter) *Repo {
	return &Repo{db: db, getter: getter}
}

func (r *Repo) Create(ctx context.Context, account *domain.Account) (*domain.Account, error) {
	_, err := r.getter.DefaultTrOrDB(ctx, r.db).Exec(ctx, `INSERT INTO accounts (id, owner, balance, currency, created_at, status)
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

	err := r.getter.DefaultTrOrDB(ctx, r.db).QueryRow(ctx, `
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

func (r *Repo) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	record := entity.GetAccountRecord{}

	err := r.getter.DefaultTrOrDB(ctx, r.db).QueryRow(ctx, `
	SELECT id, owner, currency, balance, created_at, status
	FROM accounts
	WHERE id = $1 AND deleted_at IS NULL
	FOR UPDATE 
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
	rowsUpdated, err := r.getter.DefaultTrOrDB(ctx, r.db).Exec(ctx, `
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

func (r *Repo) Save(ctx context.Context, a *domain.Account) error {
	_, err := r.getter.DefaultTrOrDB(ctx, r.db).Exec(ctx, `
	UPDATE accounts SET balance = $2, status = $3 WHERE id = $1
`, a.ID(), a.Balance().Amount(), a.Status().String())
	return err
}

func (r *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	tag, err := r.getter.DefaultTrOrDB(ctx, r.db).Exec(ctx, `UPDATE accounts SET status = $2, deleted_at = now() WHERE id = $1 AND status != $2`, id, domain.AccountStatusDeleted)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errs.ErrAccountNotFound
	}

	return nil
}
