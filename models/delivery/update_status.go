package delivery

import (
	"github.com/google/uuid"
	"time"
)

// VerifyPickupRequest represents the pickup verification request
type VerifyPickupRequest struct {
	Pin int `json:"pin" binding:"required" example:"1234"`
}

// VerifyPickupResponse represents the pickup verification response
type VerifyPickupResponse struct {
	Success         bool      `json:"success"`
	Message         string    `json:"message"`
	OrderID         uuid.UUID `json:"order_id"`
	NewStatus       string    `json:"new_status"`
	VerifiedAt      time.Time `json:"verified_at"`
	DeliveryPartner string    `json:"delivery_partner"`
	VendorName      string    `json:"vendor_name"`
	CustomerName    string    `json:"customer_name"`
	ItemsCount      int       `json:"items_count"`
}

// PickupErrorResponse represents pickup verification error
type PickupErrorResponse struct {
	Success       bool      `json:"success"`
	Error         string    `json:"error"`
	Message       string    `json:"message"`
	Code          int       `json:"code"`
	OrderID       uuid.UUID `json:"order_id,omitempty"`
	CurrentStatus *string   `json:"current_status,omitempty"`
}
