package handlers

import (
	"ecommercestore/database"
	"ecommercestore/models"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetProductsHandler(c *gin.Context) {
	var products []models.ProductVariants

	baseQuery := database.DB.Preload("Product").
		Preload("Product.Category").
		Preload("Product.Category.CategoryOffer").
		Preload("Product.Images").
		Where("status='listed'")

	if err := baseQuery.Find(&products).Error; err != nil {
		fmt.Println("Error fetching product variants with product and images:", err)
		return
	}

	var Category []models.Categories
	database.DB.Find(&Category)

	c.HTML(http.StatusOK, "shop-product.html", gin.H{
		"ProductVariants": products,
		"Category":        Category,
	})
}

func FilterProductshandler(c *gin.Context) {
	var products []models.ProductVariants
	ram := c.PostForm("ram")
	storage := c.PostForm("storage")
	categoryID := c.PostForm("category")

	fmt.Println("ram", ram)

	baseQuery := database.DB.
		Joins("JOIN products ON product_variants.product_id = products.id").
		Joins("JOIN categories ON products.category_id = categories.id").
		Preload("Product").
		Preload("Product.Category").
		Preload("Product.Images").
		Where("product_variants.status='listed'")

	if ram != "" {
		baseQuery = baseQuery.Where("product_variants.ram IN (?)", strings.Split(ram, ","))
	}
	if storage != "" {
		baseQuery = baseQuery.Where("product_variants.storage IN (?)", strings.Split(storage, ","))
	}
	if categoryID != "" {
		baseQuery = baseQuery.Where("categories.id IN (?)", strings.Split(categoryID, ","))
	}

	if err := baseQuery.Find(&products).Error; err != nil {
		fmt.Println("Error fetching product variants with product and images:", err)
		return
	}

	var Category []models.Categories
	database.DB.Find(&Category)

	c.JSON(http.StatusOK, gin.H{
		"ProductVariants": products,
		"Category":        Category,
		"FilterRam":       ram,
		"FilterStorage":   storage,
		"FilterCategory":  categoryID,
	})
}
