package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"zumarzehgeer.com/go/server/models"
)

func (h handler) Products(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		CreateProducts(h, w, r)
	case "GET":
		getAllProducts(h, w, r)
	case "DELETE":
		deleteAllProducts(h, w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func CreateProducts(h handler, w http.ResponseWriter, r *http.Request) {
	var data models.Product
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		fmt.Printf("error: %+v\n", err)
		return
	}

	if result := h.DB.Create(&data); result.Error != nil {
		http.Error(w, "Could not create the product", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	responseData := models.ResponseData{
		Status:  http.StatusCreated,
		Message: "Successfully created the product",
		Data:    data,
	}

	json.NewEncoder(w).Encode(responseData)
}

func getAllProducts(h handler, w http.ResponseWriter, _ *http.Request) {
	var products []models.Product

	if result := h.DB.Preload("BuyCouponDetails").Preload("GetCouponDetails").Find(&products); result.Error != nil {
		http.Error(w, "Could not fetch the products", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	responseData := models.ResponseData{
		Status:  http.StatusOK,
		Message: "Successfully fetched the products",
		Data:    products,
	}

	json.NewEncoder(w).Encode(responseData)
}

func deleteAllProducts(h handler, w http.ResponseWriter, _ *http.Request) {
	// Delete related records in coupon_details_get_products table
	if result := h.DB.Exec("DELETE FROM coupon_details_get_products"); result.Error != nil {
		http.Error(w, "Could not delete related records in coupon_details_get_products", http.StatusInternalServerError)
		return
	}

	// Delete related records in coupon_details_buy_products table
	if result := h.DB.Exec("DELETE FROM coupon_details_buy_products"); result.Error != nil {
		http.Error(w, "Could not delete related records in coupon_details_buy_products", http.StatusInternalServerError)
		return
	}

	// Now delete from products table
	if result := h.DB.Exec("DELETE FROM products"); result.Error != nil {
		http.Error(w, "Could not delete the products", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	responseData := models.ResponseData{
		Status:  http.StatusOK,
		Message: "Successfully deleted all products",
	}

	json.NewEncoder(w).Encode(responseData)
}
