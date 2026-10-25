package dto

import (
	"github.com/google/uuid"
)

type CreateRecord struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Currency string    `json:"currency"`
}

func Restore(id uuid.UUID, name, currency string) CreateRecord {
	return CreateRecord{
		ID:       id,
		Name:     name,
		Currency: currency,
	}
}

type GetRecord struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Currency string    `json:"currency"`
	Balance  int64     `json:"balance"`
	Status   string    `json:"status"`
	CratedAt string    `json:"crated_at"`
}

func RestoreGetRecord(id uuid.UUID, name, currency, status, createdAt string, balance int64) GetRecord {
	return GetRecord{
		ID:       id,
		Name:     name,
		Currency: currency,
		Balance:  balance,
		Status:   status,
		CratedAt: createdAt,
	}
}
