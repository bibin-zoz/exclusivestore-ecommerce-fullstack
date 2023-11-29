package helpers

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var otpMap = make(map[string]string)

func GenerateOTP() string {
	source := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(source)
	return fmt.Sprintf("%06d", randGen.Intn(1000000))
}

func SendOTP(otp string, email string) error {
	from := os.Getenv("email")
	password := os.Getenv("password")
	to := email
	log.Println("email", email, otp)
	smtpServer := "smtp.gmail.com"
	smtpPort := "587"
	otpMap[email] = otp

	auth := smtp.PlainAuth("", from, password, smtpServer)

	message := fmt.Sprintf("Subject: Your OTP\n\nYour OTP is: %s", otp)

	err := smtp.SendMail(smtpServer+":"+smtpPort, auth, from, []string{to}, []byte(message))
	if err != nil {
		return err
	}

	return nil
}

func VerifyOTP(otp string, email string, c *gin.Context) bool {
	userEmail := email
	enteredOTP := otp

	storedOTP, ok := otpMap[userEmail]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OTP not found for the given Email"})

		return false
	}

	if enteredOTP == storedOTP {
		// Clear the OTP from the map after successful verification
		delete(otpMap, userEmail)
		// Render HTML page with a success message
		// c.HTML(http.StatusOK, "verify.html", gin.H{"message": "OTP verified successfully"})
		// // Send JSON response with the same success message
		// c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
		return true
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OTP"})
		return false
	}
}
