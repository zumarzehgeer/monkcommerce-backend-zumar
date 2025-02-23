package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"zumarzehgeer.com/go/server/models"
)

// NOTE: Main Handler
func (h handler) ApplyCoupon(w http.ResponseWriter, r *http.Request) {
	// NOTE: Getting if from path url
	id, err := strconv.Atoi(r.URL.Path[len("/apply-coupon/"):])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// NOTE: Handeling methods
	switch r.Method {
	case "POST":
		applyCouponById(h, w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

// NOTE: Applying coupon by ID
func applyCouponById(h handler, w http.ResponseWriter, r *http.Request, id int) {

	// NOTE: Creating cart variable and getting data from the r.Body
	var cart models.Cart
	err := json.NewDecoder(r.Body).Decode(&cart)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// NOTE: Getting the coupon form its ID and setting encoding to false. Otherwise this func will also send a HTTP response.
	data, error := getCouponById(h, w, r, id, false)
	if error != nil {
		http.Error(w, "Coupon not found ", http.StatusNotFound)
		return
	}
	coupon := data.Data.(models.Coupon)

	// NOTE: Variables for sending the updated cart
	var totalPrice int
	var totalDiscount int
	for _, item := range cart.Items {
		totalPrice += item.Price * item.Quantity
		totalDiscount += item.TotalDiscount
	}

	// NOTE: Doing operations specific to the cart type
	if coupon.Type == models.CartWise {
		if isCartWiseApplicable(cart, coupon) {
			discount := calculateCartWiseDiscount(cart, coupon)
			totalDiscount = discount
		}
	} else if coupon.Type == models.ProductWise {
		if isProductWiseApplicable(cart, coupon) {
			for i := range cart.Items {
				discount := calculateProductWiseDiscount(cart, coupon)
				cart.Items[i].TotalDiscount = discount
				totalDiscount += discount
			}
		}
	} else if coupon.Type == models.Bxgy {
		if isBxgyApplicable(cart, coupon) {
			discount := calculateBxgyDiscount(cart, coupon)
			cart.Items = append(cart.Items, coupon.Details.GetProducts...)
			totalDiscount = discount
		}
	}

	cart.TotalPrice = totalPrice
	cart.TotalDiscount = totalDiscount
	cart.FinalPrice = totalPrice - totalDiscount

	// TODO: once the coupon is applied update the coupons repetition limit
	// NOTE: Writing Header and sending the HTTP response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cart)
}
