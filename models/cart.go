package models

import (
	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	ID        uint            `json:"id" gorm:"unique;not null"`
	ProductID uint            `json:"productID" gorm:"index;foreignKey:ProductID"`
	VariantID uint            `json:"variantID" gorm:"index;foreignKey:VariantID"`
	UserName  string          `json:"userName" gorm:"index;foreignKey:UserName"`
	Quantity  uint            `json:"quantity" gorm:"not null"`
	Price     float64         `json:"price" gorm:"not null"`
	Total     float64         `json:"Total" gorm:"not null"`
	Product   Products        `gorm:"foreignKey:ProductID"`
	Variant   ProductVariants `gorm:"foreignKey:VariantID"`
	// Image     []Image `gorm:"foreignKey:ProductID"`
}

type GetCart struct {
	ID        uint `json:"id" `
	ProductID uint `json:"productID"`
	VariantID uint `json:"variantID"`
	Quantity  uint `json:"quantity"`
	Price     uint `json:"Price"`
}

func (p *Cart) BeforeSave(tx *gorm.DB) (err error) {
	// Your logic to set Total based on Variant.Price and Quantity

	p.Total = p.Price * float64(p.Quantity)
	return nil
}
