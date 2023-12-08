package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Orders struct {
	gorm.Model
	ID              uint            `json:"id" gorm:"unique;not null"`
	UserID          uint            `json:"UserID" gorm:"index;foreignKey:UserID"`
	AddressID       uint            `json:"AddressID" gorm:"index;foreignKey:AddressID"`
	Status          string          `json:"status" gorm:"default:'pending'"`
	Payment         string          `json:"payment" gorm:"default:'cod'"`
	Total           float64         `json:"Total" gorm:"not null"`
	User            User            `gorm:"foreignKey:UserID"`
	Address         UserAddress     `gorm:"foreignKey:AddressID"`
	OrderedProducts []OrderProducts `json:"Products" gorm:"foreignKey:OrderID"`
}
type GetOrders struct {
	ID              uint            `json:"id" gorm:"unique;not null"`
	UserID          uint            `json:"UserID" gorm:"index;foreignKey:UserID"`
	AddressID       uint            `json:"AddressID" gorm:"index;foreignKey:AddressID"`
	Status          string          `json:"status" gorm:"default:'pending'"`
	Payment         string          `json:"payment" gorm:"default:'cod'"`
	Total           float64         `json:"Total" gorm:"not null"`
	User            User            `gorm:"foreignKey:UserID"`
	Address         UserAddress     `gorm:"foreignKey:AddressID"`
	OrderedProducts []OrderProducts `json:"Products" gorm:"foreignKey:OrderID"`
}

type OrderProducts struct {
	ID           uint            `json:"id" gorm:"unique;not null"`
	OrderID      uint            `json:"orderID" gorm:"index;foreignKey:OrderID"`
	ProductID    uint            `json:"productID" gorm:"index;foreignKey:ID"`
	VariantID    uint            `json:"variantID" gorm:"index;foreignKey:VariantID"`
	Quantity     uint            `json:"quantity" gorm:"not null"`
	Status       string          `json:"status" gorm:"default:'pending'"`
	Price        float64         `json:"price" gorm:"not null"`
	Total        float64         `json:"Total" gorm:"not null"`
	Variant      ProductVariants `gorm:"foreignKey:VariantID"`
	Image        []Image         `gorm:"foreignKey:ProductID"`
	OrderDetails Orders          `json:"Products" gorm:"foreignKey:OrderID"`
	CreatedAt    time.Time
}

func (o *Orders) CalculateTotal() {
	fmt.Println("Calculating Total...")
	fmt.Println("order")
	var total float64

	for _, product := range o.OrderedProducts {
		fmt.Println("Product Status:", product.Status)
		fmt.Println("Product Total:", product.Total)

		if product.Status != "cancelled" {
			total += product.Total
		}
	}

	fmt.Println("Final Total:", total)
	o.Total = total
}
