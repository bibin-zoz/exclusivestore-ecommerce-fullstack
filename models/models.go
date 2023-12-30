package models

import (
	"github.com/golang-jwt/jwt"
)

type Claims struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Status   string `json:"status"`
	jwt.StandardClaims
}
type Pagination struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
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
	IDError       string
	NameError     string
	EmailError    string
	NumberError   string
	PasswordError string
	RoleError     string
	CommonError   string
	LoginStatus   bool
	AmountError   string
	DateError     string
	StatusError   string
}

type Compare struct {
	ID       uint
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

//	type Category struct {
//		gorm.Model
//		CategoryName string `gorm:"unique;not null"`
//		Status       string `gorm:"default:'listed'"`
//		CreatedAt    time.Time
//	}

type OrderReq struct {
	CartID        string `json:"cartID"`
	AddressID     string `json:"addressID"`
	PaymentMethod string `json:"paymentMethod"`
	CouponCode    string `json:"couponcode"`
}
type Updatecart struct {
	ID       string `json:"id"`
	Quantity string `json:"quantity"`
}

type UserRequest struct {
	ID      int    `json:"id"`
	Request string `json:"request"`
}
