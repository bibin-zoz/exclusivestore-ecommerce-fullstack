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

	// router.LoadHTMLGlob("templates/*.html")
	// router.LoadHTMLGlob("templates/*")
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
	router.GET("/product", handlers.ProductViewhandler)

	//user cart
	router.GET("/cart", middleware.LoginAuth(), handlers.GetCarthandler)
	router.POST("/cart", middleware.LoginAuth(), handlers.AddToCarthandler)

	//admin
	router.GET("/admin/login", middleware.IsLogin(), handlers.AdminLogin)
	router.POST("/admin/login", handlers.AdminLoginPost)
	router.GET("/admin/home", middleware.AdminAuth(), handlers.AdminHome)

	//customers
	router.GET("/admin/customers", middleware.AdminAuth(), handlers.CustomerHandler)
	router.DELETE("/admin/customers", middleware.AdminAuth(), handlers.DeleteCustomerHandler)
	router.GET("/admin/update-status", middleware.AdminAuth(), handlers.UpdateStatusHandler)

	//category
	router.GET("/admin/categories", middleware.AdminAuth(), handlers.Categoryhandler)
	router.PATCH("/admin/categories/update-status", middleware.AdminAuth(), handlers.UpdateCategoryStatus)
	router.POST("/admin/categories", middleware.AdminAuth(), handlers.CategoryPost)
	router.DELETE("/admin/categories", middleware.AdminAuth(), handlers.DeleteCategoryHandler)

	//sellers
	router.GET("/admin/sellers", middleware.AdminAuth(), handlers.SellersHandler)

	//Products
	router.GET("/admin/products", middleware.AdminAuth(), handlers.ProductsHandler)
	router.POST("/admin/product", middleware.AdminAuth(), handlers.AddProduct)
	router.GET("/admin/Products/update-status", middleware.AdminAuth(), handlers.UpdateProductStatus)
	router.DELETE("/admin/products", middleware.AdminAuth(), handlers.DeleteProductHandler)

	//product edit
	router.GET("/admin/product", middleware.AdminAuth(), handlers.ProductDetailsHandler)
	router.PUT("/admin/product", middleware.AdminAuth(), handlers.ProductUpdateHandler)

	router.Run(":8080")

}
