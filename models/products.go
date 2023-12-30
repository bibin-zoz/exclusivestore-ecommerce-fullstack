package models

import (
	"fmt"
	"time"

	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Products struct {
	gorm.Model
	ID             uint              `json:"id" gorm:"unique;not null"`
	CategoryID     uint              `json:"categoryID" gorm:"index;foreignKey:CategoryID"`
	BrandID        uint              `json:"brandID" gorm:"index;foreignKey:BrandID"`
	ProductName    string            `json:"productName"`
	ProductDetails string            `json:"productDetails"`
	Status         string            `json:"status" gorm:"default:'listed'"`
	Category       Categories        `json:"category" gorm:"foreignKey:CategoryID"`
	Brand          Brands            `json:"brand" gorm:"foreignKey:BrandID"`
	Discount       uint              `json:"discount" gorm:"not null;default:0"`
	Images         []Image           `json:"images" gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	Variants       []ProductVariants `json:"product_variants" gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	// Carts          []Cart            `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
}

type ProductVariants struct {
	gorm.Model
	ID            uint     `json:"id" gorm:"unique;not null"`
	ProductID     uint     `json:"productID" gorm:"index;foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	Processor     string   `json:"processor"`
	Storage       string   `json:"storage"`
	Ram           string   `json:"ram"`
	Stock         int      `json:"stock"`
	Status        string   `json:"status" gorm:"default:'listed'"`
	Price         float64  `json:"price"`
	MaxPrice      float64  `json:"maxprice"`
	Slug          string   `json:"slug" gorm:"uniqueIndex"`
	DiscountPrice uint     `json:"discountprice" gorm:"not null;default:0"`
	Product       Products `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	// Carts     []CartProducts `gorm:"foreignKey:VariantID"`
}

func (pv *ProductVariants) CalculateDiscountPrice() {
	pv.Price = pv.Price - float64(pv.Product.Discount) - float64(pv.Product.Category.CategoryOffer.Discount)

}

func (pv *ProductVariants) CreateSlug(productName string) {
	// Combine relevant fields to form a string
	slugInput := fmt.Sprintf("%s-%s-%s", productName, pv.Storage, pv.Ram)

	// Generate the slug using gosimple/slug
	newSlug := slug.MakeLang(slugInput, "en")

	// Set the generated slug to the Slug field
	pv.Slug = newSlug

}

type Image struct {
	ID        uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductID uint   `json:"productID" gorm:"index;foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	FilePath  string `json:"filepath" gorm:"not null"`
	Product   Products
}

type Brands struct {
	ID        uint   `json:"id" gorm:"unique;not null"`
	BrandName string `json:"brandname" gorm:"unique;not null"`
	CreatedAt time.Time
}

type Productview struct {
	ID             uint
	ProductID      uint `json:"productID" gorm:"index;foreignKey:ProductID"`
	ProductName    string
	ProductDetails string

	Images []Image `json:"images" gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
}
