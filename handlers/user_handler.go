package handlers

import (
	db "ecommercestore/database"
	"ecommercestore/models"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"ecommercestore/helpers"

	"github.com/gin-gonic/gin"
)

type UserDetails struct {
	Username        string
	Email           string
	Number          string
	Password        string
	ConfirmPassword string
}

var user UserDetails

func LoginHandler(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Expires", "0")
	var data models.Invalid
	data.LoginStatus = false
	c.HTML(http.StatusOK, "login.html", data)

}

func LoginPost(c *gin.Context) {
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
		data.PasswordError = "Password should not be empty"
		c.HTML(http.StatusBadRequest, "login.html", data)
		return
	}
	if err := db.DB.Raw("SELECT ID, password, username,email, role, status FROM users WHERE email=$1", Newmail).Scan(&compare).Error; err != nil {
		fmt.Println("Error querying the database:", err)
		data.EmailError = "An error occurred while querying the database"
		c.HTML(http.StatusInternalServerError, "login.html", data)
		return
	}

	var count int64
	if result := db.DB.Model(&models.User{}).Where("email = ?", Newmail).Count(&count); result.Error != nil || count == 0 {
		data.EmailError = "User not found! Re-check the Mailid"
		c.HTML(http.StatusBadRequest, "login.html", data)
		return
	}
	if compare.Password != Newpassword {
		data.PasswordError = "Check password again"
		c.HTML(http.StatusBadRequest, "login.html", data)
		return
	}
	if compare.Role != "user" {
		data.RoleError = "Click here for admin login -->"
		c.HTML(http.StatusBadRequest, "login.html", data)
		return
	}
	if compare.Status != "active" {
		data.StatusError = "User is blocked"
		c.HTML(http.StatusBadRequest, "login.html", data)
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
		// Handle the error (e.g., return an error response)
		return
	}

	refreshToken, err := helpers.GenerateRefreshToken(claims)
	if err != nil {
		fmt.Println("Error generating refresh token:", err)
		// Handle the error (e.g., return an error response)
		return
	}

	UserLoginDetails := &models.TokenUser{
		// Users:        claims,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	userDetailsJSON := helpers.CreateJson(UserLoginDetails)

	if claims.Role == "admin" {
		c.SetCookie("adminAuth", string(userDetailsJSON), 0, "/", "localhost", true, true)

	} else {
		c.SetCookie("auth", string(userDetailsJSON), 0, "/", "localhost", true, true)

	}

	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Expires", "0")

	// Redirect to home only after successful token generation
	c.Redirect(http.StatusFound, "/home")
}

// ... (your other functions)

func SignupHandler(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Expires", "0")

	c.HTML(http.StatusOK, "signup.html", nil)

}

func SignupPost(c *gin.Context) {

	var errors models.Invalid
	user = UserDetails{
		Username:        c.Request.FormValue("username"),
		Email:           c.Request.FormValue("email"),
		Number:          c.Request.FormValue("number"),
		Password:        c.Request.FormValue("password"),
		ConfirmPassword: c.Request.FormValue("confirmPassword"),
	}

	if user.Username == "" {
		errors.NameError = "Name should not be empty"
		c.HTML(http.StatusBadRequest, "signup.html", errors)
		return
	}
	var usernameCount int64
	if err := db.DB.Table("users").Where("username = ?", user.Username).Count(&usernameCount).Error; err != nil {
		fmt.Println(err)
		c.HTML(http.StatusBadRequest, "signup.html", nil)
		return
	}

	if usernameCount > 0 {
		errors.NameError = "Username already exists"
		c.HTML(http.StatusBadRequest, "signup.html", errors)
		return
	}
	if user.Email == "" {
		errors.EmailError = "Email should not be empty"
		c.HTML(http.StatusBadRequest, "signup.html", errors)
		return
	}
	var emailCount int64
	if err := db.DB.Table("users").Where("email = ?", user.Email).Count(&emailCount).Error; err != nil {
		fmt.Println(err)
		c.HTML(http.StatusBadRequest, "signup.html", nil)
		return
	}

	if emailCount > 0 {
		errors.EmailError = "Email already exists"
		c.HTML(http.StatusBadRequest, "signup.html", errors)
		return
	}

	pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(pattern)
	if !regex.MatchString(user.Email) {
		errors.EmailError = "Email not in the correct format"
		c.HTML(http.StatusBadRequest, "signup.html", errors)
		return
	}

	if user.Number == "" {
		errors.NumberError = "Number should not be empty"
		c.HTML(http.StatusBadRequest, "signup.html", errors)
		return
	}

	pattern = `^[0-9]{10}$`
	regex = regexp.MustCompile(pattern)
	if !regex.MatchString(user.Number) {
		errors.NumberError = "Invalid Mobile Number"
		c.HTML(http.StatusBadRequest, "signup.html", errors)
		return
	}

	if user.Password == "" {
		errors.PasswordError = "Password should not be empty"
		c.HTML(http.StatusBadRequest, "signup.html", errors)
		return
	}
	if user.Password != user.ConfirmPassword {
		errors.PasswordError = "Passwords do not match"
		c.HTML(http.StatusBadRequest, "signup.html", errors)
		return
	}

	var count int64
	if err := db.DB.Table("users").Where("email = ?", user.Email).Count(&count).Error; err != nil {
		fmt.Println(err)
		c.HTML(http.StatusOK, "signup.html", nil)
		return
	}

	if count > 0 {
		errors.EmailError = "User already exists"
		c.HTML(http.StatusBadRequest, "signup.html", errors)
		return
	}

	c.Redirect(http.StatusFound, "/verify")
}

var lastOTPSendTime time.Time

func VerifyHandler(c *gin.Context) {
	// Check if it has been at least 60 seconds since the last OTP was sent
	if time.Since(lastOTPSendTime) < 60*time.Second {
		c.HTML(http.StatusOK, "verify.html", gin.H{"Message": "Please wait before requesting a new OTP"})
		return
	}

	c.HTML(http.StatusOK, "verify.html", gin.H{"Message": "OTP sented"})

	otp := helpers.GenerateOTP()
	helpers.SendOTP(otp, user.Email)

	// Update the last OTP send time
	lastOTPSendTime = time.Now()
}

func VerifyPost(c *gin.Context) {
	var verifyData models.VerifyData

	if err := c.ShouldBindJSON(&verifyData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	otp := verifyData.OTP
	status := helpers.VerifyOTP(otp, user.Email, c)
	log.Println("verifypost", otp, status)

	if !status {
		// Handle the case when OTP verification fails
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OTP"})
		return
	}

	// Attempt to create a new user
	newUser := models.User{
		Username: user.Username,
		Email:    user.Email,
		Number:   user.Number,
		Password: user.Password,
	}

	err := db.DB.Create(&newUser).Error

	if err != nil {
		// Check for duplicate key violation or other errors
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error creating user"})
		return
	}

	// Redirect to /login with a success message
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully. Please log in."})
}

// func HomeHandler(c *gin.Context) {
// 	var products []models.Products

// 	if err := db.DB.Preload("Images").Find(&products).Error; err != nil {

// 		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Product not found"})
// 		return
// 	}

// 	c.HTML(http.StatusOK, "home.html", gin.H{
// 		// "Productvariants": ProductVariants,
// 		"Products": products,
// 	})
// }

func HomeHandler(c *gin.Context) {
	var products []models.ProductVariants

	if err := db.DB.Preload("Product").Preload("Product.Images").Find(&products).Error; err != nil {
		fmt.Println("Error fetching product variant with product and images:", err)
		return
	}

	c.HTML(http.StatusOK, "home.html", gin.H{
		// "Productvariants": ProductVariants,
		"ProductVariants": products,
	})
}

func LogoutHandler(c *gin.Context) {

	c.SetCookie("auth", "", -1, "/", "localhost", false, true)

	c.Redirect(http.StatusSeeOther, "/login")
}

func ProductViewhandler(c *gin.Context) {
	slug := c.Query("Variant")
	var product models.ProductVariants

	if err := db.DB.Preload("Product").Preload("Product.Images").Where("Slug=?", slug).Find(&product).Error; err != nil {
		fmt.Println("Error fetching product variant with product and images:", err)
		return
	}

	result := make([]int, 0)
	for i := 1; i <= product.Stock && i <= 5; i++ {
		result = append(result, i)
	}
	fmt.Println("result", result)

	// Iterate through product variants to find unique Ram values

	c.HTML(http.StatusOK, "productdetail.html", gin.H{
		// "Productvariants": ProductVariants,
		"Product":  product,
		"Quantity": result,
	})

}
