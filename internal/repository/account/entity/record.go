package entity

import (
	"time"

	"github.com/google/uuid"
)

type GetAccountRecord struct {
	AccountID uuid.UUID
	Owner     string
	Balance   int64
	Currency  string
	CreatedAt time.Time
	Status    string
}
