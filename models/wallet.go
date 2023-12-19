package models

import (
	"gorm.io/gorm"
)

type Wallet struct {
	gorm.Model
	ID           uint          `json:"id" gorm:"unique;not null"`
	UserID       uint          `json:"userID" gorm:"index;foreignKey:userID;constraint:OnDelete:CASCADE"`
	Pin          string        `json:"pin" gorm:"not null"`
	Balance      float64       `json:"balance" gorm:"not null"`
	Transactions []Transaction `json:"transaction" gorm:"foreignKey:WalletID;constraint:OnDelete:CASCADE"`
}

type Transaction struct {
	gorm.Model
	ID          uint   `json:"id" gorm:"unique;not null"`
	WalletID    uint   `json:"walletID" gorm:"index;foreignKey:walletID;constraint:OnDelete:CASCADE"`
	Amount      int    `json:"amount" gorm:"not null"`
	Type        string `json:"type" gorm:"not null"`
	Description string `json:"description" gorm:"not null"`
}

func (w *Wallet) BeforeCreate(tx *gorm.DB) (err error) {
	w.Balance = 0
	for _, val := range w.Transactions {
		if val.Type == "credit" {
			w.Balance += float64(val.Amount)
		} else if val.Type == "debit" {
			w.Balance -= float64(val.Amount)
		}
	}
	return nil
}

func (w *Wallet) BeforeUpdate(tx *gorm.DB) (err error) {
	return w.BeforeCreate(tx)
}
