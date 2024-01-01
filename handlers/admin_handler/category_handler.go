package handlers

import (
	db "ecommercestore/database"
	"ecommercestore/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// category
func Categoryhandler(c *gin.Context) {
	var category []models.Categories
	db.DB.Find(&category)

	c.HTML(http.StatusOK, "categories.html", gin.H{
		"Category": category,
	})

}
func CategoryPost(c *gin.Context) {
	categoryName := c.PostForm("categoryName")
	status := c.PostForm("status")
	var NameError models.Invalid

	if categoryName == "" {
		NameError.NameError = "Enter valid Category Name"
		c.JSON(http.StatusBadRequest, gin.H{
			"Errors": NameError,
		})
		return
	}

	newCategory := models.Categories{
		CategoryName: categoryName,
		Status:       status,
	}

	result := db.DB.Create(&newCategory)

	if result.Error != nil {
		NameError.NameError = "Category Already Exists"

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  result.Error.Error(),
			"Errors": NameError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{

		"message":  "Category added successfully",
		"category": newCategory,
		"redirect": "/admin/categories",
	})

}

func DeleteCategoryHandler(c *gin.Context) {

	var req DeleteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	categoryID := req.ID

	var category models.Categories
	result := db.DB.Where("id = ?", categoryID).Delete(&category)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete customer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
}

func UpdateCategoryStatus(c *gin.Context) {

	var req DeleteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ID := req.ID

	var category models.Categories
	if err := db.DB.Where("id = ?", ID).First(&category).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	if category.Status == "listed" {
		category.Status = "unlisted"
	} else {
		category.Status = "listed"
	}

	if err := db.DB.Save(&category).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Category Status Updated Successfully",
		"redirect": "/admin/categories",
	})
}
