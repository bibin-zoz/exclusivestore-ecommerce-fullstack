package middleware

import (
	"ecommercestore/helpers"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func LoginAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// Skip token check for the "/login" path
		if path == "/login" {
			c.Next()
			return
		}

		Token, err := c.Cookie("token")

		if err != nil {
			log.Println("Token not present in cookie:", err)
			c.Redirect(http.StatusSeeOther, "/login")
			c.AbortWithStatus(http.StatusSeeOther)
			return
		}

		fmt.Println("claims", Token)
		claims, err := helpers.ParseToken(Token)
		if err != nil {
			log.Println("Error parsing token:", err)
			c.Redirect(http.StatusSeeOther, "/login")
			c.AbortWithStatus(http.StatusSeeOther)
			return
		}

		if time.Now().Unix() > claims.StandardClaims.ExpiresAt {
			log.Println("Token has expired")
			c.Redirect(http.StatusSeeOther, "/login")
			c.AbortWithStatus(http.StatusSeeOther)
			return
		}

		if err != nil {
			log.Println("Token not present in cookie:", err)
			c.Redirect(http.StatusSeeOther, "/login")
			c.AbortWithStatus(http.StatusSeeOther)
			return
		}

		role, _, err := helpers.GetUserRoleFromToken(Token)
		if role != "user" && role != "" {
			log.Println("Roll mismatch or not exist", err)
			if role == "admin" {
				c.Redirect(http.StatusSeeOther, "/admin/home")
			} else if role == "staff" {
				c.Redirect(http.StatusSeeOther, "/staff/home")
			}

			c.AbortWithStatus(http.StatusSeeOther)
			return
		}

		c.Next()
	}
}
