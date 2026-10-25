package transfer

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func NewUUID(t *testing.T) uuid.UUID {
	t.Helper()
	newID, err := uuid.NewV7()
	require.NoError(t, err)
	return newID
}

type noRetry struct{}

func (r *noRetry) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

// TODO: реализовать, когда будет время
func TestService_Create(t *testing.T) {
	//from := NewUUID(t)
	//to := NewUUID(t)
	//amount := int64(10_000)
	//currency := "RUB"
	//
	//newCurr, err := domain.NewCurrency(currency)
	//require.NoError(t, err)
	//
	//newMoney, err := domain.NewMoney(amount, newCurr)
	//require.NoError(t, err)
	//
	//svc := NewService()
}
