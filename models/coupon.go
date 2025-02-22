package models

import "gorm.io/gorm"

//  atlas schema inspect --env gorm --url env://src -w

type Coupon struct {
	gorm.Model
	Type           string         `json:"type"`
	Details        *CouponDetails `json:"details" gorm:"foreignKey:CouponID;references:ID"`
	Repition_limit int            `json:"repition_limit"`
}

type CouponDetails struct {
	gorm.Model
	CouponID  uint `json:"coupon_id"`
	Threshold int  `json:"threshold"`
	Discount  int  `json:"discount"`
	ProductID uint `json:"product_id"`

	BuyProducts []Product `json:"buy_products" gorm:"many2many:coupon_details_buy_products"`
	GetProducts []Product `json:"get_products" gorm:"many2many:coupon_details_get_products"`
}

type Product struct {
	gorm.Model

	Quantity      int `json:"quantity"`
	Price         int `json:"price"`
	TotalDiscount int `json:"total_discount"`

	BuyCouponDetails []CouponDetails `gorm:"many2many:coupon_details_buy_products;"`
	GetCouponDetails []CouponDetails `gorm:"many2many:coupon_details_get_products;"`
}

type Cart struct {
	gorm.Model

	Items         []Product `json:"items" gorm:"many2many:cart_items"`
	TotalPrice    int       `json:"total_price"`
	TotalDiscount int       `json:"total_discount"`
	FinalPrice    int       `json:"final_price"`
}

type ApplicableCoupons struct {
	CouponID   int    `json:"coupon_id"`
	CouponType string `json:"coupon_type"`
	Discount   int    `json:"discount"`
}

const (
	CartWise    = "cart-wise"
	ProductWise = "product-wise"
	Bxgy        = "bxgy"
)

// NOTE: method is assosiated with Coupon
func (c *Coupon) ValidateCouponType() bool {
	return c.Type == CartWise ||
		c.Type == ProductWise ||
		c.Type == Bxgy
}
