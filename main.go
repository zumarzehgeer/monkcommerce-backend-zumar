package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"zumarzehgeer.com/go/server/db"
	"zumarzehgeer.com/go/server/handlers"

	_ "ariga.io/atlas-provider-gorm/gormschema"
)

func main() {

	// NOTE: Initiate DB
	DB := db.Init()
	h := handlers.New(DB)
	fmt.Println("Successfully connected!")

	// NOTE: Initiate Routes
	http.HandleFunc("/coupons", h.Coupons)
	http.HandleFunc("/coupons/", h.CouponsId)
	http.HandleFunc("/applicable-coupons", h.ApplicableCoupons)
	http.HandleFunc("/apply-coupon/", h.ApplyCoupon)
	http.HandleFunc("/products", h.Products)

	// NOTE: listening to server
	fmt.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
