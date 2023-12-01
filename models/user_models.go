package models

import (
	"time"

	"gorm.io/gorm"
)

var (
	DB *gorm.Model
)

type User struct {
	gorm.Model
	Username  string `gorm:"unique;not null"`
	Email     string `gorm:"unique;not null"`
	Number    string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	Role      string `gorm:"default:'user'"`
	Status    string `gorm:"default:'active'"`
	CreatedAt time.Time
}
type TokenUser struct {
	Users        Claims
	AccessToken  string
	RefreshToken string
}

type UserAddress struct {
	gorm.Model
	UserID      uint   `json:"userID"  gorm:"index;foreignKey:UserID"`
	Street      string `gorm:"not null"  json:"street" form:"street" binding:"required"`
	City        string `gorm:"not null" json:"city" form:"city" binding:"required"`
	State       string `gorm:"not null" json:"state" form:"state" binding:"required"`
	PostalCode  string `gorm:"not null" json:"postalcode" form:"postalcode" binding:"required"`
	Country     string `gorm:"not null" json:"country" form:"country" binding:"required"`
	PhoneNumber string `json:"phone_number" gorm:"not null" form:"phonenumber" binding:"required"`
}

type UserDetail struct {
	UserName    string `form:"username" binding:"required"`
	Email       string `form:"email" binding:"required"`
	PhoneNumber string `form:"mobile" binding:"required"`
}
type UpdatePassword struct {
	Password        string `form:"password" binding:"required"`
	NewPassword     string `form:"newpassword" binding:"required"`
	ConfirmPassword string `form:"confirmpassword" binding:"required"`
}
