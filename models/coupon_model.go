package models

import (
	"time"

	"gorm.io/gorm"
)

type Coupons struct {
	gorm.Model
	ID          uint      `json:"id" gorm:"unique;not null"`
	CouponCode  string    `json:"couponCode" gorm:"unique;not null"`
	Discount    uint      `json:"discount" gorm:"not null"`
	MaxDiscount uint      `json:"maxdiscount" gorm:"not null"`
	MinPurchase uint      `json:"minpurchase" gorm:"not null"`
	Status      string    `json:"status" gorm:"default:'active'"`
	Expiry      time.Time `json:"expriy" gorm:"not null"`
}
