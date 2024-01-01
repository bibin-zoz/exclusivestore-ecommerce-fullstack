package handlers

import (
	db "ecommercestore/database"
	"ecommercestore/helpers"
	"ecommercestore/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func ProductOffersHandler(c *gin.Context) {
	var Categories []models.Categories
	db.DB.Preload("CategoryOffer").Find(&Categories)

	c.HTML(http.StatusOK, "categoryoffers.html", gin.H{
		"Category": Categories,
	})

}
func AddProductOfferhandler(c *gin.Context) {

	discount, _ := strconv.Atoi(c.PostForm("discount"))
	expiryDateStr := c.PostForm("expiryDate")
	status := c.PostForm("status")
	var NameError models.Invalid
	categoryID, err := strconv.Atoi(c.PostForm("categoryID"))
	if err != nil {
		NameError.NameError = "Enter valid Category "
		c.JSON(http.StatusBadRequest, gin.H{
			"Errors": NameError,
		})
		return
	}
	if discount < 100 {
		NameError.AmountError = "Enter valid Discount amount "
		c.JSON(http.StatusBadRequest, gin.H{
			"Errors": NameError,
		})
		return
	}
	expiryDate, err := time.Parse("2006-01-02", expiryDateStr)
	if err != nil {
		NameError.DateError = "Enter valid expiry Date "
		c.JSON(http.StatusBadRequest, gin.H{
			"Errors": NameError,
		})
		return
	}

	newCategory := models.CategoryOffer{
		CategoryID: uint(categoryID),
		Status:     status,
		Discount:   uint(discount),
		ExpiryAt:   expiryDate,
	}

	result := db.DB.Create(&newCategory)

	if result.Error != nil {
		NameError.NameError = "Category offer Already Exists"

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  result.Error.Error(),
			"Errors": NameError,
		})
		return
	}
	helpers.UpdateDiscountPrice()

	c.JSON(http.StatusOK, gin.H{

		"message":  "Category added successfully",
		"category": newCategory,
		"redirect": "/admin/categoryoffers",
	})

}

func DeleteProductOfferHandler(c *gin.Context) {

	var req DeleteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ID := req.ID

	var category models.CategoryOffer
	result := db.DB.Where("id = ?", ID).Delete(&category)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove Category offer"})
		return
	}
	helpers.UpdateDiscountPrice()

	c.JSON(http.StatusOK, gin.H{"message": "offer removed successfully"})
}
