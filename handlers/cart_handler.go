package handlers

import (
	db "ecommercestore/database"
	"ecommercestore/helpers"
	"ecommercestore/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetCarthandler(c *gin.Context) {
	authCookie, _ := c.Cookie("auth")
	var token models.TokenUser
	err := json.NewDecoder(strings.NewReader(authCookie)).Decode(&token)
	if err != nil {
		fmt.Println("Error fetching  UserDetails :", err)
		// Handle the error (e.g., return an error response)
		return
	}
	Claims, err := helpers.ParseToken(token.AccessToken)
	if err != nil {
		fmt.Println("Error fetching  UserDetails :", err)
		// Handle the error (e.g., return an error response)
		return
	}

	var Cart []models.Cart

	if err := db.DB.Preload("Product").Preload("Variant").Preload("Product.Images").Where("user_name=?", Claims.Username).Find(&Cart).Error; err != nil {
		fmt.Println("Error fetching carts:", err)
		return
	}

	c.HTML(http.StatusOK, "cart.html", gin.H{
		// "Productvariants": ProductVariants,
		"Cart": Cart,
	})

}

func AddToCarthandler(c *gin.Context) {

	authCookie, _ := c.Cookie("auth")
	var token models.TokenUser
	err := json.NewDecoder(strings.NewReader(authCookie)).Decode(&token)
	if err != nil {
		fmt.Println("Error fetching  UserDetails :", err)
		// Handle the error (e.g., return an error response)
		return
	}
	Claims, err := helpers.ParseToken(token.AccessToken)
	if err != nil {
		fmt.Println("Error fetching  UserDetails :", err)
		// Handle the error (e.g., return an error response)
		return
	}

	var Cart models.GetCart

	if err := c.ShouldBindJSON(&Cart); err != nil {
		fmt.Println("sas")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var variant models.ProductVariants
	db.DB.First(&variant, Cart.VariantID)

	UpdateCart := &models.Cart{
		UserName:  Claims.Username,
		VariantID: Cart.VariantID,
		ProductID: Cart.ProductID,
		Price:     variant.Price,
		Quantity:  Cart.Quantity,
	}

	result := db.DB.Create(UpdateCart)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cart updated successfully",
	})

}
