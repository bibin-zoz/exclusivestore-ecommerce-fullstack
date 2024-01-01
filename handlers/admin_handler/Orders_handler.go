package handlers

import (
	db "ecommercestore/database"
	"ecommercestore/models"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func UserOrdersHandler(c *gin.Context) {
	var orders []models.Orders
	pgno, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		pgno = 1
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 10
	}
	offset := (pgno - 1) * limit

	var count int64
	db.DB.Model(models.Orders{}).Count(&count)
	fmt.Println("count", count)

	if err := db.DB.Preload("User").Preload("Address").Offset(offset).Limit(limit).Order("created_at DESC").Find(&orders).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Failed to fetch orders"})
		return
	}

	num := int(count) / (limit)
	if int(count)%limit != 0 {
		num = num + 1
	}
	fmt.Println("num", num)
	pagenumber := make([]int, 0)

	for i := 1; i <= num; i++ {
		pagenumber = append(pagenumber, i)
	}
	if len(pagenumber) == 0 {
		pagenumber = append(pagenumber, 1)
	}

	// Render the userorders.html template with the orders data
	c.HTML(http.StatusOK, "orders.html", gin.H{
		"Orders":      orders,
		"Pagenumber":  pagenumber,
		"Entries":     limit,
		"Currentpage": pgno,
	})
}

func UpdateOrderStatusHandler(c *gin.Context) {
	var updateStatusRequest struct {
		ID     uint   `json:"id"`
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&updateStatusRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order models.OrderProducts
	if err := db.DB.First(&order, updateStatusRequest.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	fmt.Println("upda", updateStatusRequest.Status, updateStatusRequest.ID)
	order.Status = updateStatusRequest.Status
	if err := db.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated successfully"})
}
func GetOrderStats(c *gin.Context) {
	report := c.Query("report")
	fmt.Println("report", report)

	var (
		orders         []models.OrderProducts
		count          int64
		totalAmount    float64
		totalDelivered int64
	)

	// Set the time range based on the report type
	var timeRange time.Time
	switch report {
	case "daily":
		timeRange = time.Now().AddDate(0, 0, -1)
	case "weekly":
		timeRange = time.Now().AddDate(0, 0, -7)
	case "monthly":
		timeRange = time.Now().AddDate(0, -1, 0)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "badrequest"})
		return
	}

	err := db.DB.Preload("Variant").Preload("Variant.Product").Preload("OrderDetails").
		Preload("OrderDetails.User").Where("status <> 'cancelled' AND created_at > ?", timeRange).Order("created_at DESC").Find(&orders).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching orders"})
		return
	}

	if err := db.DB.Debug().Table("orders").Where("status <> 'cancelled' AND created_at > ?", timeRange).Count(&count).Error; err != nil {
		fmt.Println("Error checking table count:", err)
		return
	}

	if err := db.DB.Debug().Table("order_products").Where("status = 'delivered' AND created_at > ?", timeRange).Count(&totalDelivered).Error; err != nil {
		fmt.Println("Error checking table count:", err)
		return
	}

	if count > 0 {
		if result := db.DB.Table("orders").Where("status <> 'cancelled' AND created_at > ?", timeRange).
			Select("SUM(Total) as total_amount").Scan(&totalAmount); result.Error != nil {
			fmt.Println("Error calculating sum:", result.Error)
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"OrdersReport":   orders,
		"TotalAmount":    totalAmount,
		"TotalOrders":    count,
		"TotalDelivered": totalDelivered,
	})
}

func ManageOrderHandler(c *gin.Context) {
	var orders []models.OrderProducts
	var OrderDetails models.Orders
	orderid := c.Query("id")
	// if err != nil {
	// 	fmt.Println("error", err)
	// 	// Handle the error appropriately, possibly return an error response.
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	// 	return
	// }
	fmt.Println("id", orderid)
	db.DB.Preload("Variant.Product.Images").Preload("Variant.Product").Preload("Variant").Preload("Image").Where("order_id=?", orderid).Find(&orders)
	db.DB.Preload("User").Preload("Address").Where("id=?", orderid).Find(&OrderDetails)
	fmt.Println("hiii")
	// c.JSON(http.StatusOK, orders)
	c.HTML(http.StatusOK, "manageorder.html", gin.H{
		// "Productvariants": ProductVariants,
		"Order":        orders,
		"OrderDetails": OrderDetails,
	})
}
