package user

import "time"

// Structs for requests and responses
type LoginRequest struct {
	Phone string `json:"phone" binding:"required"`
}

type LoginResponse struct {
	Success  bool          `json:"success"`
	Message  string        `json:"message"`
	Customer *CustomerData `json:"customer,omitempty"`
}

type CustomerData struct {
	ID        string  `json:"id"`
	FullName  string  `json:"full_name"`
	Phone     string  `json:"phone"`
	Email     string  `json:"email"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Title     string  `json:"title"`
	Address   string  `json:"address"`
}

type AddAddressRequest struct {
	Title        string   `json:"title" binding:"required"`
	AddressLine1 string   `json:"address_line1" binding:"required"`
	AddressLine2 string   `json:"address_line2"`
	City         string   `json:"city" binding:"required"`
	State        string   `json:"state" binding:"required"`
	PostalCode   string   `json:"postal_code" binding:"required"`
	Country      string   `json:"country" binding:"required"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
}

type SignUpRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
}

type AddressResponse struct {
	Success   bool    `json:"success"`
	Message   string  `json:"message"`
	AddressID *string `json:"address_id,omitempty"`
}

type MarkDefaultRequest struct {
	AddressID string `json:"address_id" binding:"required"`
}

type Coupon struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Code          string    `json:"code"`
	DiscountType  string    `json:"discount_type"` // "percentage" or "fixed"
	DiscountValue float64   `json:"discount_value"`
	MaxDiscount   *float64  `json:"max_discount"`
	MinOrderValue float64   `json:"min_order_value"`
	MaxUsages     int       `json:"max_usages"`
	UsedCount     int       `json:"used_count"`
	ValidFrom     time.Time `json:"valid_from"`
	ValidTo       time.Time `json:"valid_to"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
}

type CouponsResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Coupons []Coupon `json:"coupons,omitempty"`
}

type ApplyCouponRequest struct {
	OrderID    string `json:"order_id" binding:"required"`
	CouponCode string `json:"coupon_code" binding:"required"`
}

type ApplyCouponResponse struct {
	Success        bool     `json:"success"`
	Message        string   `json:"message"`
	DiscountAmount *float64 `json:"discount_amount,omitempty"`
	FinalAmount    *float64 `json:"final_amount,omitempty"`
}

type OrderStatus struct {
	OrderID      string     `json:"order_id"`
	Status       string     `json:"status"`
	OrderTime    time.Time  `json:"order_time"`
	Amount       float64    `json:"amount"`
	PackByTime   *time.Time `json:"pack_by_time"`
	PaidTime     *time.Time `json:"paid_time"`
	DeliveryTime *time.Time `json:"delivery_time"`
}

type OrderStatusResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Order   *OrderStatus `json:"order,omitempty"`
}
