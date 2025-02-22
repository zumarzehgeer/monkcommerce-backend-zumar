package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"zumarzehgeer.com/go/server/models"
)

var h = handler{}

func (h handler) Coupons(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		createCoupon(h, w, r)
	case "GET":
		getAllCoupons(h, w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func createCoupon(h handler, w http.ResponseWriter, r *http.Request) {
	/*
		TODO:
			1) create coupon
						* "cart-wise"
						* "product-wise"
						* "bxgy"
	*/
	var data models.Coupon
	err := json.NewDecoder(r.Body).Decode(&data)

	if !data.ValidateCouponType() {
		http.Error(w, "Invalid coupon type", http.StatusBadRequest)
		return
	}

	if data.Details == nil {
		http.Error(w, "Please provide coupon details", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		fmt.Printf("error: %+v\n", err)
		return
	}

	if result := h.DB.Create(&data); result.Error != nil {
		http.Error(w, "Could not create the coupon", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	responseData := models.ResponseData{
		Status:  http.StatusCreated,
		Message: "Successfully processed the request",
		Data:    data,
	}

	json.NewEncoder(w).Encode(responseData)

}

func getAllCoupons(h handler, w http.ResponseWriter, _ *http.Request, encode ...bool) ([]models.Coupon, error) {

	shouldEncode := true

	if len(encode) > 0 {
		shouldEncode = encode[0]
	}

	var coupons []models.Coupon

	if result := h.DB.Preload("Details").Preload("Details.BuyProducts").Preload("Details.GetProducts").Find(&coupons); result.Error != nil {
		fmt.Println(result.Error)
		return nil, result.Error
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	responseData := models.ResponseData{
		Status:  http.StatusOK,
		Message: "Successfully processed the request",
		Data:    coupons,
	}

	if shouldEncode {
		json.NewEncoder(w).Encode(responseData)
	}
	return coupons, nil
}
