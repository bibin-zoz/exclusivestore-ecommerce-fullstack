package handlers

import (
	db "ecommercestore/database"
	"ecommercestore/helpers"
	"ecommercestore/models"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ForgotPasswordHandler(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Expires", "0")

	c.HTML(http.StatusOK, "forgotpass.html", nil)

}
func ForgotPasswordPostHandler(c *gin.Context) {

	var errors models.Invalid
	user = UserDetails{
		Email:           c.Request.FormValue("email"),
		Password:        c.Request.FormValue("password"),
		ConfirmPassword: c.Request.FormValue("confirmpassword"),
	}

	if user.Email == "" {
		errors.EmailError = "Email should not be empty"
		c.HTML(http.StatusBadRequest, "forgotpass.html", errors)
		return
	}
	var emailCount int64
	if err := db.DB.Table("users").Where("email = ?", user.Email).Count(&emailCount).Error; err != nil {
		fmt.Println(err)
		c.HTML(http.StatusBadRequest, "forgotpass.html", nil)
		return
	}

	if emailCount == 0 {
		errors.EmailError = "Email not  exists"
		c.HTML(http.StatusBadRequest, "forgotpass.html", errors)
		return
	}
	if user.Password == "" {
		errors.PasswordError = "Password should not be empty"
		c.HTML(http.StatusBadRequest, "forgotpass.html", errors)
		return
	}
	if user.Password != user.ConfirmPassword {
		errors.PasswordError = "Passwords do not match"
		c.HTML(http.StatusBadRequest, "forgotpass.html", errors)
		return
	}

	c.Redirect(http.StatusFound, "/emailverify")
}
func EmailVerify(c *gin.Context) {
	//for avoiding req in 60secondss
	if time.Since(lastOTPSendTime) < 60*time.Second {
		c.HTML(http.StatusOK, "emailverify.html", gin.H{"Message": "Please wait before requesting a new OTP"})
		return
	}

	c.HTML(http.StatusOK, "emailverify.html", gin.H{"Message": "OTP sented"})

	otp := helpers.GenerateOTP()
	helpers.SendOTP(otp, user.Email)

	// Update the last OTP send time
	lastOTPSendTime = time.Now()
}
func EmailVerifyPost(c *gin.Context) {
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

	// Find the user with the given email
	existingUser := models.User{}
	err := db.DB.Where("email = ?", user.Email).First(&existingUser).Error

	if err != nil {
		// Handle errors, including record not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
			return
		} else {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error finding user"})
			return
		}
	}

	// Update the user's password
	hashedPassword, _ := helpers.HashPassword(user.Password)
	existingUser.Password = hashedPassword

	err = db.DB.Save(&existingUser).Error
	if err != nil {
		// Handle update error
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error changing password"})
		return
	}

	// Redirect to /login with a success message
	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully. Please log in."})
}
