package handlers

import (
	"encoding/json"
	"net/http"

	"zumarzehgeer.com/go/server/models"
)

// NOTE: Check applicable coupons for a specifix cart
func (h handler) ApplicableCoupons(w http.ResponseWriter, r *http.Request) {
	// NOTE: Handle methods
	switch r.Method {
	case "POST":
		checkApplicableCoupons(h, w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// NOTE: Checks applicable coupons
func checkApplicableCoupons(h handler, w http.ResponseWriter, r *http.Request) {
	// NOTE: Gets the cart from r.Body
	var cart models.Cart
	err := json.NewDecoder(r.Body).Decode(&cart)
	if err != nil {
		http.Error(w, "Invalid cart data", http.StatusBadRequest)
		return
	}

	// NOTE: Gets all the coupons
	coupons, err := getAllCoupons(h, w, r, false)
	if err != nil {
		http.Error(w, "Could not get coupons", http.StatusInternalServerError)
		return
	}

	// NOTE: Create an empty var for applicable coupons
	var applicableCoupons []models.ApplicableCoupons

	// NOTE: Do the logic here
	for _, coupon := range coupons {
		// NOTE: Loop over the coupons & check for cart types
		if coupon.Type == models.CartWise {
			// NOTE: Checks if cart wise is applicable for the cart
			if isCartWiseApplicable(cart, coupon) {
				// NOTE: Then update the information
				var applicableCoupon models.ApplicableCoupons
				// NOTE: Calculate the discount
				discount := calculateCartWiseDiscount(cart, coupon)

				applicableCoupon.CouponID = int(coupon.ID)
				applicableCoupon.CouponType = coupon.Type
				applicableCoupon.Discount = discount

				// NOTE: Then append the coupon to the Applicable Coupons Variable
				applicableCoupons = append(applicableCoupons, applicableCoupon)
			}
		} else if coupon.Type == models.ProductWise {
			// NOTE: Check applicability
			if isProductWiseApplicable(cart, coupon) {
				// NOTE: Update the information
				var applicableCoupon models.ApplicableCoupons
				// NOTE: Calculate the discount
				discount := calculateProductWiseDiscount(cart, coupon)

				applicableCoupon.CouponID = int(coupon.ID)
				applicableCoupon.CouponType = coupon.Type
				applicableCoupon.Discount = discount

				// NOTE: Append the coupon into the cart
				applicableCoupons = append(applicableCoupons, applicableCoupon)
			}
		} else if coupon.Type == models.Bxgy {
			// NOTE: Check the applicablity
			if isBxgyApplicable(cart, coupon) {
				var applicableCoupon models.ApplicableCoupons
				// NOTE: Calculate the discount
				discount := calculateBxgyDiscount(cart, coupon)

				applicableCoupon.CouponID = int(coupon.ID)
				applicableCoupon.CouponType = coupon.Type
				applicableCoupon.Discount = discount

				// NOTE: Append the coupon to applicable coupons
				applicableCoupons = append(applicableCoupons, applicableCoupon)
			}
		}

	}

	// NOTE: Write the HTTP header and return the applicable coupons response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(applicableCoupons)
}

func isBxgyApplicable(cart models.Cart, coupon models.Coupon) bool {
	// NOTE: Check if all buy products are in the cart
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

	// NOTE: If all buy products are in the cart, the coupon is applicable
	return buyProductsInCart
}

func isCartWiseApplicable(cart models.Cart, coupon models.Coupon) bool {
	// NOTE: Check if the cart total is greater than the threshold
	cartTotal := 0
	for _, cartItem := range cart.Items {
		cartTotal += (cartItem.Price * cartItem.Quantity)
	}

	// NOTE: If cart total is greater than the threshold It will return true
	return cartTotal > coupon.Details.Threshold
}

func isProductWiseApplicable(cart models.Cart, coupon models.Coupon) bool {
	// NOTE: Check if any of the cart items match the product ID
	for _, cartItem := range cart.Items {
		if cartItem.ID == coupon.Details.ProductID {
			return true
		}
	}
	return false
}

func calculateCartWiseDiscount(cart models.Cart, coupon models.Coupon) (discount int) {
	// NOTE: Calculate the cart total
	cartTotal := 0
	for _, cartItem := range cart.Items {
		cartTotal += (cartItem.Price * cartItem.Quantity)
	}

	// NOTE: Then calculate the discount
	discount = 0
	if cartTotal > coupon.Details.Threshold {
		discount = cartTotal * coupon.Details.Discount / 100
	}

	return discount
}

func calculateProductWiseDiscount(cart models.Cart, coupon models.Coupon) (discount int) {
	// NOTE: Calculate the discount for a product
	discount = 0
	for _, cartItem := range cart.Items {
		if cartItem.ID == coupon.Details.ProductID {
			discount = (cartItem.Price * cartItem.Quantity) * coupon.Details.Discount / 100
		}
	}

	return discount
}

func calculateBxgyDiscount(_ models.Cart, coupon models.Coupon) (discount int) {
	discount = 0
	for _, getProduct := range coupon.Details.GetProducts {
		discount += (getProduct.Price * getProduct.Quantity)
	}
	return discount
}
