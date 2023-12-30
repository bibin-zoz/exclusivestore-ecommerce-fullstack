package models

import (
	"time"

	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	UserID       uint           `json:"userID" gorm:"index;unique;not null;foreignKey:userID;constraint:OnDelete:CASCADE"`
	Total        float64        `json:"Total" gorm:"not null"`
	CartProducts []CartProducts `gorm:"foreignKey:CartID;"` // Update foreignKey to CartID

	// Image     []Image `gorm:"foreignKey:ProductID"`
}

type CartProducts struct {
	gorm.Model
	// CartID    uint            `json:"cartID" gorm:"index;not null;"` // Add autoIncrement tag
	CartID    uint            `json:"cartID" gorm:"index;not null;foreignKey:CartID;constraint:OnDelete:CASCADE"`
	ProductID uint            `json:"productID" gorm:"index;foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	VariantID uint            `json:"variantID" gorm:"index;foreignKey:VariantID;constraint:OnDelete:CASCADE"`
	Quantity  uint            `json:"quantity" gorm:"not null"`
	Price     float64         `json:"price" gorm:"not null"`
	Total     float64         `json:"Total" gorm:"not null"`
	Product   Products        `gorm:"foreignKey:ProductID;"`
	Variant   ProductVariants `gorm:"foreignKey:VariantID"`

	CreatedAt time.Time
}

type GetCart struct {
	ID        uint `json:"id" `
	ProductID uint `json:"productID"`
	VariantID uint `json:"variantID"`
	Quantity  uint `json:"quantity"`
	Price     uint `json:"Price"`
}

func (p *CartProducts) BeforeSave(tx *gorm.DB) (err error) {

	p.Total = p.Price * float64(p.Quantity)
	return nil
}
func (c *Cart) CalculateTotal() {
	var total float64

	for _, cartProduct := range c.CartProducts {
		total += cartProduct.Total
	}

	c.Total = total
}
func (c *Cart) BeforeSave(tx *gorm.DB) (err error) {
	c.CalculateTotal()
	return nil
}
