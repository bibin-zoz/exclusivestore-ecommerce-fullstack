package main

import (
	db "ecommercestore/database"
	handlers "ecommercestore/handlers/user_handler"
	"ecommercestore/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	db.InitDB()
	routes.AdminRoutes(router)
	routes.UserRoutes(router)

	//
	router.LoadHTMLGlob("templates/**/*.html")
	router.Static("/static", "./static")
	router.NoRoute(handlers.PageNotfoundHandler)

	router.Run(":8080")

}
