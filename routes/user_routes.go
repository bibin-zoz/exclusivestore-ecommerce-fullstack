package routes

import (
	handlers "ecommercestore/handlers/user_handler"
	"ecommercestore/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	userGroup := r.Group("/")

	// Public routes
	{
		// Authentication routes
		userGroup.GET("/signup", middleware.IsLogin(), handlers.SignupHandler)
		userGroup.POST("/signup", handlers.SignupPost)
		//login
		userGroup.GET("/login", middleware.IsLogin(), handlers.LoginHandler)
		userGroup.POST("/login", handlers.LoginPost)
		//verify
		userGroup.GET("/verify", middleware.IsLogin(), handlers.VerifyHandler)
		userGroup.POST("/verify", handlers.VerifyPost)

		// Home and product routes
		userGroup.GET("/home", middleware.LoginAuth(), handlers.HomeHandler)
		userGroup.GET("/", handlers.HomeHandler)
		//product
		userGroup.GET("/product", handlers.ProductViewhandler)
		userGroup.GET("/shop-products", middleware.LoginAuth(), handlers.GetProductsHandler)
		userGroup.POST("/shop-products", middleware.LoginAuth(), handlers.FilterProductshandler)

		// Coupon routes
		userGroup.GET("/couponvalidate", handlers.CouponValidatehandler)
		userGroup.PATCH("/coupon", middleware.LoginAuth(), handlers.RemoveCouponHandler)

		// payment routes
		userGroup.POST("/onlinepay", handlers.CreateRazorpayOrder)

		//order
		userGroup.GET("/order", middleware.LoginAuth(), handlers.GetOrdershandler)
		userGroup.POST("/order", middleware.LoginAuth(), handlers.OrderPlacehandler)
		userGroup.PUT("/order", middleware.LoginAuth(), handlers.ReturnOrderHandler)
		userGroup.PATCH("/order", middleware.LoginAuth(), handlers.CancelOrderHandler)
		userGroup.PATCH("/cancelitem", middleware.LoginAuth(), handlers.CancelProductHandler)
		//track
		userGroup.GET("/trackorder", middleware.LoginAuth(), handlers.TrackOrderHandler)
		//invoice
		userGroup.GET("/download", handlers.GeneratePDFHandler)

		// Wallet and checkout routes
		userGroup.POST("/wallet", handlers.WalletOrderhandler)
		userGroup.GET("/wallet", middleware.LoginAuth(), handlers.WalletHandler)

		userGroup.GET("/checkout", middleware.LoginAuth(), handlers.CheckOuthandler)
		//referal
		userGroup.GET("/referalvalidate", handlers.ReferalValidatehandler)

		// dashboard routes

		userGroup.GET("/userdashboard", middleware.LoginAuth(), handlers.UserDashboardHandler)

		// Profile and address routes
		userGroup.GET("/userprofile", middleware.LoginAuth(), handlers.GetUserProfileHandler)
		userGroup.PUT("/userprofile", middleware.LoginAuth(), handlers.UpdateUserProfileHandler)
		userGroup.PATCH("/userprofile", middleware.LoginAuth(), handlers.UpdatePasswordHandler)
		//address
		userGroup.GET("/useraddress", middleware.LoginAuth(), handlers.UserAddressHandler)
		userGroup.POST("/useraddress", middleware.LoginAuth(), handlers.NewAddressHandler)
		userGroup.DELETE("/useraddress", middleware.LoginAuth(), handlers.DeleteAddressHandler)

		// userGroup.GET("/test", middleware.LoginAuth(), handlers.TestHandler)
	}
}
