package main

import (
	db "ecommercestore/database"
	handlers "ecommercestore/handlers"
	"ecommercestore/middleware"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	db.InitDB()

	//
	router.LoadHTMLGlob("templates/**/*.html")
	router.Static("/static", "./static")

	router.GET("/signup", middleware.IsLogin(), handlers.SignupHandler)
	router.POST("/signup", handlers.SignupPost)
	router.GET("/verify", middleware.IsLogin(), handlers.VerifyHandler)
	router.POST("/verify", handlers.VerifyPost)
	router.GET("/login", middleware.IsLogin(), handlers.LoginHandler)
	router.POST("/login", handlers.LoginPost)
	router.GET("/home", middleware.LoginAuth(), handlers.HomeHandler)
	router.GET("/logout", handlers.LogoutHandler)
	router.GET("/admin/logout", handlers.AdminLogoutHandler)
	router.GET("/product", handlers.ProductViewhandler)

	//user cart
	router.GET("/cart", middleware.LoginAuth(), handlers.GetCarthandler)
	router.POST("/cart", middleware.LoginAuth(), handlers.AddToCarthandler)
	router.PATCH("/cart", middleware.LoginAuth(), handlers.UpdateQuantityHandler)
	router.DELETE("/cart", middleware.LoginAuth(), handlers.DeleteCartHandler)

	//coupon cart
	router.GET("/couponvalidate", handlers.CouponValidatehandler)
	router.PATCH("/coupon", middleware.LoginAuth(), handlers.RemoveCouponHandler)

	//order
	router.POST("/razorpay/order", handlers.CreateRazorpayOrder)
	router.POST("/order", middleware.LoginAuth(), handlers.OrderPlacehandler)
	router.GET("/order", middleware.LoginAuth(), handlers.GetOrdershandler)
	router.GET("/trackorder", middleware.LoginAuth(), handlers.TrackOrderHandler)
	router.PUT("/order", middleware.LoginAuth(), handlers.ReturnOrderHandler)

	router.PATCH("/order", middleware.LoginAuth(), handlers.CancelOrderHandler)
	router.PATCH("/cancelitem", middleware.LoginAuth(), handlers.CancelProductHandler)

	//wallet
	router.GET("/wallet", middleware.LoginAuth(), handlers.WalletHandler)
	router.GET("/test", middleware.LoginAuth(), handlers.TestHandler)

	//checkout
	router.GET("/checkout", middleware.LoginAuth(), handlers.CheckOuthandler)

	//referall
	router.GET("/referalvalidate", handlers.ReferalValidatehandler)

	//admin
	router.GET("/admin/login", handlers.AdminLogin)
	router.POST("/admin/login", handlers.AdminLoginPost)
	router.GET("/admin/home", middleware.AdminAuth(), handlers.AdminHome)

	//sales report
	router.GET("/admin/salesreport", middleware.AdminAuth(), handlers.SalesReporthandler)

	//customers
	router.GET("/admin/customers", middleware.AdminAuth(), handlers.CustomerHandler)
	router.DELETE("/admin/customers", middleware.AdminAuth(), handlers.DeleteCustomerHandler)
	router.GET("/admin/customer/update-status", middleware.AdminAuth(), handlers.UpdateStatusHandler)

	//category
	router.GET("/admin/categories", middleware.AdminAuth(), handlers.Categoryhandler)
	router.PATCH("/admin/categories", middleware.AdminAuth(), handlers.UpdateCategoryStatus)
	router.POST("/admin/categories", middleware.AdminAuth(), handlers.CategoryPost)
	router.DELETE("/admin/categories", middleware.AdminAuth(), handlers.DeleteCategoryHandler)

	//categoryoffers
	router.GET("/admin/categoryoffers", middleware.AdminAuth(), handlers.CategoryOffershandler)
	router.POST("/admin/categoryoffers", middleware.AdminAuth(), handlers.AddCategoryOffershandler)
	router.DELETE("/admin/categoryoffers", middleware.AdminAuth(), handlers.DeleteCategoryOfferHandler)

	//sellers
	router.GET("/admin/sellers", middleware.AdminAuth(), handlers.SellersHandler)

	//Products
	router.GET("/admin/products", middleware.AdminAuth(), handlers.ProductsHandler)
	router.POST("/admin/product", middleware.AdminAuth(), handlers.AddProduct)
	router.PATCH("/admin/products", middleware.AdminAuth(), handlers.UpdateProductStatus)
	router.DELETE("/admin/products", middleware.AdminAuth(), handlers.DeleteProductHandler)

	//product edit
	router.GET("/admin/product", middleware.AdminAuth(), handlers.ProductDetailsHandler)
	router.PUT("/admin/product", middleware.AdminAuth(), handlers.ProductUpdateHandler)

	//orders

	router.GET("/admin/orders", middleware.AdminAuth(), handlers.UserOrdersHandler)
	router.GET("/admin/manageorder", middleware.AdminAuth(), handlers.ManageOrderHandler)
	router.PATCH("/admin/orders", middleware.AdminAuth(), handlers.UpdateOrderStatusHandler)

	router.GET("/admin/getOrderStats", middleware.AdminAuth(), handlers.GetOrderStats)

	//user dashboard
	router.GET("/userdashboard", middleware.LoginAuth(), handlers.UserDashboardHandler)

	//profile
	router.GET("/userprofile", middleware.LoginAuth(), handlers.GetUserProfileHandler)
	router.PUT("/userprofile", middleware.LoginAuth(), handlers.UpdateUserProfileHandler)
	router.PATCH("/userprofile", middleware.LoginAuth(), handlers.UpdatePasswordHandler)

	//address
	router.GET("/useraddress", middleware.LoginAuth(), handlers.UserAddressHandler)
	router.POST("/useraddress", middleware.LoginAuth(), handlers.NewAddressHandler)
	router.DELETE("/useraddress", middleware.LoginAuth(), handlers.DeleteAddressHandler)

	//coupon
	router.GET("/admin/coupons", middleware.AdminAuth(), handlers.CouponHandler)
	router.GET("/admin/coupon", middleware.AdminAuth(), handlers.CouponHandler)
	router.DELETE("/admin/coupons", middleware.AdminAuth(), handlers.DeleteCouponHandler)
	router.POST("/admin/coupon", middleware.AdminAuth(), handlers.AddCouponHandler)

	router.Run(":8080")

}
