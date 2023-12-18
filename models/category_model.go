package models

import (
	"time"
)

type Categories struct {
	ID            uint   `json:"id" gorm:"unique;not null"`
	CategoryName  string `json:"category" gorm:"unique;not null"`
	Status        string `json:"status" gorm:"default:'listed'"`
	Image         string `json:"category_image"`
	CreatedAt     time.Time
	CategoryOffer CategoryOffer `json:"category_offer" gorm:"foreignKey:CategoryID"`
}

type CategoryOffer struct {
	ID         uint      `json:"id" gorm:"unique;not null"`
	CategoryID uint      `json:"categoryID" gorm:"index;not null"`
	Discount   uint      `json:"discount" gorm:"default:0"`
	ExpiryAt   time.Time `json:"expiry_at"`
	Status     string    `json:"status" gorm:"default:'active'"`
	CreatedAt  time.Time
}
