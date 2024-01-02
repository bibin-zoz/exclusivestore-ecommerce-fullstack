package handlers

import (
	"bytes"
	db "ecommercestore/database"
	"ecommercestore/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

type DeleteRequest struct {
	ID int `json:"id"`
}

func SalesReporthandler(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Expires", "0")

	c.HTML(http.StatusOK, "salesreport.html", nil)

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
func SalesReportDownloadhandler(c *gin.Context) {
	// Generate PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 30)
	pdf.SetTextColor(31, 73, 125)
	pdf.Cell(0, 20, "Sales Report")
	pdf.Ln(20)

	// Fetch monthly, weekly, and daily sales report
	var (
		monthlyOrders         []models.OrderProducts
		monthlyCount          int64
		monthlyTotalAmount    float64
		monthlyTotalDelivered int64

		weeklyOrders         []models.OrderProducts
		weeklyCount          int64
		weeklyTotalAmount    float64
		weeklyTotalDelivered int64

		dailyOrders         []models.OrderProducts
		dailyCount          int64
		dailyTotalAmount    float64
		dailyTotalDelivered int64
	)

	// Set the time range for the monthly, weekly, and daily reports
	now := time.Now()
	monthlyTimeRange := now.AddDate(0, -1, 0)
	weeklyTimeRange := now.AddDate(0, 0, -7)
	dailyTimeRange := now.AddDate(0, 0, -1)

	// Fetch monthly report
	err := db.DB.Preload("Variant").Preload("Variant.Product").Preload("OrderDetails").
		Preload("OrderDetails.User").Where("status <> 'cancelled' AND created_at > ?", monthlyTimeRange).Order("created_at DESC").Find(&monthlyOrders).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching monthly orders"})
		return
	}

	if err := db.DB.Debug().Table("orders").Where("status <> 'cancelled' AND created_at > ?", monthlyTimeRange).Count(&monthlyCount).Error; err != nil {
		fmt.Println("Error checking monthly table count:", err)
		return
	}

	if err := db.DB.Debug().Table("order_products").Where("status = 'delivered' AND created_at > ?", monthlyTimeRange).Count(&monthlyTotalDelivered).Error; err != nil {
		fmt.Println("Error checking monthly table count:", err)
		return
	}

	if monthlyCount > 0 {
		if result := db.DB.Table("orders").Where("status <> 'cancelled' AND created_at > ?", monthlyTimeRange).
			Select("SUM(Total) as total_amount").Scan(&monthlyTotalAmount); result.Error != nil {
			fmt.Println("Error calculating monthly sum:", result.Error)
			return
		}
	}

	// Fetch weekly report
	err = db.DB.Preload("Variant").Preload("Variant.Product").Preload("OrderDetails").
		Preload("OrderDetails.User").Where("status <> 'cancelled' AND created_at > ?", weeklyTimeRange).Order("created_at DESC").Find(&weeklyOrders).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching weekly orders"})
		return
	}

	if err := db.DB.Debug().Table("orders").Where("status <> 'cancelled' AND created_at > ?", weeklyTimeRange).Count(&weeklyCount).Error; err != nil {
		fmt.Println("Error checking weekly table count:", err)
		return
	}

	if err := db.DB.Debug().Table("order_products").Where("status = 'delivered' AND created_at > ?", weeklyTimeRange).Count(&weeklyTotalDelivered).Error; err != nil {
		fmt.Println("Error checking weekly table count:", err)
		return
	}

	if weeklyCount > 0 {
		if result := db.DB.Table("orders").Where("status <> 'cancelled' AND created_at > ?", weeklyTimeRange).
			Select("SUM(Total) as total_amount").Scan(&weeklyTotalAmount); result.Error != nil {
			fmt.Println("Error calculating weekly sum:", result.Error)
			return
		}
	}

	// Fetch daily report
	err = db.DB.Preload("Variant").Preload("Variant.Product").Preload("OrderDetails").
		Preload("OrderDetails.User").Where("status <> 'cancelled' AND created_at > ?", dailyTimeRange).Order("created_at DESC").Find(&dailyOrders).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching daily orders"})
		return
	}

	if err := db.DB.Debug().Table("orders").Where("status <> 'cancelled' AND created_at > ?", dailyTimeRange).Count(&dailyCount).Error; err != nil {
		fmt.Println("Error checking daily table count:", err)
		return
	}

	if err := db.DB.Debug().Table("order_products").Where("status = 'delivered' AND created_at > ?", dailyTimeRange).Count(&dailyTotalDelivered).Error; err != nil {
		fmt.Println("Error checking daily table count:", err)
		return
	}

	if dailyCount > 0 {
		if result := db.DB.Table("orders").Where("status <> 'cancelled' AND created_at > ?", dailyTimeRange).
			Select("SUM(Total) as total_amount").Scan(&dailyTotalAmount); result.Error != nil {
			fmt.Println("Error calculating daily sum:", result.Error)
			return
		}
	}

	// Add the fetched information to the PDF
	pdf.Ln(10)

	// Monthly Report
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Monthly Report")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 14)
	pdf.Cell(0, 10, fmt.Sprintf("Total Orders: %d", monthlyCount))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Total Delivered Orders: %d", monthlyTotalDelivered))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Total Sales: RS.%.2f", monthlyTotalAmount))
	pdf.Ln(20)

	// Weekly Report
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Weekly Report")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 14)
	pdf.Cell(0, 10, fmt.Sprintf("Total Orders: %d", weeklyCount))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Total Delivered Orders: %d", weeklyTotalDelivered))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Total Sales: Rs.%.2f", weeklyTotalAmount))
	pdf.Ln(20)

	// Daily Report
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Daily Report")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 14)
	pdf.Cell(0, 10, fmt.Sprintf("Total Orders: %d", dailyCount))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Total Delivered Orders: %d", dailyTotalDelivered))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Total Sales: Rs.%.2f", dailyTotalAmount))

	// Serve the PDF as a response
	var buf bytes.Buffer
	pdf.Output(&buf)
	pdfOutput := buf.Bytes()
	c.Data(http.StatusOK, "application/pdf", pdfOutput)
	fmt.Println("PDF generated successfully")
}
