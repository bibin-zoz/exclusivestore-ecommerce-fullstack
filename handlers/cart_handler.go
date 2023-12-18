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

	var Cart models.Cart
	if err := db.DB.Where("user_id=?", Claims.ID).Find(&Cart).Error; err != nil {
		fmt.Println("Error fetching carts:", err)
		return
	}

	var CartProducts []models.CartProducts
	if err := db.DB.Preload("Product").Preload("Variant").Preload("Product.Images").Where("cart_id=?", Cart.ID).Find(&CartProducts).Error; err != nil {
		fmt.Println("Error fetching carts:", err)
		return
	}
	var totalSum float64

	for _, cartItem := range CartProducts {
		totalSum += cartItem.Total
	}

	fmt.Println("Total Sum of Cart.TotalPrice:", totalSum)

	c.HTML(http.StatusOK, "cart.html", gin.H{
		// "Productvariants": ProductVariants,
		"Cart":         Cart,
		"CartProducts": CartProducts,
		"CartTotal":    totalSum,
	})
	// c.JSON(http.StatusOK, gin.H{
	// 	// "Productvariants": ProductVariants,
	// 	"Cart":         Cart,
	// 	"CartProducts": CartProducts,
	// 	"CartTotal":    totalSum,
	// })

}

func AddToCarthandler(c *gin.Context) {
	// Retrieve user ID
	userID, err := helpers.GetID(c)
	if userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Bind JSON data to AddCart model
	var AddCart models.GetCart
	if err := c.ShouldBindJSON(&AddCart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the product is already in the cart
	var count int64
	var cartCount models.CartProducts
	var cartDetails models.Cart

	result := db.DB.Where("user_id=?", userID).Attrs(models.Cart{UserID: userID}).FirstOrCreate(&cartDetails)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	db.DB.Where("cart_id=? AND variant_id=?", cartDetails.ID, AddCart.VariantID).Find(&cartCount).Count(&count)

	if count != 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Product Already in cart. Redirecting..."})
		return
	}

	// Retrieve product variant details
	var variant models.ProductVariants
	db.DB.First(&variant, AddCart.VariantID)

	// Find or create a cart for the user

	// Create a new CartProducts entry
	var cartID uint
	db.DB.Table("carts").Where("user_id=?", userID).Select("id").Scan(&cartID)
	fmt.Println("cartid", cartDetails.ID)

	cartProduct := &models.CartProducts{
		CartID:    cartDetails.ID,
		VariantID: AddCart.VariantID,
		ProductID: AddCart.ProductID,
		Price:     variant.Price,
		Quantity:  AddCart.Quantity,
		Total:     variant.Price * float64(AddCart.Quantity),
	}

	// Create the cart product entry
	result = db.DB.Create(cartProduct)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	var CartTotal models.Cart
	// Fetch an order with preloaded OrderedProducts
	if err := db.DB.Preload("CartProducts").Where("id = ?", cartDetails.ID).First(&CartTotal).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	CartTotal.CalculateTotal()

	// Save to the database
	db.DB.Save(&CartTotal)

	c.JSON(http.StatusOK, gin.H{"message": "Product added to cart successfully"})
}

func UpdateQuantityHandler(c *gin.Context) {
	fmt.Println("hiii")
	var updateCart models.Updatecart
	var cart models.CartProducts

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

	var cart models.CartProducts
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

	var Cart models.Cart
	if err := db.DB.Preload("CartProducts").Preload("CartProducts.Product").Preload("CartProducts.Variant").Preload("CartProducts.Product.Images").Where("user_id=?", Claims.ID).Find(&Cart).Error; err != nil {
		fmt.Println("Error fetching carts:", err)
		return
	}
	if len(Cart.CartProducts) == 0 {

		c.Redirect(http.StatusTemporaryRedirect, "/home")
	}
	var totalSum float64
	var cartid uint

	for _, cartItem := range Cart.CartProducts {
		cartid = cartItem.CartID
		totalSum += cartItem.Total
	}

	fmt.Println("Total Sum of Cart.TotalPrice:", totalSum)

	couponcode, _ := c.Cookie("couponcode")
	fmt.Println("couponcode", couponcode)

	c.HTML(http.StatusOK, "checkout.html", gin.H{
		// "Productvariants": ProductVariants,
		"Cart":       Cart,
		"CartID":     cartid,
		"CartTotal":  totalSum,
		"CouponCode": couponcode,
	})
}

// order
func OrderPlacehandler(c *gin.Context) {
	var orderReq models.OrderReq
	var cart models.Cart
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

	result := db.DB.Preload("CartProducts").Where("user_id=?", userid).Find(&cart)
	if result.Error != nil {
		fmt.Println("error fetching cart:", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	var couponDetails models.Coupons
	CouponCode, _ := c.Cookie("couponcode")
	var coupon bool
	result = db.DB.Where("coupon_code=?", CouponCode).Find(&couponDetails)
	if result.Error != nil {
		fmt.Println("error fetching coupon details:", result.Error)
		coupon = false
	} else {
		coupon = true
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
	orderDetail.Total = cart.Total
	if coupon {
		fmt.Println("dsdsd")
		orderDetail.Discount = helpers.DiscountPrice(couponDetails, c)
		orderDetail.Total = orderDetail.Total - float64(orderDetail.Discount)

	}

	result = db.DB.Debug().Create(&orderDetail)
	if result.Error != nil {
		fmt.Println("error creating order detail:", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	fmt.Println("orderDetailid", orderDetail)

	for i, cartItem := range cart.CartProducts {
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

	result = db.DB.Where("cart_id = ?", cart.ID).Delete(&models.CartProducts{})
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
	// db.DB.Where("user_id=?", ID).Find(&orders)
	db.DB.Where("user_id = ?", ID).Order("created_at DESC").Find(&orders)
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

// cancel
func CancelProductHandler(c *gin.Context) {
	ID, _ := helpers.GetID(c)
	var updateStatusRequest struct {
		OrderID   uint `json:"orderID"`
		ProductID uint `json:"productID"`
	}

	if err := c.ShouldBindJSON(&updateStatusRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order models.OrderProducts
	if err := db.DB.Preload("OrderDetails").Where("order_id = ? AND product_id=?", updateStatusRequest.OrderID, updateStatusRequest.ProductID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	if order.Status == "cancelled" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item already Cancelled "})
		return
	}

	// Calculate refund amount and product discount
	refundAmount, productDiscount := helpers.GetRefundAmount(&order)
	if order.OrderDetails.Payment != "cod" {
		var wallet models.Wallet
		db.DB.Preload("Transactions").Where("user_id=?", ID).Find(&wallet)
		var Transaction models.Transaction
		Transaction.WalletID = wallet.ID
		Transaction.Amount = int(refundAmount)
		Transaction.Type = "credit"
		Transaction.Description = fmt.Sprintf("refund for order %d", order.ID)

		db.DB.Create(&Transaction)
		wallet.Balance += float64(refundAmount)
		db.DB.Save(&wallet)

	}

	fmt.Print("refund details", refundAmount, productDiscount)

	// Update order's total and discount
	var OrderDetail models.Orders
	if err := db.DB.Preload("OrderedProducts").Where("id = ?", updateStatusRequest.OrderID).First(&OrderDetail).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Update discount and total
	OrderDetail.Discount -= uint(productDiscount)
	OrderDetail.Total -= float64(refundAmount)

	// Save to the database
	if err := db.DB.Save(&OrderDetail).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order details"})
		return
	}

	// Update order status
	order.Status = "cancelled"
	if err := db.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	// Respond with success message
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
