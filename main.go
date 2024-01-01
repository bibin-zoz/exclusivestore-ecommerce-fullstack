package main

import (
	db "ecommercestore/database"
	adminhandlers "ecommercestore/handlers/admin_handler"
	handlers "ecommercestore/handlers/user_handler"
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
	router.GET("/admin/logout", adminhandlers.AdminLogoutHandler)
	router.GET("/product", handlers.ProductViewhandler)

	//forgotpass
	router.GET("/forgotpass", middleware.IsLogin(), handlers.ForgotPasswordHandler)
	router.POST("/forgotpass", middleware.IsLogin(), handlers.ForgotPasswordPostHandler)
	router.GET("/emailverify", middleware.IsLogin(), handlers.EmailVerify)
	router.POST("/emailverify", handlers.EmailVerifyPost)

	//product
	router.GET("/shop-products", handlers.GetProductsHandler)
	router.POST("/shop-products", handlers.FilterProductshandler)
	// router.GET("/products", handlers.GetProductsHandler)

	//user cart
	router.GET("/cart", middleware.LoginAuth(), handlers.GetCarthandler)
	router.POST("/cart", middleware.LoginAuth(), handlers.AddToCarthandler)
	router.PATCH("/cart", middleware.LoginAuth(), handlers.UpdateQuantityHandler)
	router.DELETE("/cart", middleware.LoginAuth(), handlers.DeleteCartHandler)

	//coupon cart
	router.GET("/couponvalidate", handlers.CouponValidatehandler)
	router.PATCH("/coupon", middleware.LoginAuth(), handlers.RemoveCouponHandler)

	//order
	router.POST("/onlinepay", handlers.CreateRazorpayOrder)
	router.POST("/wallet", handlers.WalletOrderhandler)
	router.POST("/order", middleware.LoginAuth(), handlers.OrderPlacehandler)
	router.GET("/order", middleware.LoginAuth(), handlers.GetOrdershandler)
	router.GET("/trackorder", middleware.LoginAuth(), handlers.TrackOrderHandler)
	router.PUT("/order", middleware.LoginAuth(), handlers.ReturnOrderHandler)
	router.GET("/download", handlers.GeneratePDFHandler)

	router.PATCH("/order", middleware.LoginAuth(), handlers.CancelOrderHandler)
	router.PATCH("/cancelitem", middleware.LoginAuth(), handlers.CancelProductHandler)

	//wallet
	router.GET("/wallet", middleware.LoginAuth(), handlers.WalletHandler)
	router.GET("/test", middleware.LoginAuth(), handlers.TestHandler)

	//checkout
	router.GET("/checkout", middleware.LoginAuth(), handlers.CheckOuthandler)

	//referal
	router.GET("/referalvalidate", handlers.ReferalValidatehandler)

	//admin
	router.GET("/admin/login", adminhandlers.AdminLogin)
	router.POST("/admin/login", adminhandlers.AdminLoginPost)
	router.GET("/admin/home", middleware.AdminAuth(), adminhandlers.AdminHome)

	//sales report
	router.GET("/admin/salesreport", middleware.AdminAuth(), adminhandlers.SalesReporthandler)

	//customers
	router.GET("/admin/customers", middleware.AdminAuth(), adminhandlers.CustomerHandler)
	router.DELETE("/admin/customers", middleware.AdminAuth(), adminhandlers.DeleteCustomerHandler)
	router.GET("/admin/customer/update-status", middleware.AdminAuth(), adminhandlers.UpdateStatusHandler)

	//category
	router.GET("/admin/categories", middleware.AdminAuth(), adminhandlers.Categoryhandler)
	router.PATCH("/admin/categories", middleware.AdminAuth(), adminhandlers.UpdateCategoryStatus)
	router.POST("/admin/categories", middleware.AdminAuth(), adminhandlers.CategoryPost)
	router.DELETE("/admin/categories", middleware.AdminAuth(), adminhandlers.DeleteCategoryHandler)

	//categoryoffers
	router.GET("/admin/categoryoffers", middleware.AdminAuth(), handlers.CategoryOffershandler)
	router.POST("/admin/categoryoffers", middleware.AdminAuth(), handlers.AddCategoryOffershandler)
	router.DELETE("/admin/categoryoffers", middleware.AdminAuth(), handlers.DeleteCategoryOfferHandler)

	//Products
	router.GET("/admin/products", middleware.AdminAuth(), adminhandlers.ProductsHandler)
	router.POST("/admin/product", middleware.AdminAuth(), adminhandlers.AddProduct)
	router.PATCH("/admin/products", middleware.AdminAuth(), adminhandlers.UpdateProductStatus)
	router.DELETE("/admin/products", middleware.AdminAuth(), adminhandlers.DeleteProductHandler)

	//product edit
	router.GET("/admin/product", middleware.AdminAuth(), adminhandlers.ProductDetailsHandler)
	router.PUT("/admin/product", middleware.AdminAuth(), adminhandlers.ProductUpdateHandler)

	//product offers
	router.GET("/admin/productoffers", middleware.AdminAuth(), handlers.ProductOffersHandler)

	//orders

	router.GET("/admin/orders", middleware.AdminAuth(), adminhandlers.UserOrdersHandler)
	router.GET("/admin/manageorder", middleware.AdminAuth(), adminhandlers.ManageOrderHandler)
	router.PATCH("/admin/orders", middleware.AdminAuth(), adminhandlers.UpdateOrderStatusHandler)

	router.GET("/admin/getOrderStats", middleware.AdminAuth(), adminhandlers.GetOrderStats)

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

	router.GET("/admin/coupon", middleware.AdminAuth(), handlers.CouponHandler)
	router.DELETE("/admin/coupon", middleware.AdminAuth(), handlers.DeleteCouponHandler)
	router.POST("/admin/coupon", middleware.AdminAuth(), handlers.AddCouponHandler)

	router.Run(":8080")

}
