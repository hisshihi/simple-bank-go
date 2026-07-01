package entity

import "github.com/hisshihi/simple-bank/internal/domain"

func (r GetAccountRecord) ToDomain() (*domain.Account, error) {
	curr, err := domain.NewCurrency(r.Currency)
	if err != nil {
		return nil, err
	}
	bal, err := domain.NewMoney(r.Balance, curr)
	if err != nil {
		return nil, err
	}
	return domain.RestoreAccount(
		r.AccountID, r.Owner, bal, curr, r.CreatedAt, domain.AccountStatus(r.Status),
	), nil
}
