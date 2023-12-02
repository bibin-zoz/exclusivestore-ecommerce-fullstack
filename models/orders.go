package models

import "gorm.io/gorm"

type Orders struct {
	gorm.Model
	ID        uint            `json:"id" gorm:"unique;not null"`
	UserID    uint            `json:"UserID" gorm:"index;foreignKey:UserID"`
	AddressID uint            `json:"AddressID" gorm:"index;foreignKey:AddressID"`
	ProductID uint            `json:"productID" gorm:"index;foreignKey:ProductID"`
	VariantID uint            `json:"variantID" gorm:"index;foreignKey:VariantID"`
	Status    string          `json:"status" gorm:"default:'pending'"`
	Payment   string          `json:"payment" gorm:"default:'cod'"`
	Quantity  uint            `json:"quantity" gorm:"not null"`
	Price     float64         `json:"price" gorm:"not null"`
	Total     float64         `json:"Total" gorm:"not null"`
	User      User            `gorm:"foreignKey:UserID"`
	Product   Products        `gorm:"foreignKey:ProductID"`
	Variant   ProductVariants `gorm:"foreignKey:VariantID"`
	Address   UserAddress     `gorm:"foreignKey:AddressID"`
}
