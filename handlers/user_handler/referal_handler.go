package handlers

import (
	db "ecommercestore/database"
	"ecommercestore/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ReferalValidatehandler(c *gin.Context) {
	referralCode := c.Query("referralID")

	var referalDetails models.ReferalDetails
	err := db.DB.Where("referal_code = ?", referralCode).First(&referalDetails).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Invalid Referral ID"})
		return
	}
	var wallet models.Wallet
	db.DB.Preload("Transactions").Where("user_id=?", referalDetails.UserID).Find(&wallet)
	var Transaction models.Transaction
	Transaction.WalletID = wallet.ID
	Transaction.Amount = 100
	Transaction.Type = "credit"
	Transaction.Description = "referral bonus"

	db.DB.Create(&Transaction)
	wallet.Balance += 100
	db.DB.Save(&wallet)
	c.JSON(http.StatusOK, gin.H{"Message": "Referral ID is valid"})

}

func TestHandler(c *gin.Context) {
	var Wallet models.Wallet
	db.DB.Preload("Transactions").Find(&Wallet)
	db.DB.Save(&Wallet)

}
