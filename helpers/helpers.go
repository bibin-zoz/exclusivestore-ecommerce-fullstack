package helpers

import (
	"ecommercestore/models"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/golang-jwt/jwt"
)

// var twilioSID = "ACde3aa36171cc3d73b605996e0d73ed6f"
// var twilioAuthToken = "188a1c697848e8854a397ce5bdb69c64"
// var twilioFromNumber = "+12563882106"
var secretKey []byte

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

func CreateToken(c *gin.Context, user models.Compare) {
	expirationTime := time.Now().Add(15 * time.Minute) // Adjust as needed
	claims := &models.Claims{
		Username: user.Username,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	fmt.Println("us", user.Username)
	fmt.Printf("Claims: %+v\n", claims)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtKey := []byte(os.Getenv("jwtKey"))
	fmt.Println("JWT Key:", jwtKey)
	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		fmt.Println("Error signing token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.SetCookie("token", signedToken, int(expirationTime.Unix()), "/", "localhost", false, true)

	c.Status(http.StatusOK)

}

func GetUserRoleFromToken(tokenString string) (string, error) {
	secretKey = []byte(os.Getenv("jwtKey"))
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		fmt.Println("", claims.Role)
		return claims.Role, nil
	}

	return "", fmt.Errorf("invalid token")

}
