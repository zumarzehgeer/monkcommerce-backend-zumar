package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"
	"zumarzehgeer.com/go/server/models"
)

// NOTE: Main Handler
func (h handler) CouponsId(w http.ResponseWriter, r *http.Request) {
	// NOTE: Gets ID from path url
	id, err := strconv.Atoi(r.URL.Path[len("/coupons/"):])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// NOTE: Handles function on the routes
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

// NOTE: Gets a Coupon by its ID
func getCouponById(h handler, w http.ResponseWriter, _ *http.Request, id int, encode ...bool) (models.ResponseData, error) {

	// NOTE: Creating coupon variable and checking its encoding
	var coupon models.Coupon
	shouldEncode := true
	if len(encode) > 0 {
		shouldEncode = encode[0]
	}

	// NOTE: Preload the details and find the coupon by ID in the DB
	if result := h.DB.Preload("Details").Preload("Details.BuyProducts").Preload("Details.GetProducts").Preload("Details.BuyProducts.BuyCouponDetails").Preload("Details.GetProducts.GetCouponDetails").First(&coupon, id); result.Error != nil {
		// NOTE: Checks if the record exists
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "record not found", http.StatusNotFound)
			return models.ResponseData{}, result.Error
		}
		http.Error(w, "Not Found", http.StatusNotFound)
		return models.ResponseData{}, result.Error
	}

	// NOTE: Writing HTTP headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Status", "200")

	// NOTE: Creating response object
	responseData := models.ResponseData{
		Status:  http.StatusOK,
		Message: "Coupon found",
		Data:    coupon,
	}

	// NOTE: Checks encoding If true sends https response
	if shouldEncode {
		json.NewEncoder(w).Encode(responseData)
	}

	// NOTE: Otherwise send the responseData object
	return responseData, nil
}

// NOTE: Update a coupon by its ID
func updateCouponById(h handler, w http.ResponseWriter, r *http.Request, id int) (models.ResponseData, error) {
	// NOTE: Getting a specific coupon by its ID
	data, error := getCouponById(h, w, r, id, false)
	if error != nil {
		return models.ResponseData{}, error
	}
	coupon := data.Data.(models.Coupon)

	// NOTE: Create update coupon variable and decoding the r.Body
	var updatedCoupon models.Coupon
	err := json.NewDecoder(r.Body).Decode(&updatedCoupon)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return models.ResponseData{}, err
	}

	// NOTE: Validating the coupon type
	if !updatedCoupon.ValidateCouponType() {
		http.Error(w, "Invalid coupon type", http.StatusBadRequest)
		return models.ResponseData{}, err
	}

	// NOTE: Update teh coupon object with new values
	coupon.Type = updatedCoupon.Type
	coupon.Repition_limit = updatedCoupon.Repition_limit
	coupon.UpdatedAt = time.Now()

	// NOTE: Saving it into DB and checking if the operation in DB was successful
	if result := h.DB.Save(&coupon); result.Error != nil {
		http.Error(w, "Failed to update coupon", http.StatusInternalServerError)
		return models.ResponseData{}, result.Error
	}

	// NOTE: Based on Coupon type that the user sent, updating the coupon
	if coupon.Type == models.CartWise {
		// NOTE: updates the Threshold & Discount Percentage
		coupon.Details.Threshold = updatedCoupon.Details.Threshold
		coupon.Details.Discount = updatedCoupon.Details.Discount

		// NOTE: Saving the updated values in the DB
		if result := h.DB.Save(&coupon.Details); result.Error != nil {
			http.Error(w, "Failed to update coupon details", http.StatusInternalServerError)
			return models.ResponseData{}, result.Error
		}

	} else if coupon.Type == models.ProductWise {
		// NOTE: updates the Product ID and Discount Percentage for a specific product
		coupon.Details.ProductID = updatedCoupon.Details.ProductID
		coupon.Details.Discount = updatedCoupon.Details.Discount

		// NOTE: Updating the values in the DB
		if result := h.DB.Save(&coupon.Details); result.Error != nil {
			http.Error(w, "Failed to update coupon details", http.StatusInternalServerError)
			return models.ResponseData{}, result.Error
		}

	} else if coupon.Type == models.Bxgy {
		// NOTE: Updates the buy_products & get_products

		// NOTE: Replaces the current buy_product with the new products that the use sent
		err = h.DB.Model(&coupon.Details).Association("BuyProducts").Replace(updatedCoupon.Details.BuyProducts)
		if err != nil {
			http.Error(w, "Failed to update BuyProducts", http.StatusInternalServerError)
			return models.ResponseData{}, nil
		}

		// NOTE: Replaces the current get_product with the new products that the use sent
		err = h.DB.Model(&coupon.Details).Association("GetProd0ucts").Replace(updatedCoupon.Details.GetProducts)
		if err != nil {
			http.Error(w, "Failed to update GetProducts", http.StatusInternalServerError)
			return models.ResponseData{}, nil
		}
	}

	// NOTE: Refetching the updated coupon
	data, error = getCouponById(h, w, r, id, false)
	if error != nil {
		return models.ResponseData{}, error
	}
	coupon = data.Data.(models.Coupon)

	// NOTE: Writing the HTTP headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Status", "200")

	// NOTE: Creating response object
	responseData := models.ResponseData{
		Status:  http.StatusOK,
		Message: "Successfully updated the coupon",
		Data:    coupon,
	}

	// NOTE: Sending the HTTP response as well as returning the responseData, error
	json.NewEncoder(w).Encode(responseData)
	return responseData, nil
}

// NOTE: Delete a coupon by its ID
func deleteCouponById(h handler, w http.ResponseWriter, _ *http.Request, id int) (models.ResponseData, error) {

	// NOTE: Finds the coupon
	var coupon models.Coupon
	if result := h.DB.First(&coupon, id); result.Error != nil {
		http.Error(w, "Record not found", http.StatusNotFound)
		return models.ResponseData{}, result.Error
	}

	// NOTE: Delete the coupon from DB
	if result := h.DB.Preload("Details").Preload("Details.BuyProducts").Preload("Details.GetProducts").Delete(&coupon); result.Error != nil {
		http.Error(w, "Failed to delete coupon", http.StatusInternalServerError)
		return models.ResponseData{}, result.Error
	}

	// NOTE: Writing the HTTP headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Status", "200")

	// NOTE: Creating response object
	responseData := models.ResponseData{
		Status:  http.StatusOK,
		Message: "Successfully deleted the coupon",
		Data:    coupon,
	}

	// NOTE: Send the HTTP json response and responseData, error
	json.NewEncoder(w).Encode(responseData)
	return responseData, nil
}
