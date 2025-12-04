package db

import "gorm.io/gorm"

// Transaction представляет собой перевод между двумя аккаунтами.
type Transaction struct {
	gorm.Model
	FromAccountID uint    `gorm:"not null;index:idx_from_account_transactions;index:idx_composite"`
	FromAccount   Account `gorm:"foreignKey:FromAccountID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ToAccountID   uint    `gorm:"not null;index:idx_to_account_transactions;index:idx_composite"`
	ToAccount     Account `gorm:"foreignKey:ToAccountID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Amount        int64   `gorm:"not null"`
}

func (Transaction) TableName() string {
	return "transactions"
}
