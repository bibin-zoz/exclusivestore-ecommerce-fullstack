package routes

import (
	adminhandlers "ecommercestore/handlers/admin_handler"
	handlers "ecommercestore/handlers/user_handler"
	"ecommercestore/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.Engine) {
	r.GET("/admin/logout", adminhandlers.AdminLogoutHandler)
	r.GET("/admin/login", middleware.IsLogin(), adminhandlers.AdminLogin)
	r.POST("/admin/login", middleware.IsLogin(), adminhandlers.AdminLoginPost)

	adminGroup := r.Group("/admin")
	adminGroup.Use(middleware.AdminAuth())

	// Admin Home
	adminGroup.GET("", adminhandlers.AdminHome)

	adminGroup.GET("/home", adminhandlers.AdminHome)

	// Sales Report
	adminGroup.GET("/salesreport", adminhandlers.SalesReporthandler)
	adminGroup.GET("/downloadsalesreport", adminhandlers.SalesReportDownloadhandler)
	adminGroup.POST("/salesreport", adminhandlers.SalesReportDownloadhandler)

	// Customers
	adminGroup.GET("/customers", adminhandlers.CustomerHandler)
	adminGroup.DELETE("/customers", adminhandlers.DeleteCustomerHandler)
	adminGroup.GET("/customer/update-status", adminhandlers.UpdateStatusHandler)

	// Categories
	adminGroup.GET("/categories", adminhandlers.Categoryhandler)
	adminGroup.PATCH("/categories", adminhandlers.UpdateCategoryStatus)
	adminGroup.POST("/categories", adminhandlers.CategoryPost)
	adminGroup.DELETE("/categories", adminhandlers.DeleteCategoryHandler)

	// Category Offers
	adminGroup.GET("/categoryoffers", handlers.CategoryOffershandler)
	adminGroup.POST("/categoryoffers", handlers.AddCategoryOffershandler)
	adminGroup.DELETE("/categoryoffers", handlers.DeleteCategoryOfferHandler)

	// Products
	adminGroup.GET("/products", adminhandlers.ProductsHandler)
	adminGroup.POST("/product", adminhandlers.AddProduct)
	adminGroup.PATCH("/products", adminhandlers.UpdateProductStatus)
	adminGroup.DELETE("/products", adminhandlers.DeleteProductHandler)

	// Product Edit
	adminGroup.GET("/product", adminhandlers.ProductDetailsHandler)
	adminGroup.PUT("/product", adminhandlers.ProductUpdateHandler)

	// Product Offers
	adminGroup.GET("/productoffers", handlers.ProductOffersHandler)

	// Orders
	adminGroup.GET("/orders", adminhandlers.UserOrdersHandler)
	adminGroup.GET("/manageorder", adminhandlers.ManageOrderHandler)
	adminGroup.PATCH("/orders", adminhandlers.UpdateOrderStatusHandler)

	adminGroup.GET("/getOrderStats", adminhandlers.GetOrderStats)

	adminGroup.GET("/coupon", handlers.CouponHandler)
	adminGroup.DELETE("/coupon", handlers.DeleteCouponHandler)
	adminGroup.POST("/coupon", handlers.AddCouponHandler)

}
