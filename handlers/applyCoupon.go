package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"zumarzehgeer.com/go/server/models"
)

func (h handler) ApplyCoupon(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[len("/apply-coupon/"):])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	fmt.Println("apply coupon id:", id)

	switch r.Method {
	case "POST":
		applyCouponById(h, w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func applyCouponById(h handler, w http.ResponseWriter, r *http.Request, id int) {

	var cart models.Cart
	err := json.NewDecoder(r.Body).Decode(&cart)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	data, error := getCouponById(h, w, r, id, false)
	if error != nil {
		http.Error(w, "Coupon not found ", http.StatusNotFound)
		return
	}

	coupon := data.Data.(models.Coupon)

	var totalPrice int
	var totalDiscount int
	for _, item := range cart.Items {
		totalPrice += item.Price * item.Quantity
		totalDiscount += item.TotalDiscount
	}

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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}
