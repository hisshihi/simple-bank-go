package transfer

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// func TestService_Create(t *testing.T) {
//	fromID, toID := NewUUID(t), NewUUID(t)
//	amount := 100
//	currency := "RUB"
//	t.Run("success", func(t *testing.T) {
//		repo := mocks.NewMockRepo(t)
//
//		repo.EXPECT()
//	})
//}

func NewUUID(t *testing.T) uuid.UUID {
	t.Helper()
	newID, err := uuid.NewV7()
	require.NoError(t, err)
	return newID
}
