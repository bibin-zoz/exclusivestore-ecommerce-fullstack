package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func IsLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		_, err := c.Cookie("adminAuth")
		if err == nil && path == "/admin/login" {
			c.Redirect(http.StatusSeeOther, "/admin/home")
			c.AbortWithStatus(http.StatusSeeOther)
			return
		}
		_, err = c.Cookie("auth")
		if err == nil {
			c.Redirect(http.StatusSeeOther, "/home")
			c.AbortWithStatus(http.StatusSeeOther)
			return
		}
		c.Next()
	}
}
