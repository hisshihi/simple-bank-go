package domain

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	errs "github.com/hisshihi/simple-bank/internal/errors"
	"github.com/stretchr/testify/require"
)

func TestNewAccount(t *testing.T) {
	newUUID, _ := uuid.NewV7()
	currency, _ := NewCurrency("RUB")

	fixedTime := time.Date(2026, 6, 2, 0, 0, 0, 0, time.UTC)
	clock := func() time.Time { return fixedTime }

	type args struct {
		id       uuid.UUID
		owner    string
		currency Currency
		clock    Clock
	}
	tests := []struct {
		name    string
		args    args
		want    *Account
		wantErr bool
		err     error
	}{
		{
			name: "NewAccount",
			args: args{
				id:       newUUID,
				owner:    "hiss",
				currency: currency,
				clock:    clock,
			},
			want: &Account{
				id:        newUUID,
				owner:     "hiss",
				balance:   Money{amount: 0, currency: currency},
				currency:  currency,
				createdAt: clock(),
				status:    AccountStatusActive,
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "invalid account id",
			args: args{
				id:       uuid.Nil,
				owner:    "hiss",
				currency: currency,
				clock:    clock,
			},
			want:    nil,
			wantErr: true,
			err:     errs.ErrInvalidAccountID,
		},
		{
			name: "err owner",
			args: args{
				id:       newUUID,
				owner:    "",
				currency: currency,
				clock:    clock,
			},
			want:    nil,
			wantErr: true,
			err:     errs.ErrOwnerRequired,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAccount(tt.args.id, tt.args.owner, tt.args.currency, tt.args.clock)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				require.ErrorIs(t, err, tt.err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAccount() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccount_Withdraw(t *testing.T) {
	rub := Currency("RUB")

	tests := []struct {
		name    string
		balance Money
		status  AccountStatus
		amount  Money
		wantErr error
		want    int64
	}{
		{
			name:    "success",
			balance: Money{amount: 100, currency: rub},
			status:  AccountStatusActive,
			amount:  Money{amount: 50, currency: rub},
			want:    50,
		},
		{
			name:    "account not active",
			status:  AccountStatusBlocked,
			wantErr: errs.ErrAccountNotActive,
		},
		{
			name:    "balance is zero",
			status:  AccountStatusActive,
			balance: Money{amount: 0, currency: rub},
			wantErr: errs.ErrBalanceIsLessThanZero,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				id:       NewUUID(t),
				owner:    "hiss",
				balance:  tt.balance,
				currency: tt.balance.Currency(),
				status:   tt.status,
			}

			err := a.Withdraw(tt.amount)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Equal(t, tt.balance, a.Balance())
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, a.balance.Amount())
		})
	}
}

func NewUUID(t *testing.T) uuid.UUID {
	t.Helper()
	newUUID, err := uuid.NewV7()
	require.NoError(t, err)
	return newUUID
}

func TestAccount_Deposit(t *testing.T) {
	rub := Currency("RUB")

	tests := []struct {
		name    string
		amount  Money
		status  AccountStatus
		balance Money
		wantErr error
		want    int64
	}{
		{
			name:    "success",
			amount:  Money{amount: 100, currency: rub},
			status:  AccountStatusActive,
			balance: Money{amount: 100, currency: rub},
			want:    200,
		},
		{
			name:    "account not active",
			status:  AccountStatusBlocked,
			wantErr: errs.ErrAccountNotActive,
		},
		{
			name:    "amount is zero",
			amount:  Money{amount: 0, currency: rub},
			wantErr: errs.ErrZeroAmount,
			status:  AccountStatusActive,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				id:       NewUUID(t),
				owner:    "hiss",
				balance:  tt.balance,
				currency: tt.balance.Currency(),
				status:   tt.status,
			}

			err := a.Deposit(tt.amount)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Equal(t, tt.balance, a.Balance())
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, a.balance.Amount())
		})
	}
}
