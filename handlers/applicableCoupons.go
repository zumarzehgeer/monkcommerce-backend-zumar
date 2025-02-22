package handlers

import (
	"encoding/json"
	"net/http"

	"zumarzehgeer.com/go/server/models"
)

func (h handler) ApplicableCoupons(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		// he will send me a cart
		// I will check all the applicable coupons for that cart
		// I will show him the total discount that will be applied by each coupon
		checkApplicableCoupons(h, w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func checkApplicableCoupons(h handler, w http.ResponseWriter, r *http.Request) {

	// get the cart from the request
	var cart models.Cart
	err := json.NewDecoder(r.Body).Decode(&cart)
	if err != nil {
		http.Error(w, "Invalid cart data", http.StatusBadRequest)
		return
	}

	// get all the coupons
	coupons, err := getAllCoupons(h, w, r, false)
	if err != nil {
		http.Error(w, "Could not get coupons", http.StatusInternalServerError)
		return
	}

	var applicableCoupons []models.ApplicableCoupons

	for _, coupon := range coupons {
		if coupon.Type == models.Bxgy {
			if isBxgyApplicable(cart, coupon) {
				var applicableCoupon models.ApplicableCoupons
				discount := calculateBxgyDiscount(cart, coupon)

				applicableCoupon.CouponID = int(coupon.ID)
				applicableCoupon.CouponType = coupon.Type
				applicableCoupon.Discount = discount // entire discount value

				applicableCoupons = append(applicableCoupons, applicableCoupon)
			}
		} else if coupon.Type == models.CartWise {
			if isCartWiseApplicable(cart, coupon) {
				var applicableCoupon models.ApplicableCoupons
				discount := calculateCartWiseDiscount(cart, coupon)

				applicableCoupon.CouponID = int(coupon.ID)
				applicableCoupon.CouponType = coupon.Type
				applicableCoupon.Discount = discount

				applicableCoupons = append(applicableCoupons, applicableCoupon)
			}
		} else if coupon.Type == models.ProductWise {
			if isProductWiseApplicable(cart, coupon) {
				var applicableCoupon models.ApplicableCoupons
				discount := calculateProductWiseDiscount(cart, coupon)

				applicableCoupon.CouponID = int(coupon.ID)
				applicableCoupon.CouponType = coupon.Type
				applicableCoupon.Discount = discount

				applicableCoupons = append(applicableCoupons, applicableCoupon)
			}
		}
	}

	// return the applicable coupons
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(applicableCoupons)
}

func isBxgyApplicable(cart models.Cart, coupon models.Coupon) bool {
	// Check if all buy products are in the cart
	buyProductsInCart := true
	for _, buyProduct := range coupon.Details.BuyProducts {
		found := false
		for _, cartItem := range cart.Items {
			if cartItem.ID == buyProduct.ID {
				found = true
				break
			}
		}
		if !found {
			buyProductsInCart = false
			break
		}
	}

	// If all buy products are in the cart, the coupon is applicable
	return buyProductsInCart
}

func isCartWiseApplicable(cart models.Cart, coupon models.Coupon) bool {
	// Check if the cart total is greater than the threshold
	cartTotal := 0
	for _, cartItem := range cart.Items {
		cartTotal += (cartItem.Price * cartItem.Quantity)
	}

	return cartTotal > coupon.Details.Threshold
}

func isProductWiseApplicable(cart models.Cart, coupon models.Coupon) bool {
	// Check if any of the cart items match the product ID
	for _, cartItem := range cart.Items {
		if cartItem.ID == coupon.Details.ProductID {
			return true
		}
	}

	return false
}

func calculateBxgyDiscount(_ models.Cart, coupon models.Coupon) (discount int) {
	discount = 0
	for _, getProduct := range coupon.Details.GetProducts {
		discount += (getProduct.Price * getProduct.Quantity)
	}

	return discount
}

func calculateCartWiseDiscount(cart models.Cart, coupon models.Coupon) (discount int) {
	// Calculate the cart total
	cartTotal := 0
	for _, cartItem := range cart.Items {
		cartTotal += (cartItem.Price * cartItem.Quantity)
	}

	// Calculate the discount
	discount = 0
	if cartTotal > coupon.Details.Threshold {
		discount = cartTotal * coupon.Details.Discount / 100
	}

	return discount
}

func calculateProductWiseDiscount(cart models.Cart, coupon models.Coupon) (discount int) {
	// Calculate the discount for product wise coupon
	discount = 0
	for _, cartItem := range cart.Items {
		if cartItem.ID == coupon.Details.ProductID {
			discount = (cartItem.Price * cartItem.Quantity) * coupon.Details.Discount / 100
		}
	}

	return discount
}
