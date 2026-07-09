package dto

import "github.com/google/uuid"

type CreateRecord struct {
	ID       uuid.UUID
	Name     string
	Currency string
}

func Restore(id uuid.UUID, name, currency string) CreateRecord {
	return CreateRecord{
		ID:       id,
		Name:     name,
		Currency: currency,
	}
}
