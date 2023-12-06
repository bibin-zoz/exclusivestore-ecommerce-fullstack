package handlers

import (
	db "ecommercestore/database"
	"ecommercestore/helpers"
	"ecommercestore/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetCarthandler(c *gin.Context) {
	authCookie, _ := c.Cookie("auth")
	var token models.TokenUser
	err := json.NewDecoder(strings.NewReader(authCookie)).Decode(&token)
	if err != nil {
		fmt.Println("Error fetching  UserDetails :", err)

		return
	}
	Claims, err := helpers.ParseToken(token.AccessToken)
	if err != nil {
		fmt.Println("Error fetching  UserDetails :", err)
		return
	}
	fmt.Println("claimscart", Claims)

	var Cart []models.Cart
	if err := db.DB.Preload("Product").Preload("Variant").Preload("Product.Images").Where("user_id=?", Claims.ID).Find(&Cart).Error; err != nil {
		fmt.Println("Error fetching carts:", err)
		return
	}
	var totalSum float64

	for _, cartItem := range Cart {
		totalSum += cartItem.Total
	}

	fmt.Println("Total Sum of Cart.TotalPrice:", totalSum)

	c.HTML(http.StatusOK, "cart.html", gin.H{
		// "Productvariants": ProductVariants,
		"Cart":      Cart,
		"CartTotal": totalSum,
	})

}

func AddToCarthandler(c *gin.Context) {

	ID, err := helpers.GetID(c)
	if ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var Cart models.GetCart

	if err := c.ShouldBindJSON(&Cart); err != nil {
		// fmt.Println("sas")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var count int64
	var cartcount models.Cart
	db.DB.Where("user_id=? AND product_id=?", ID, Cart.ProductID).Find(&cartcount).Count(&count)
	fmt.Println("count", count)
	if count != 0 {
		fmt.Println("product already in cart")
		c.JSON(http.StatusOK, gin.H{
			"message": "Product Already in cart....Rediecting"})
		return

	}

	var variant models.ProductVariants
	db.DB.First(&variant, Cart.VariantID)

	UpdateCart := &models.Cart{
		UserID:    ID,
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

func DeleteCartHandler(c *gin.Context) {
	var req DeleteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cartID := req.ID

	var cart models.Cart
	result := db.DB.Where("id = ?", cartID).Delete(&cart)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove from Cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item removed from cart successfully"})
}

func CheckOuthandler(c *gin.Context) {
	authCookie, _ := c.Cookie("auth")
	var token models.TokenUser
	err := json.NewDecoder(strings.NewReader(authCookie)).Decode(&token)
	if err != nil {
		fmt.Println("Error fetching  UserDetails :", err)

		return
	}
	Claims, err := helpers.ParseToken(token.AccessToken)
	if err != nil {
		fmt.Println("Error fetching  UserDetails :", err)

		return
	}
	fmt.Println("claimscart", Claims)

	var Cart []models.Cart
	if err := db.DB.Preload("Product").Preload("Variant").Preload("Product.Images").Where("user_id=?", Claims.ID).Find(&Cart).Error; err != nil {
		fmt.Println("Error fetching carts:", err)
		return
	}
	if len(Cart) == 0 {

		c.Redirect(http.StatusTemporaryRedirect, "/home")
	}
	var totalSum float64
	var cartid uint

	for _, cartItem := range Cart {
		cartid = cartItem.ID
		totalSum += cartItem.Total
	}

	fmt.Println("Total Sum of Cart.TotalPrice:", totalSum)

	c.HTML(http.StatusOK, "checkout.html", gin.H{
		// "Productvariants": ProductVariants,
		"Cart":      Cart,
		"CartID":    cartid,
		"CartTotal": totalSum,
	})
}

func OrderPlacehandler(c *gin.Context) {
	var orderReq models.OrderReq
	var cart []models.Cart
	userid, _ := helpers.GetID(c)

	if err := c.ShouldBindJSON(&orderReq); err != nil {
		fmt.Println("bind error order:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if orderReq.AddressID == "" || orderReq.CartID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Select an address or add a new address"})
		return
	}

	result := db.DB.Debug().Where("user_id=?", userid).Find(&cart)
	if result.Error != nil {
		fmt.Println("error fetching cart:", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	addressID, err := strconv.ParseUint(orderReq.AddressID, 10, 64)
	if err != nil {
		fmt.Println("error parsing AddressID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid AddressID"})
		return
	}

	var orderDetail models.Orders
	orderDetail.AddressID = uint(addressID)
	orderDetail.UserID = userid
	for _, cartItem := range cart {
		orderDetail.Total += cartItem.Total

	}

	result = db.DB.Debug().Create(&orderDetail)
	if result.Error != nil {
		fmt.Println("error creating order detail:", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	fmt.Println("orderDetailid", orderDetail)

	for i, cartItem := range cart {
		fmt.Println("count", i)
		var orderProduct models.OrderProducts
		orderProduct.OrderID = orderDetail.ID
		orderProduct.ProductID = cartItem.ProductID
		orderProduct.VariantID = cartItem.VariantID
		orderProduct.Quantity = cartItem.Quantity
		orderProduct.Price = cartItem.Price
		orderProduct.Total = cartItem.Total

		result = db.DB.Debug().Create(&orderProduct)
		if result.Error != nil {
			fmt.Println("error creating order product entry:", result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
	}

	result = db.DB.Where("user_id = ?", userid).Delete(&models.Cart{})
	if result.Error != nil {
		fmt.Println("error deleting cart item:", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order placed successfully"})
}

func GetOrdershandler(c *gin.Context) {
	var orders []models.Orders
	ID, err := helpers.GetID(c)
	if err != nil {
		fmt.Println("error", err)
	}
	db.DB.Where("user_id=?", ID).Find(&orders)
	// fmt.Println("useradd", userAddress)
	c.JSON(http.StatusOK, orders)

}

func CancelOrderHandler(c *gin.Context) {
	var updateStatusRequest struct {
		ID uint `json:"id"`
	}

	if err := c.ShouldBindJSON(&updateStatusRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order models.Orders
	if err := db.DB.First(&order, updateStatusRequest.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	if order.Status == "cancelled" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item already Cancelled "})
		return

	}

	order.Status = "cancelled"
	if err := db.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated successfully"})
}

func CancelProductHandler(c *gin.Context) {
	var updateStatusRequest struct {
		OrderID   uint `json:"orderID"`
		ProductID uint `json:"productID"`
	}

	if err := c.ShouldBindJSON(&updateStatusRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order models.OrderProducts
	if err := db.DB.Where("order_id = ? AND product_id=?", updateStatusRequest.OrderID, updateStatusRequest.ProductID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	if order.Status == "cancelled" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item already Cancelled "})
		return

	}

	order.Status = "cancelled"
	if err := db.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}
	var OrderDetail models.Orders
	// Fetch an order with preloaded OrderedProducts
	if err := db.DB.Preload("OrderedProducts").Where("id = ?", updateStatusRequest.OrderID).First(&OrderDetail).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	OrderDetail.CalculateTotal()

	// Save to the database
	db.DB.Save(&OrderDetail)

	c.JSON(http.StatusOK, gin.H{"message": "Status updated successfully"})
}

func TrackOrderHandler(c *gin.Context) {
	var orders []models.OrderProducts
	var OrderDetails models.Orders
	_, err := helpers.GetID(c)
	orderid := c.Query("id")
	if err != nil {
		fmt.Println("error", err)
	}
	fmt.Println("id", orderid)
	db.DB.Preload("Variant.Product.Images").Preload("Variant.Product").Preload("Variant").Preload("Image").Where("order_id=?", orderid).Find(&orders)
	db.DB.Preload("Address").Where("id=?", orderid).Find(&OrderDetails)
	fmt.Println("hiii")
	// c.JSON(http.StatusOK, orders)
	c.HTML(http.StatusOK, "orderdetails.html", gin.H{
		// "Productvariants": ProductVariants,
		"Order":        orders,
		"OrderDetails": OrderDetails,
	})

}
