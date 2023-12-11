package handlers

import (
	db "ecommercestore/database"
	"ecommercestore/helpers"
	"ecommercestore/models"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/razorpay/razorpay-go"
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
	db.DB.Where("user_id=? AND variant_id=?", ID, Cart.VariantID).Find(&cartcount).Count(&count)
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
func UpdateQuantityHandler(c *gin.Context) {
	fmt.Println("hiii")
	var updateCart models.Updatecart
	var cart models.Cart

	err := c.ShouldBindJSON(&updateCart)
	if err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	result := db.DB.Preload("Variant").Where("id = ?", updateCart.ID).Find(&cart)
	if result.Error != nil {
		fmt.Println("Error fetching cart:", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart"})
		return
	}

	newQuantity, err := strconv.Atoi(updateCart.Quantity)
	if err != nil {
		fmt.Println("Error converting quantity to integer:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quantity"})
		return
	}

	if cart.Variant.Stock+int(cart.Quantity) < newQuantity {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to update quantity"})
		return
	}

	cart.Quantity = uint(newQuantity)
	db.DB.Save(&cart)

	c.JSON(http.StatusOK, gin.H{
		"Quantity": newQuantity,
		"Total":    cart.Total,
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
	if orderReq.PaymentMethod == "" || orderReq.CartID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error to fetch payment method"})
		return
	}

	result := db.DB.Where("user_id=?", userid).Find(&cart)
	if result.Error != nil {
		fmt.Println("error fetching cart:", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	// result = db.DB.Where("id=?", cart.product_id).Find(&cart)
	// if result.Error != nil {
	// 	fmt.Println("error fetching cart:", result.Error)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
	// 	return
	// }

	addressID, err := strconv.ParseUint(orderReq.AddressID, 10, 64)
	if err != nil {
		fmt.Println("error parsing AddressID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid AddressID"})
		return
	}

	var orderDetail models.Orders
	orderDetail.AddressID = uint(addressID)
	orderDetail.UserID = userid
	orderDetail.Payment = orderReq.PaymentMethod
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

	c.JSON(http.StatusOK, gin.H{
		"message": "Order placed successfully",
		"orderID": orderDetail.ID,
	})
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

func ReturnOrderHandler(c *gin.Context) {
	// Parse JSON request body into an OrderProduct struct
	var returnRequest models.UserRequest
	var orderdetails models.OrderProducts
	if err := c.ShouldBindJSON(&returnRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Where("id=?", returnRequest.ID).First(&orderdetails).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order not found"})
		return
	}

	// Check if the return period has expired
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	if orderdetails.CreatedAt.Before(sevenDaysAgo) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Return period expired"})
		return
	}
	result := db.DB.Debug().Table("order_products").Where("id=?", returnRequest.ID).Update("notes", "Return Request:"+returnRequest.Request).Update("status", "pending")
	if result.Error != nil {
		fmt.Println("error", result.Error)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to request return"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Return request processed successfully"})
}

func CreateRazorpayOrder(c *gin.Context) {
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
	fmt.Println("reqdata", requestData)
	var cart models.Cart

	result := db.DB.Where("ID=?", requestData.CartID).Find(&cart)
	if result.Error != nil {
		fmt.Println("Error fetching cart details:", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart details"})
	}

	// You may need to calculate the order amount based on the items in the cart or other factors.
	// For simplicity, I'm assuming a fixed amount of 1000 here.
	orderAmount := cart.Total

	params := map[string]interface{}{
		"amount":          orderAmount * 100,
		"currency":        "INR",
		"payment_capture": 1,
	}
	razorpayKey := os.Getenv("RAZORPAY_KEY_ID")
	razorpaySecret := os.Getenv("RAZORPAY_KEY_SECRET")

	client := razorpay.NewClient(razorpayKey, razorpaySecret)

	order, err := client.Order.Create(params, nil)
	if err != nil {
		fmt.Println("Error creating Razorpay order:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Razorpay order"})
		return
	}

	// Type assert the 'id' field from the map
	orderID, ok := order["id"].(string)
	if !ok {
		fmt.Println("Invalid order ID format:", order)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid order ID format"})
		return
	}
	fmt.Println("aaaaaaa")

	c.JSON(http.StatusOK, gin.H{
		"id":     orderID,
		"amount": strconv.Itoa(int(orderAmount)),
	})
}
