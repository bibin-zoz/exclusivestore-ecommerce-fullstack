package database

import (
	"ecommercestore/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dsn := "user=" + dbUsername + " password=" + dbPassword + " dbname=" + dbName + " host=" + dbHost + " port=" + dbPort + " sslmode=disable"

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database")
	}
	if err := DB.AutoMigrate(&models.User{}); err != nil {

		panic(err)
	}
	if err := DB.AutoMigrate(&models.Categories{}); err != nil {

		panic(err)
	}
	if err := DB.AutoMigrate(&models.Products{}); err != nil {

		panic(err)
	}

	DB.AutoMigrate(&models.Image{}, models.Brands{}, models.ProductVariants{})
	// if err := DB.AutoMigrate(&models.ProductImage{}); err != nil {

	// 	panic(err)
	// }

}
