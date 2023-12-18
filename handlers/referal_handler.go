package handlers

import (
	db "ecommercestore/database"
	"ecommercestore/models"
	"fmt"
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
	var order models.OrderProducts
	db.DB.Preload("OrderDetails").Where("id=38").Find(&order)
	fmt.Println("order", order.ID)
	fmt.Println("orsssssssder", order.OrderDetails.ID)
	Discount := order.OrderDetails.Discount
	TotalOrder := order.OrderDetails.Total + float64(Discount)
	ProductPrice := order.Total
	if Discount != 0 {
		productDiscount := (ProductPrice / TotalOrder) * float64(Discount)
		refundAmount := uint(ProductPrice) - uint(productDiscount)
		fmt.Println(refundAmount)
		return

	}

	refundAmount := uint(ProductPrice)
	fmt.Println(refundAmount)

}
