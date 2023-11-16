package models

import (
	"time"

	"gorm.io/gorm"
)

type Products struct {
	ID             uint           `json:"id" gorm:"unique;not null"`
	CategoryID     int            `json:"categoryID" gorm:"foreignkey:CategoryID;constraint:OnDelete:CASCADE"`
	ProductName    string         `json:"product_name"`
	ProductDetails string         `json:"product_details"`
	Image          string         `json:"image"`
	Storage        string         `json:"storage"`
	Ram            string         `json:"ram"`
	Stock          int            `json:"stock"`
	Status         string         `json:"status" gorm:"default:'listed'"`
	Price          float64        `json:"price"`
	Images         []ProductImage `json:"ProductImage" gorm:"foreignKey:ProductID"`
}
type ProductImage struct {
	ID        uint   `json:"id" gorm:"unique;not null"`
	ProductID uint   `json:"product_id"`
	Image     string `json:"image"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

//	type Category struct {
//		ID       uint   `json:"id" gorm:"unique;not null"`
//		Category string `json:"category"`
//		Image    string `json:"category_image"`
//	}
type Categories struct {
	ID           uint   `json:"id" gorm:"unique;not null"`
	CategoryName string `json:"category" gorm:"unique;not null"`
	Status       string `json:"status" gorm:"default:'listed'"`
	Image        string `json:"category_image"`
	CreatedAt    time.Time
}

type Productview struct {
	ID             uint
	ProductName    string
	CategoryName   string
	ProductDetails string
	Status         string
	Ram            string
	Storage        string
	Stock          int
	Price          float64
}
