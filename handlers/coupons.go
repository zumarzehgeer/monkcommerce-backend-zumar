package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"zumarzehgeer.com/go/server/models"
)

// NOTE: Initializing the handler so that it can be accessable across the package
var h = handler{}

// NOTE: Main Handler
func (h handler) Coupons(w http.ResponseWriter, r *http.Request) {
	// NOTE: Handles methods
	switch r.Method {
	case "POST":
		createCoupon(h, w, r)
	case "GET":
		getAllCoupons(h, w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// NOTE: Creating a coupon
func createCoupon(h handler, w http.ResponseWriter, r *http.Request) {
	var data models.Coupon
	err := json.NewDecoder(r.Body).Decode(&data)

	// NOTE: Validate coupon type
	if !data.ValidateCouponType() {
		http.Error(w, "Invalid coupon type", http.StatusBadRequest)
		return
	}

	// NOTE: Coupon details required
	if data.Details == nil {
		http.Error(w, "Please provide coupon details", http.StatusBadRequest)
		return
	}

	// NOTE: Checking payload
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		fmt.Printf("error: %+v\n", err)
		return
	}

	// NOTE: Create coupon in the DB
	if result := h.DB.Create(&data); result.Error != nil {
		http.Error(w, "Could not create the coupon", http.StatusInternalServerError)
		return
	}

	// NOTE: Writing HTTP headers
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// NOTE: Creating response Object
	responseData := models.ResponseData{
		Status:  http.StatusCreated,
		Message: "Successfully processed the request",
		Data:    data,
	}

	// NOTE: Sending HTTP json response
	json.NewEncoder(w).Encode(responseData)
}

// NOTE: Getting all coupons
func getAllCoupons(h handler, w http.ResponseWriter, _ *http.Request, encode ...bool) ([]models.Coupon, error) {

	// NOTE: checks encoding (this is necessary otherwise the func will send 2 responses if used in other functions)
	shouldEncode := true
	if len(encode) > 0 {
		shouldEncode = encode[0]
	}

	// NOTE: Create coupon variable. Preload the Details and put it in coupons
	var coupons []models.Coupon
	if result := h.DB.Preload("Details").Preload("Details.BuyProducts").Preload("Details.GetProducts").Find(&coupons); result.Error != nil {
		fmt.Println(result.Error)
		return nil, result.Error
	}

	// NOTE: Writing headers
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// NOTE: Creating response object
	responseData := models.ResponseData{
		Status:  http.StatusOK,
		Message: "Successfully processed the request",
		Data:    coupons,
	}

	// NOTE: If TRUE sends HTTP json response
	if shouldEncode {
		json.NewEncoder(w).Encode(responseData)
	}

	// NOTE: Else return the response object
	return coupons, nil
}
