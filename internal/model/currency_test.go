package model

import (
	"testing"

	errs "github.com/hisshihi/simple-bank/internal/errors"
	"github.com/stretchr/testify/require"
)

func TestNewCurrency(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    Currency
		wantErr bool
		err     error
	}{
		{
			name: "newCurrency",
			args: args{
				"RUB",
			},
			want:    CurrencyTypeRUB,
			wantErr: false,
			err:     nil,
		},
		{
			name: "currency is not supported",
			args: args{
				"NIL",
			},
			want:    "",
			wantErr: true,
			err:     errs.ErrCurrencyIsNotSupported,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCurrency(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCurrency() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				require.ErrorIs(t, err, tt.err)
			}

			if got != tt.want {
				t.Errorf("NewCurrency() got = %v, want %v", got, tt.want)
			}
		})
	}
}
