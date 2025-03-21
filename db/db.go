package db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"zumarzehgeer.com/go/server/models"
)

func Init() *gorm.DB {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	connStr := os.Getenv("DATABASE_URL")

	dsn := fmt.Sprintf("%s", connStr)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	err = db.AutoMigrate(&models.Cart{}, &models.Coupon{}, &models.CouponDetails{}, &models.Product{}, &models.ApplicableCoupons{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate: %v", err)
	}

	return db
}
