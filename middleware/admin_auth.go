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

			log.Println("Token not present in cookie:", err)
			c.Redirect(http.StatusSeeOther, "/admin/login")
			c.AbortWithStatus(http.StatusSeeOther)
			return
		}

		if err == nil && path == "/admin/login" {

			log.Println("Admin is already logged in.")
			c.Redirect(http.StatusSeeOther, "/admin/home")
			c.AbortWithStatus(http.StatusSeeOther)
			return
		}

		if err != nil && path != "/admin/login" {

			log.Println("Token not present in cookie:", err)
			c.Redirect(http.StatusSeeOther, "/admin/login")
			c.AbortWithStatus(http.StatusSeeOther)
			return
		}
		role, err := helpers.GetUserRoleFromToken(token)
		if role != "admin" && role != "" {

			log.Println("Roll mismatch or not exist", err)
			if role == "user" {
				c.Redirect(http.StatusSeeOther, "/home")
			} else if role == "staff" {
				c.Redirect(http.StatusSeeOther, "/staff/home")
			}

			c.AbortWithStatus(http.StatusSeeOther)
			return
		}

		c.Next()
	}
}
