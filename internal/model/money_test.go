package model

import (
	"reflect"
	"testing"

	errs "github.com/hisshihi/simple-bank/internal/errors"
	"github.com/stretchr/testify/require"
)

func TestNewMoney(t *testing.T) {
	type args struct {
		amount   int64
		currency Currency
	}
	tests := []struct {
		name    string
		args    args
		want    Money
		wantErr bool
		err     error
	}{
		{
			name: "newMoney",
			args: args{
				amount:   0,
				currency: CurrencyTypeRUB,
			},
			want:    Money{amount: 0, currency: CurrencyTypeRUB},
			wantErr: false,
			err:     nil,
		},
		{
			name: "balance is less than zero",
			args: args{
				amount:   -20,
				currency: CurrencyTypeRUB,
			},
			want:    Money{},
			wantErr: true,
			err:     errs.ErrBalanceIsLessThanZero,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMoney(tt.args.amount, tt.args.currency)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMoney() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				require.ErrorIs(t, err, tt.err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMoney() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoney_Add(t *testing.T) {
	type fields struct {
		amount   int64
		currency Currency
	}
	type args struct {
		other Money
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Money
		wantErr bool
		err     error
	}{
		{
			name: "success",
			fields: fields{
				amount:   10,
				currency: CurrencyTypeRUB,
			},
			args: args{
				other: Money{amount: 10, currency: CurrencyTypeRUB},
			},
			want:    Money{amount: 20, currency: CurrencyTypeRUB},
			wantErr: false,
			err:     nil,
		},
		{
			name: "invalid currency type",
			fields: fields{
				amount:   10,
				currency: CurrencyTypeEUR,
			},
			args: args{
				other: Money{amount: 10, currency: CurrencyTypeRUB},
			},
			want:    Money{},
			wantErr: true,
			err:     errs.ErrCurrencyMismatch,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Money{
				amount:   tt.fields.amount,
				currency: tt.fields.currency,
			}
			got, err := m.Add(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				require.ErrorIs(t, err, tt.err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Add() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoney_Subtrack(t *testing.T) {
	type fields struct {
		amount   int64
		currency Currency
	}
	type args struct {
		other Money
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Money
		wantErr bool
		err     error
	}{
		{
			name: "success",
			fields: fields{
				amount:   10,
				currency: CurrencyTypeRUB,
			},
			args:    args{other: Money{amount: 10, currency: CurrencyTypeRUB}},
			want:    Money{amount: 0, currency: CurrencyTypeRUB},
			wantErr: false,
			err:     nil,
		},
		{
			name: "invalid currency type",
			fields: fields{
				amount:   10,
				currency: CurrencyTypeRUB,
			},
			args:    args{other: Money{amount: 10, currency: CurrencyTypeEUR}},
			want:    Money{},
			wantErr: true,
			err:     errs.ErrCurrencyMismatch,
		},
		{
			name:    "insufficient funds",
			fields:  fields{amount: 5, currency: CurrencyTypeRUB},
			args:    args{other: Money{amount: 10, currency: CurrencyTypeRUB}},
			want:    Money{},
			wantErr: true,
			err:     errs.ErrInsufficientFunds,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Money{
				amount:   tt.fields.amount,
				currency: tt.fields.currency,
			}
			got, err := m.Subtrack(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Subtrack() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				require.ErrorIs(t, err, tt.err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Subtrack() got = %v, want %v", got, tt.want)
			}
		})
	}
}
