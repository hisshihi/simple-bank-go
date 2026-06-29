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
