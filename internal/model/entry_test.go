package model

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	errs "github.com/hisshihi/simple-bank/internal/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewEntry(t *testing.T) {
	entryID, _ := uuid.NewV7()
	accountID, _ := uuid.NewV7()
	fixedTime := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	clock := func() time.Time { return fixedTime }

	type args struct {
		id        uuid.UUID
		accountID uuid.UUID
		amount    int64
		createdAt Clock
	}
	tests := []struct {
		name    string
		args    args
		want    *Entry
		wantErr bool
		err     error
	}{
		{
			name: "new entry",
			args: args{
				id:        entryID,
				accountID: accountID,
				amount:    0,
				createdAt: clock,
			},
			want: &Entry{
				ID:        entryID,
				AccountID: accountID,
				Amount:    0,
				CreatedAt: clock(),
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "invalid entry id",
			args: args{
				id:        uuid.Nil,
				accountID: accountID,
				amount:    0,
				createdAt: clock,
			},
			want:    nil,
			wantErr: true,
			err:     errs.ErrInvalidEntryID,
		},
		{
			name: "invalid account id",
			args: args{
				id:        entryID,
				accountID: uuid.Nil,
				amount:    0,
				createdAt: clock,
			},
			want:    nil,
			wantErr: true,
			err:     errs.ErrInvalidAccountID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEntry(tt.args.id, tt.args.accountID, tt.args.amount, tt.args.createdAt)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEntry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEntry() got = %v, want %v", got, tt.want)
			}
		})
	}
}
