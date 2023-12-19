package helpers

import (
	db "ecommercestore/database"
	"ecommercestore/models"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

func GenerateRandomReferalCode() string {
	guid := xid.New()
	code := guid.String()
	return code[:5]
}
func UpdateReferalCount(referralCode string) error {

	var referalDetails models.ReferalDetails
	err := db.DB.Where("referal_code = ?", referralCode).First(&referalDetails).Error
	if err != nil {
		return err
	}
	referalDetails.Count++
	err = db.DB.Save(&referalDetails).Error
	if err != nil {
		return err
	}

	return nil
}

func UpdateTransaction(txDetails models.Transaction) error {
	result := db.DB.Create(&txDetails)
	return result.Error
}

func GetRefundAmount(order *models.OrderProducts) (uint, float64) {

	Discount := order.OrderDetails.Discount
	TotalOrder := order.OrderDetails.Total + float64(Discount)
	ProductPrice := order.Total
	fmt.Println("Details", Discount, TotalOrder, ProductPrice)
	if Discount != 0 {

		productDiscount := (ProductPrice / TotalOrder) * float64(Discount)
		refundAmount := uint(ProductPrice) - uint(productDiscount)
		order.OrderDetails.Discount = order.OrderDetails.Discount - uint(productDiscount)
		order.OrderDetails.Total = TotalOrder - float64(refundAmount)
		return refundAmount, productDiscount
	}

	refundAmount := uint(ProductPrice)
	fmt.Println(refundAmount)
	return refundAmount, 0

}

func DiscountPrice(couponDetails models.Coupons, c *gin.Context) uint {
	var cart models.Cart
	userID, _ := GetID(c)
	db.DB.Where("user_id=?", userID).Find(&cart)

	discountPercentage := couponDetails.Discount
	maxDiscount := couponDetails.MaxDiscount

	var discountprice uint
	if (uint(cart.Total)/100)*discountPercentage < maxDiscount {
		discountprice = (uint(cart.Total) / 100) * discountPercentage
	} else {
		discountprice = (maxDiscount)
	}
	fmt.Println("discountprice", discountprice)
	return discountprice
}

func GetProductDiscountPrice(variantID uint) uint {
	var product models.ProductVariants

	err := db.DB.Preload("Product").Preload("Product.Category").Preload("Product.Category.CategoryOffer").Where("id = ?", variantID).First(&product).Error
	if err != nil {
		return 0
	}

	if product.Product.Category.CategoryOffer.ExpiryAt.Before(time.Now()) ||
		product.Product.Category.CategoryOffer.Status != "active" {
		discountedPrice := product.Price - float64(product.Product.Discount)
		return uint(discountedPrice)
	}

	discountedPrice := product.Price - float64(product.Product.Category.CategoryOffer.Discount) - float64(product.Product.Discount)

	return uint(discountedPrice)
}

func UpdateDiscountPrice() {
	var products []models.ProductVariants
	db.DB.Preload("Product").Preload("Product.Category").Preload("Product.Category.CategoryOffer").Find(&products)
	for i := range products {

		products[i].DiscountPrice = GetProductDiscountPrice(products[i].ID)
		if products[i].DiscountPrice == uint(products[i].Price) {
			products[i].DiscountPrice = 0

		}
		fmt.Println("products[i].DiscountPrice", products[i].DiscountPrice)
		db.DB.Save(&products[i])
	}
}
