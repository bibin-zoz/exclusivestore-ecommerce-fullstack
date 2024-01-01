package handlers

import (
	db "ecommercestore/database"
	"ecommercestore/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"ecommercestore/helpers"

	"github.com/gin-gonic/gin"
)

type DeleteRequest struct {
	ID int `json:"id"`
}
type UserDetails struct {
	Username        string
	Email           string
	Number          string
	Password        string
	ConfirmPassword string
	ReferalCode     string
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

	err := helpers.VerifyPassword(Newpassword, compare.Password)
	if err != nil {
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
		fmt.Println("user blocked")
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

	c.SetCookie("auth", string(userDetailsJSON), 0, "/", "localhost", true, true)

	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Expires", "0")

	c.Redirect(http.StatusFound, "/home")
}

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
		ReferalCode:     c.Request.FormValue("referralCode"),
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
	var numberCount int64
	if err := db.DB.Table("users").Where("number = ?", user.Number).Count(&numberCount).Error; err != nil {
		fmt.Println(err)
		c.HTML(http.StatusBadRequest, "signup.html", nil)
		return
	}

	if numberCount > 0 {
		errors.NumberError = "Number already exists"
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
	if len(user.Password) < 6 {
		errors.PasswordError = "Password length should be 6 "
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
	//for avoiding req in 60secondss
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

	hasedPassword, _ := helpers.HashPassword(user.Password)
	newUser := models.User{
		Username: user.Username,
		Email:    user.Email,
		Number:   user.Number,
		Password: hasedPassword,
	}

	err := db.DB.Create(&newUser).Error

	if err != nil {
		// Check for duplicate key violation or other errors
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error creating user"})
		return
	}
	helpers.UpdateReferalCount(user.ReferalCode)
	referalDetails := models.ReferalDetails{
		UserID:      1, // Replace with the actual user ID
		Count:       0,
		ReferalCode: helpers.GenerateRandomReferalCode(),
	}

	// Save to the database
	err = db.DB.Create(&referalDetails).Error
	if err != nil {
		// Check for duplicate key violation or other errors
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error creating Referal id"})
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
	if err := db.DB.Preload("Product").Preload("Product.Category").Preload("Product.Category.CategoryOffer").Preload("Product.Images").Where("status='listed'").Find(&products).Error; err != nil {
		fmt.Println("Error fetching product variant with product and images:", err)
		return
	}

	c.HTML(http.StatusOK, "home.html", gin.H{
		// "Productvariants": ProductVariants,
		"ProductVariants": products,
	})
}

func LogoutHandler(c *gin.Context) {
	fmt.Println("user logut")

	c.SetCookie("auth", "", -1, "/", "localhost", false, true)

	c.Redirect(http.StatusSeeOther, "/login")
}

func ProductViewhandler(c *gin.Context) {
	ID, _ := helpers.GetID(c)
	slug := c.Query("Variant")
	var product models.ProductVariants

	if err := db.DB.Preload("Product.Category").Preload("Product.Category.CategoryOffer").Preload("Product").Preload("Product.Images").Where("Slug=?", slug).Find(&product).Error; err != nil {
		fmt.Println("Error fetching product variant with product and images:", err)
		return
	}

	result := make([]int, 0)
	for i := 1; i <= product.Stock && i <= 5; i++ {
		result = append(result, i)
	}

	var Cartid uint
	db.DB.Table("carts").Where("user_id=?", ID).Select("id").Scan(&Cartid)

	var count int64
	var cartcount models.CartProducts
	query := db.DB.Where("cart_id=? AND variant_id=?", Cartid, product.ID).Find(&cartcount)
	fmt.Println("SQL Query:", query.Statement.SQL.String())

	query.Count(&count)
	fmt.Println("Count:", count)

	if err := query.Error; err != nil {
		fmt.Println("Error fetching cart count:", err)
		return
	}

	c.HTML(http.StatusOK, "productdetail.html", gin.H{
		"Product":  product,
		"Quantity": result,
		"Count":    count,
	})
}

func UserDashboardHandler(c *gin.Context) {
	var products []models.ProductVariants

	if err := db.DB.Preload("Product").Preload("Product.Images").Find(&products).Error; err != nil {
		fmt.Println("Error fetching product variant with product and images:", err)
		return
	}

	c.HTML(http.StatusOK, "userdashboard.html", gin.H{
		// "Productvariants": ProductVariants,
		"ProductVariants": products,
	})

}

func UserAddressHandler(c *gin.Context) {
	var userAddress []models.UserAddress

	usercookie, _ := c.Cookie("auth")
	var token models.TokenUser
	err := json.NewDecoder(strings.NewReader(usercookie)).Decode(&token)
	if err != nil {
		fmt.Println("Error fetching  UserDetails :", err)

		return
	}
	Claims, err := helpers.ParseToken(token.AccessToken)
	if err != nil {
		fmt.Println("Error fetching  UserDetails :", err)

		return

	}
	// fmt.Println("claims", Claims)
	db.DB.Where("user_id=?", Claims.ID).Find(&userAddress)
	// fmt.Println("useradd", userAddress)
	c.JSON(http.StatusOK, userAddress)

}
func DeleteAddressHandler(c *gin.Context) {
	var req DeleteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ID := req.ID

	var address models.UserAddress
	result := db.DB.Where("id = ?", ID).First(&address)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch address"})
		return
	}

	if address.IsPrimary == "true" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot remove primary address"})
		return
	}

	result = db.DB.Delete(&address)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "address removed successfully"})
}

func NewAddressHandler(c *gin.Context) {
	var address models.UserAddress

	ID, _ := helpers.GetID(c)
	if err := c.ShouldBind(&address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fill All Fields"})
		return
	}

	if address.State == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "State is required"})
		return
	}

	// Validate City
	if address.City == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "City is required"})
		return
	}
	if address.PhoneNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PhoneNumber is required"})
		return
	}

	// Validate City
	if address.PostalCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PostalCode is required"})
		return
	}

	fmt.Println("address", address)
	address.UserID = ID

	var count int64 = 0
	db.DB.Where("user_id=? AND is_primary='true'", ID).Count(&count)
	if count == 0 {
		address.IsPrimary = "true"
	}

	result := db.DB.Create(&address)

	if result.Error != nil {
		fmt.Println("hiii")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save address"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Address saved successfully"})

}

// profile
func GetUserProfileHandler(c *gin.Context) {
	var userdetails models.User
	usercookie, _ := c.Cookie("auth")
	var token models.TokenUser
	err := json.NewDecoder(strings.NewReader(usercookie)).Decode(&token)
	if err != nil {
		fmt.Println("Error fetching  UserDetails :", err)

		return
	}
	Claims, err := helpers.ParseToken(token.AccessToken)
	if err != nil {
		fmt.Println("Error fetching  UserDetails :", err)

		return

	}
	result := db.DB.Preload("ReferalDetails").Where("id=?", Claims.ID).Find(&userdetails)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove address"})
		return

	}
	c.JSON(http.StatusOK, userdetails)

}

func UpdateUserProfileHandler(c *gin.Context) {
	var userdetails models.UserDetail
	ID, err := helpers.GetID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	if err := c.ShouldBind(&userdetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := db.DB.Exec("UPDATE users SET username = ?, email = ?, number = ? WHERE id = ?", userdetails.UserName, userdetails.Email, userdetails.PhoneNumber, ID)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user details"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated user details successfully"})
}

func UpdatePasswordHandler(c *gin.Context) {
	var comparePassword models.UpdatePassword

	ID, err := helpers.GetID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBind(&comparePassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if comparePassword.NewPassword != comparePassword.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords don't match"})
		return
	}

	if comparePassword.NewPassword == comparePassword.Password {
		c.JSON(http.StatusBadRequest, gin.H{"error": "New password is the same as the current password"})
		return
	}

	var currentPassword string
	result := db.DB.Raw("SELECT password FROM users WHERE id=?", ID).Scan(&currentPassword)

	if result.Error != nil {
		fmt.Println("Error fetching current password:", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch current password"})
		return
	}

	if currentPassword != comparePassword.Password {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Current password is incorrect"})
		return
	}

	result = db.DB.Exec("UPDATE users SET password = ? WHERE id = ?", comparePassword.NewPassword, ID)

	if result.Error != nil {
		fmt.Println("Error updating password:", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password updated successfully"})
}
