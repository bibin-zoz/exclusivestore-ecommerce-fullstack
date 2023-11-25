package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	ID        uint   `json:"id" gorm:"unique;not null"`
	ProductID uint   `json:"productID" gorm:"index;foreignKey:ProductID"`
	VariantID uint   `json:"variantID" gorm:"index;foreignKey:VariantID"`
	UserName  string `json:"userName" gorm:"index;foreignKey:UserName"`
	Quantity  uint   `json:"quantity" gorm:"not null"`
	Product   Products
	Variant   ProductVariants
	// Image     []Image `gorm:"foreignKey:ProductID"`
}

type GetCart struct {
	ID        uint `json:"id" `
	ProductID uint `json:"productID"`
	VariantID uint `json:"variantID"`
}
