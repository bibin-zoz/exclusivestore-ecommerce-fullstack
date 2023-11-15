package middleware

import (
	"ecommercestore/helpers"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		token, err := c.Cookie("token")

		if err != nil && path != "/admin/login" {
			// If token is not present and the path is not /login, redirect to login
			log.Println("Token not present in cookie:", err)
			c.Redirect(http.StatusSeeOther, "/admin/login")
			c.AbortWithStatus(http.StatusSeeOther)
			return
		}

		if err == nil && path == "/admin/login" {
			// If token is present and the path is /login, redirect to home or another appropriate route
			log.Println("Admin is already logged in.")
			c.Redirect(http.StatusSeeOther, "/admin/home")
			c.AbortWithStatus(http.StatusSeeOther)
			return
		}

		if err != nil && path != "/admin/login" {
			// If token is not present and the path is not /login, redirect to login
			log.Println("Token not present in cookie:", err)
			c.Redirect(http.StatusSeeOther, "/admin/login")
			c.AbortWithStatus(http.StatusSeeOther)
			return
		}
		role, err := helpers.GetUserRoleFromToken(token)
		if role != "admin" && role != "" {
			// If token is not present and the path is not /login, redirect to login
			log.Println("Roll mismatch or not exist", err)
			if role == "user" {
				c.Redirect(http.StatusSeeOther, "/home")
			} else if role == "staff" {
				c.Redirect(http.StatusSeeOther, "/staff/home")
			}

			c.AbortWithStatus(http.StatusSeeOther)
			return
		}

		// Continue to the next middleware or handler
		c.Next()
	}
}
