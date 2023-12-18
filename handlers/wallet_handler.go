package handlers

import (
	db "ecommercestore/database"
	"ecommercestore/helpers"
	"ecommercestore/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func WalletHandler(c *gin.Context) {
	userID, _ := helpers.GetID(c)
	var WalletDetails models.Wallet
	// var WalletTransactions []models.Transaction

	result := db.DB.Preload("Transactions").Where("user_id = ?", userID).Find(&WalletDetails)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, nil)
		return

	}
	c.JSON(http.StatusOK, gin.H{
		"WalletDetails": WalletDetails,
	})

}
