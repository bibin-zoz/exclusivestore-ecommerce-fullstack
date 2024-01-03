package handlers

import (
	db "ecommercestore/database"
	"ecommercestore/helpers"
	"ecommercestore/models"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// admin
func CouponHandler(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Expires", "0")

	var coupons []models.Coupons
	db.DB.Find(&coupons)

	c.HTML(http.StatusOK, "managecoupons.html", gin.H{
		"Coupons": coupons,
	})

}
func DeleteCouponHandler(c *gin.Context) {

	var req DeleteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ID := req.ID

	var coupon models.Coupons
	result := db.DB.Where("id = ?", ID).Delete(&coupon)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete coupon"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "coupon deleted successfully"})
}

func AddCouponHandler(c *gin.Context) {
	// var NameError models.Invalid
	CouponCode := c.PostForm("couponCode")       // Update to match the front end field name
	Discount := c.PostForm("discountPercentage") // Update to match the front end field name
	MaxDiscount := c.PostForm("maxDiscount")     // Update to match the front end field name
	expiry := c.PostForm("expiryDate")           // Update to match the front end field name
	Status := c.PostForm("status")

	// Validate CouponCode
	if CouponCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Enter Valid coupon code"})
		return
	}

	fmt.Println("expiry", expiry)

	expiryDate, err := time.Parse("2006-01-02", expiry)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expiry date format"})
		return
	}

	discountValue, err := strconv.ParseUint(Discount, 10, 64)
	if err != nil {
		fmt.Println("Error parsing discount:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid discount value"})
		return
	}

	maxDiscountPrice, err := strconv.ParseUint(MaxDiscount, 10, 64)
	if err != nil {
		fmt.Println("Error parsing maxDiscount:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid max discount value"})
		return
	}

	// Create a new Coupons instance with the provided values
	newCoupon := models.Coupons{
		CouponCode:  CouponCode,
		Discount:    uint(discountValue),
		MaxDiscount: uint(maxDiscountPrice),
		Status:      Status,
		Expiry:      expiryDate,
	}

	result := db.DB.Create(&newCoupon)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Coupon Already Exists"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Coupon added successfully",
		"coupon":   newCoupon,
		"redirect": "/admin/coupons",
	})
}

func CouponValidatehandler(c *gin.Context) {
	CouponCode := c.Query("CouponCode")
	userID, _ := helpers.GetID(c)

	var couponDetails models.Coupons
	err := db.DB.Where("coupon_code = ?", CouponCode).First(&couponDetails).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Invalid Coupon ID"})
		return
	}
	if couponDetails.Expiry.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Coupon expired"})
		fmt.Println("expired")
		return
	}
	var cart models.Cart
	db.DB.Where("user_id=?", userID).Find(&cart)

	discountPercentage := couponDetails.Discount
	maxDiscount := couponDetails.MaxDiscount
	minPurchase := couponDetails.MinPurchase

	if cart.Total < float64(minPurchase) {
		c.JSON(http.StatusBadRequest, gin.H{"Message": fmt.Sprintf(" minimum purchase of %d%% required to apply coupon ", minPurchase)})
		fmt.Println("expired")
		return

	}
	var discountprice uint
	if (uint(cart.Total)/100)*discountPercentage < maxDiscount {
		discountprice = (uint(cart.Total) / 100) * discountPercentage
	} else {
		discountprice = (maxDiscount)
	}
	finalPrice := uint(cart.Total) - discountprice
	fmt.Println("discoutnprice:", discountprice)
	fmt.Println("discountPercentage:", discountPercentage)
	fmt.Println("maxDiscount:", maxDiscount)
	fmt.Println("(float64(discountPercentage)/cart.Total)*100", (float64(discountPercentage)/cart.Total)*100)

	c.SetCookie("couponcode", CouponCode, 0, "/", "exclusivestore.xyz", true, true)

	c.JSON(http.StatusOK, gin.H{
		"Message":       fmt.Sprintf("%d%% upto %d Applied", couponDetails.Discount, couponDetails.MaxDiscount),
		"finalPrice":    finalPrice,
		"discountprice": discountprice,
	})

}

func RemoveCouponHandler(c *gin.Context) {
	c.SetCookie("couponcode", "", -1, "/", "exclusivestore.xyz", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Coupon removed successfully"})
}
