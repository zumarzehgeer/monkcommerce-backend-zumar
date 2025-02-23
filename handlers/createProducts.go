package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"zumarzehgeer.com/go/server/models"
)

// NOTE: Main function for products
func (h handler) Products(w http.ResponseWriter, r *http.Request) {
	// NOTE: Handling methods
	switch r.Method {
	case "POST":
		CreateProduct(h, w, r)
	case "GET":
		getAllProducts(h, w, r)
	case "DELETE":
		deleteAllProducts(h, w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// NOTE: Create Product
func CreateProduct(h handler, w http.ResponseWriter, r *http.Request) {
	// NOTE: Creating the data of type product and decoding the r.Body
	var data models.Product
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		fmt.Printf("error: %+v\n", err)
		return
	}

	// NOTE: Creating the Product
	if result := h.DB.Create(&data); result.Error != nil {
		http.Error(w, "Could not create the product", http.StatusInternalServerError)
		return
	}

	// NOTE: Writing the headers
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// NOTE: Creating response object
	responseData := models.ResponseData{
		Status:  http.StatusCreated,
		Message: "Successfully created the product",
		Data:    data,
	}

	// NOTE: Sending HTTP response
	json.NewEncoder(w).Encode(responseData)
}

// NOTE: Get All Products
func getAllProducts(h handler, w http.ResponseWriter, _ *http.Request) {
	// NOTE: Getting products and preloadingthe buy and get coupons details
	var products []models.Product
	if result := h.DB.Preload("BuyCouponDetails").Preload("GetCouponDetails").Find(&products); result.Error != nil {
		http.Error(w, "Could not fetch the products", http.StatusInternalServerError)
		return
	}

	// NOTE: Writing the HTTP headers
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// NOTE: Creating the response object
	responseData := models.ResponseData{
		Status:  http.StatusOK,
		Message: "Successfully fetched the products",
		Data:    products,
	}

	// NOTE: Sending the HTTP response
	json.NewEncoder(w).Encode(responseData)
}

// NOTE: Delete all products
func deleteAllProducts(h handler, w http.ResponseWriter, _ *http.Request) {
	// NOTE: Delete related records in coupon_details_get_products table
	if result := h.DB.Exec("DELETE FROM coupon_details_get_products"); result.Error != nil {
		http.Error(w, "Could not delete related records in coupon_details_get_products", http.StatusInternalServerError)
		return
	}

	// NOTE: Delete related records in coupon_details_buy_products table
	if result := h.DB.Exec("DELETE FROM coupon_details_buy_products"); result.Error != nil {
		http.Error(w, "Could not delete related records in coupon_details_buy_products", http.StatusInternalServerError)
		return
	}

	// NOTE: Delete from products table
	if result := h.DB.Exec("DELETE FROM products"); result.Error != nil {
		http.Error(w, "Could not delete the products", http.StatusInternalServerError)
		return
	}

	// NOTE: Writing HTTP headers
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// NOTE: Creating response object
	responseData := models.ResponseData{
		Status:  http.StatusOK,
		Message: "Successfully deleted all products",
	}

	// NOTE: Sending the HTTP response
	json.NewEncoder(w).Encode(responseData)
}
