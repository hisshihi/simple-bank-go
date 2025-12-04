package db

import "gorm.io/gorm"

// Entry запись о списании или зачислении средств.
type Entry struct {
	gorm.Model
	AccountID uint    `gorm:"not null;index:idx_account_entries"`
	Account   Account `gorm:"foreignKey:AccountID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Amount    int64   `gorm:"not null"` // положительное = пополнение, отрицательное = списание
}

func (Entry) TableName() string {
	return "entries"
}
