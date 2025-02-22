package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"
	"zumarzehgeer.com/go/server/models"
)

func (h handler) CouponsId(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[len("/coupons/"):])

	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		getCouponById(h, w, r, id)
	case "PUT":
		updateCouponById(h, w, r, id)
	case "DELETE":
		deleteCouponById(h, w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getCouponById(h handler, w http.ResponseWriter, _ *http.Request, id int, encode ...bool) (models.ResponseData, error) {
	var coupon models.Coupon
	shouldEncode := true

	if len(encode) > 0 {
		shouldEncode = encode[0]
	}

	if result := h.DB.Preload("Details").Preload("Details.BuyProducts").Preload("Details.GetProducts").Preload("Details.BuyProducts.BuyCouponDetails").Preload("Details.GetProducts.GetCouponDetails").First(&coupon, id); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "record not found", http.StatusNotFound)
			return models.ResponseData{}, result.Error
		}
		http.Error(w, "Not Found", http.StatusNotFound)
		return models.ResponseData{}, result.Error
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Status", "200")

	responseData := models.ResponseData{
		Status:  http.StatusOK,
		Message: "Coupon found",
		Data:    coupon,
	}

	if shouldEncode {
		json.NewEncoder(w).Encode(responseData)
	}

	return responseData, nil
}

func updateCouponById(h handler, w http.ResponseWriter, r *http.Request, id int) (models.ResponseData, error) {
	data, error := getCouponById(h, w, r, id, false)
	if error != nil {
		return models.ResponseData{}, error
	}
	coupon := data.Data.(models.Coupon)

	var updatedCoupon models.Coupon
	err := json.NewDecoder(r.Body).Decode(&updatedCoupon)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return models.ResponseData{}, err
	}

	if !updatedCoupon.ValidateCouponType() {
		http.Error(w, "Invalid coupon type", http.StatusBadRequest)
		return models.ResponseData{}, err
	}
	coupon.Type = updatedCoupon.Type
	coupon.Repition_limit = updatedCoupon.Repition_limit
	coupon.UpdatedAt = time.Now()

	if result := h.DB.Save(&coupon); result.Error != nil {
		http.Error(w, "Failed to update coupon", http.StatusInternalServerError)
		return models.ResponseData{}, result.Error
	}

	if coupon.Type == models.CartWise {
		// threshold, discount
		coupon.Details.Threshold = updatedCoupon.Details.Threshold
		coupon.Details.Discount = updatedCoupon.Details.Discount

		// do the database operation
		if result := h.DB.Save(&coupon.Details); result.Error != nil {
			http.Error(w, "Failed to update coupon details", http.StatusInternalServerError)
			return models.ResponseData{}, result.Error
		}

	} else if coupon.Type == models.ProductWise {
		// product_id, discount
		coupon.Details.ProductID = updatedCoupon.Details.ProductID
		coupon.Details.Discount = updatedCoupon.Details.Discount

		// do the database operation
		if result := h.DB.Save(&coupon.Details); result.Error != nil {
			http.Error(w, "Failed to update coupon details", http.StatusInternalServerError)
			return models.ResponseData{}, result.Error
		}

	} else if coupon.Type == models.Bxgy {
		// buy_products, get_products
		err = h.DB.Model(&coupon.Details).Association("BuyProducts").Replace(updatedCoupon.Details.BuyProducts)
		if err != nil {
			http.Error(w, "Failed to update BuyProducts", http.StatusInternalServerError)
			return models.ResponseData{}, nil
		}

		err = h.DB.Model(&coupon.Details).Association("GetProducts").Replace(updatedCoupon.Details.GetProducts)
		if err != nil {
			http.Error(w, "Failed to update GetProducts", http.StatusInternalServerError)
			return models.ResponseData{}, nil
		}
	}

	// refetch coupon
	data, error = getCouponById(h, w, r, id, false)
	if error != nil {
		return models.ResponseData{}, error
	}
	coupon = data.Data.(models.Coupon)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Status", "200")

	responseData := models.ResponseData{
		Status:  http.StatusOK,
		Message: "Successfully updated the coupon",
		Data:    coupon,
	}

	json.NewEncoder(w).Encode(responseData)
	return responseData, nil
}

func deleteCouponById(h handler, w http.ResponseWriter, _ *http.Request, id int) (models.ResponseData, error) {
	var coupon models.Coupon
	if result := h.DB.First(&coupon, id); result.Error != nil {
		http.Error(w, "Record not found", http.StatusNotFound)
		return models.ResponseData{}, result.Error
	}

	if result := h.DB.Preload("Details").Preload("Details.BuyProducts").Preload("Details.GetProducts").Delete(&coupon); result.Error != nil {
		http.Error(w, "Failed to delete coupon", http.StatusInternalServerError)
		return models.ResponseData{}, result.Error
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Status", "200")

	responseData := models.ResponseData{
		Status:  http.StatusOK,
		Message: "Successfully deleted the coupon",
		Data:    coupon,
	}

	json.NewEncoder(w).Encode(responseData)
	return responseData, nil
}
