package db

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"zumarzehgeer.com/go/server/models"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "zzehgeer"
	password = "admin"
	dbname   = "monkcommerce"
)

func Init() *gorm.DB {

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	err = db.AutoMigrate(&models.Cart{}, &models.Coupon{}, &models.CouponDetails{}, &models.Product{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate: %v", err)
	}

	return db
}
