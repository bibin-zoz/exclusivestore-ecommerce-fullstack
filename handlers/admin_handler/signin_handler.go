package handlers

import (
	db "ecommercestore/database"
	"ecommercestore/helpers"
	"ecommercestore/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func AdminLoginPost(c *gin.Context) {

	Newmail := c.Request.FormValue("email")
	Newpassword := c.Request.FormValue("password")
	var compare models.Compare
	var data models.Invalid

	if Newmail == "" {
		data.EmailError = "Email should not be empty"
		c.HTML(http.StatusOK, "login.html", data)
		return
	}
	if Newpassword == "" {
		data.PasswordError = "password should not be empty"
		c.HTML(http.StatusBadRequest, "login.html", data)
		return
	}
	if err := db.DB.Raw("SELECT password, username,role,status FROM users WHERE email=$1", Newmail).Scan(&compare).Error; err != nil {
		fmt.Println(err)
		data.EmailError = "An error occurred while querying the database"
		c.HTML(http.StatusInternalServerError, "login.html", data)
		return
	}

	// Check if no user is found
	var count int64

	if result := db.DB.Model(&models.User{}).Where("email = ? ", Newmail).Count(&count); result.Error != nil || count == 0 {
		data.EmailError = "User not found! Re-check the Mailid"
		c.HTML(http.StatusBadRequest, "adminlogin.html", data)
		return
	}
	err := helpers.VerifyPassword(Newpassword, compare.Password)
	if err != nil {
		data.PasswordError = "check password again"
		c.HTML(http.StatusBadRequest, "adminlogin.html", data)
		return
	}
	if compare.Role == "user" {
		data.RoleError = "click here for admin login -->"
		c.HTML(http.StatusBadRequest, "adminlogin.html", data)
		return
	}
	if compare.Status != "active" {
		data.StatusError = "User is blocked"
		c.HTML(http.StatusBadRequest, "adminlogin.html", data)
		return
	}
	claims := models.Claims{
		ID:       compare.ID,
		Username: compare.Username,
		Email:    compare.Email,
		Role:     compare.Role,
		Status:   compare.Status,
	}

	accessToken, err := helpers.GenerateAccessToken(claims)
	if err != nil {
		fmt.Println("Error generating access token:", err)

		return
	}

	refreshToken, err := helpers.GenerateRefreshToken(claims)
	if err != nil {
		fmt.Println("Error generating refresh token:", err)

		return
	}

	UserLoginDetails := &models.TokenUser{
		// Users:        claims,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	userDetailsJSON := helpers.CreateJson(UserLoginDetails)

	c.SetCookie("adminAuth", string(userDetailsJSON), 0, "/admin", "exclusivestore.xyz", true, true)

	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Expires", "0")

	// Redirect to home only after successful token generation
	c.Redirect(http.StatusFound, "/admin/home")
}
func AdminLogoutHandler(c *gin.Context) {
	fmt.Println("admin logout")

	// Clear the adminAuth cookie

	c.SetCookie("adminAuth", "", -1, "/admin", "exclusivestore.xyz", false, true)

	// Redirect to the login page
	c.Redirect(http.StatusSeeOther, "/admin/login")
}

func AdminLogin(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Expires", "0")
	_, err := c.Cookie("adminAuth")
	if err == nil {
		c.Redirect(http.StatusSeeOther, "/admin/home")
		c.AbortWithStatus(http.StatusSeeOther)
		return
	}

	c.HTML(http.StatusOK, "adminlogin.html", nil)

}

// admin
func AdminHome(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Expires", "0")

	// var orders models.OrderProducts
	var productssold int64
	result := db.DB.Table("order_products").Where("status='delivered'").Count(&productssold)
	if result.Error != nil {
		fmt.Println("error")
	}
	var netProfit int64
	result = db.DB.Table("order_products").Where("status='delivered'").Select("SUM(total)").Scan(&netProfit)
	if result.Error != nil {
		fmt.Println("error")
	}
	var newCustomers int64
	timeRange := time.Now().AddDate(0, -1, 0)
	result = db.DB.Table("users").Where("created_at>?", timeRange).Count(&newCustomers)
	if result.Error != nil {
		fmt.Println("error")
	}

	c.HTML(http.StatusOK, "adminhome.html", gin.H{
		"ProductsSold": productssold,
		"NewCustomers": newCustomers,
		"NetProfit":    netProfit,
	})

}
