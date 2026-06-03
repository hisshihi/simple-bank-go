package model

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	errs "github.com/hisshihi/simple-bank/internal/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTransfer(t *testing.T) {
	id, _ := uuid.NewV7()
	fromAccountID, _ := uuid.NewV7()
	toAccountID, _ := uuid.NewV7()
	fixedTime := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	clock := func() time.Time { return fixedTime }

	type args struct {
		id            uuid.UUID
		fromAccountID uuid.UUID
		toAccountID   uuid.UUID
		amount        int64
		createdAt     Clock
	}
	tests := []struct {
		name    string
		args    args
		want    *Transfer
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				id:            id,
				fromAccountID: fromAccountID,
				toAccountID:   toAccountID,
				amount:        0,
				createdAt:     clock,
			},
			want: &Transfer{
				ID:            id,
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        0,
				CreatedAt:     clock(),
			},
			wantErr: assert.NoError,
		},
		{
			name: "error invalid id",
			args: args{
				id:            uuid.Nil,
				fromAccountID: fromAccountID,
				toAccountID:   toAccountID,
				amount:        0,
				createdAt:     clock,
			},
			want: nil,
			wantErr: func(tt assert.TestingT, e error, i ...interface{}) bool {
				return assert.ErrorIs(tt, e, errs.ErrInvalidTransferID)
			},
		},
		{
			name: "error from account id",
			args: args{
				id:            id,
				fromAccountID: uuid.Nil,
				toAccountID:   toAccountID,
				amount:        0,
				createdAt:     clock,
			},
			want: nil,
			wantErr: func(tt assert.TestingT, e error, i ...interface{}) bool {
				return assert.ErrorIs(tt, e, errs.ErrInvalidFromAccountID)
			},
		},
		{
			name: "error invalid id",
			args: args{
				id:            id,
				fromAccountID: fromAccountID,
				toAccountID:   uuid.Nil,
				amount:        0,
				createdAt:     clock,
			},
			want: nil,
			wantErr: func(tt assert.TestingT, e error, i ...interface{}) bool {
				return assert.ErrorIs(tt, e, errs.ErrInvalidToAccountID)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTransfer(tt.args.id, tt.args.fromAccountID, tt.args.toAccountID, tt.args.amount, tt.args.createdAt)
			if !tt.wantErr(t, err, fmt.Sprintf("NewTransfer(%v, %v, %v, %v, %v)", tt.args.id, tt.args.fromAccountID, tt.args.toAccountID, tt.args.amount, tt.args.createdAt)) {
				return
			}
			require.Equalf(t, tt.want, got, "NewTransfer(%v, %v, %v, %v, %v)", tt.args.id, tt.args.fromAccountID, tt.args.toAccountID, tt.args.amount, tt.args.createdAt)
		})
	}
}
