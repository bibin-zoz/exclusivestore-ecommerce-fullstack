package handlers

import (
	db "ecommercestore/database"
	"ecommercestore/helpers"
	"ecommercestore/models"
	"fmt"
	"net/http"
	"strconv"

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
func WalletOrderhandler(c *gin.Context) {
	userID, _ := helpers.GetID(c)
	var requestData models.OrderReq
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if requestData.AddressID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "select Delivery address"})
		return

	}
	if requestData.CartID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart Details not found"})
		return

	}
	if requestData.CartID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart Details not found"})
		return

	}
	fmt.Println("reqdata", requestData)
	var cart models.Cart

	result := db.DB.Debug().Where("ID=?", requestData.CartID).Find(&cart)
	if result.Error != nil {
		fmt.Println("Error fetching cart details:", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart details"})
		return
	}
	fmt.Println("cart", cart)
	coupon, err := c.Cookie("couponcode")
	var orderAmount float64
	orderAmount = cart.Total
	if err == nil {

		var couponDetails models.Coupons
		db.DB.Where("coupon_code=?", coupon).Find(&couponDetails)
		fmt.Println("couponDetails")
		discount := float64(helpers.DiscountPrice(couponDetails, c))
		fmt.Println("not", orderAmount)
		orderAmount = orderAmount - discount

	}
	var wallet models.Wallet
	db.DB.Preload("Transactions").Where("user_id=?", userID).Find(&wallet)
	if wallet.Balance < orderAmount {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Not Sufficent Balance in wallet"})
		return
	}
	var Transaction models.Transaction
	Transaction.WalletID = wallet.ID
	Transaction.Amount = int(orderAmount)
	Transaction.Type = "debit"
	Transaction.Description = "Product Purchase"

	db.DB.Create(&Transaction)

	db.DB.Save(&wallet)

	c.JSON(http.StatusOK, gin.H{

		"amount": strconv.Itoa(int(orderAmount)),
	})
}
