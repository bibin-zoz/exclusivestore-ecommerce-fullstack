package handlers

import (
	db "ecommercestore/database"
	"ecommercestore/helpers"
	"ecommercestore/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCarthandler(c *gin.Context) {
	Token, _ := c.Cookie("token")
	_, Username, _ := helpers.GetUserRoleFromToken(Token)

	var Cart []models.Cart

	if err := db.DB.Preload("Product").Preload("Variant").Preload("Product.Images").Where("user_name=?", Username).Find(&Cart).Error; err != nil {
		fmt.Println("Error fetching carts:", err)
		return
	}
	fmt.Println("Cart", Cart)
	c.HTML(http.StatusOK, "cart.html", gin.H{
		// "Productvariants": ProductVariants,
		"Cart": Cart,
	})

}

func AddToCarthandler(c *gin.Context) {
	Token, _ := c.Cookie("token")
	_, Username, _ := helpers.GetUserRoleFromToken(Token)
	fmt.Println("user:v", Username)
	var Cart models.GetCart

	if err := c.ShouldBindJSON(&Cart); err != nil {
		fmt.Println("sas")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("cart", Cart)

	result := db.DB.Create(&models.Cart{
		UserName:  Username,
		VariantID: Cart.VariantID,
		ProductID: Cart.ProductID,
	})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cart updated successfully",
	})

}
