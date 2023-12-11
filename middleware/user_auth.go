package middleware

import (
	"ecommercestore/helpers"
	"ecommercestore/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func LoginAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the auth cookie

		authCookie, err := c.Cookie("auth")
		if err != nil {
			fmt.Println("err not present", err)
			c.Redirect(http.StatusSeeOther, "/login")
			c.AbortWithStatus(http.StatusSeeOther)
			return
		}

		// Decode the JSON content of the auth cookie
		var token models.TokenUser
		err = json.NewDecoder(strings.NewReader(authCookie)).Decode(&token)
		if err != nil {
			// Redirect to login if there's an error decoding the auth cookie
			c.Redirect(http.StatusSeeOther, "/login")
			c.AbortWithStatus(http.StatusSeeOther)
			return
		}
		// fmt.Println("token", token.AccessToken)
		// fmt.Println("token", token.RefreshToken)

		claims, err := helpers.ParseToken(token.AccessToken)
		if err != nil {
			log.Println("Error parsing token:", err)
			c.Redirect(http.StatusSeeOther, "/login")
			c.AbortWithStatus(http.StatusSeeOther)
			return
		}
		refreshTokenClaims, err := helpers.ParseToken(token.RefreshToken)
		if err != nil {
			log.Println("Error parsing token:", err)
			c.Redirect(http.StatusSeeOther, "/login")
			c.AbortWithStatus(http.StatusSeeOther)
			return
		}
		if time.Now().Unix() > claims.StandardClaims.ExpiresAt {
			if time.Now().Unix() > refreshTokenClaims.StandardClaims.ExpiresAt {
				log.Println("Token has expired")
				c.Redirect(http.StatusSeeOther, "/login")
				c.AbortWithStatus(http.StatusSeeOther)
				return

			} else {
				c.SetCookie("auth", "", -1, "/", "localhost", true, true)
				accessToken, err := helpers.GenerateAccessToken(*claims)
				if err != nil {
					log.Println("Error creating access token token:", err)
					c.Redirect(http.StatusSeeOther, "/login")
					c.AbortWithStatus(http.StatusSeeOther)
					return
				}
				UserAuth := &models.TokenUser{
					// Users:        claims,
					AccessToken:  accessToken,
					RefreshToken: token.RefreshToken,
				}
				userDetailsJSON := helpers.CreateJson(UserAuth)

				c.SetCookie("auth", string(userDetailsJSON), 0, "/", "localhost", true, true)
			}

		}
		if err != nil {
			log.Println("Token not present in cookie:", err)
			c.Redirect(http.StatusSeeOther, "/login")
			c.AbortWithStatus(http.StatusSeeOther)
			return
		}
		if claims.Role != "user" && claims.Role != "" {
			log.Println("User not logined In")
			c.Redirect(http.StatusSeeOther, "/login")
			c.AbortWithStatus(http.StatusSeeOther)
			return
		}
		if claims.Status != "active" {
			log.Println("User is blocked")
			c.Redirect(http.StatusSeeOther, "/login")
			c.AbortWithStatus(http.StatusSeeOther)
			return
		}

		c.Next()
	}
}
