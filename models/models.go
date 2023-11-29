package models

import (
	"github.com/golang-jwt/jwt"
)

type Claims struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Status   string `json:"status"`
	jwt.StandardClaims
}

// type AuthUserClaims struct {
// 	Email string `json:"email"`
// 	Role  string `json:"role"`
// 	jwt.StandardClaims
// }

type VerifyData struct {
	OTP string `json:"otp"`
}

type Invalid struct {
	NameError     string
	EmailError    string
	NumberError   string
	PasswordError string
	RoleError     string
	CommonError   string
	LoginStatus   bool
	StatusError   string
}

type Compare struct {
	ID       int
	Password string
	Role     string
	Username string
	Email    string
	Status   string
}

// type Product struct {
// 	Product_ID   primitive.ObjectID `bson:"_id"`
// 	Product_Name *string            `json:"product_name"`
// 	Seller_ID
// 	Category_ID

// 	Price *uint64 `json:"price"`
// 	Discount_Price
// 	Rating *uint8  `json:"rating"`
// 	Image  *string `json:"image"`
// }

// type Category struct {
// 	gorm.Model
// 	CategoryName string `gorm:"unique;not null"`
// 	Status       string `gorm:"default:'listed'"`
// 	CreatedAt    time.Time
// }
