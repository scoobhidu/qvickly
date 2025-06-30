package delivery

import "time"

// DeliveryPartnerProfile represents the delivery partner profile response
type DeliveryPartnerProfile struct {
	ID                    string    `json:"id" db:"id"`
	Name                  string    `json:"name" db:"name"`
	PhoneNumber           string    `json:"phone_number" db:"phone_number"`
	ImageURL              *string   `json:"image_url" db:"image_url"`
	Rating                float64   `json:"rating" db:"rating"`
	IsActive              bool      `json:"is_active" db:"is_active"`
	AadharCardNumber      *string   `json:"aadhar_card_number" db:"aadhar_card_number"`
	UPIID                 *string   `json:"upi_id" db:"upi_id"`
	EmergencyContactName  *string   `json:"emergency_contact_name" db:"emergency_contact_name"`
	EmergencyContactPhone *string   `json:"emergency_contact_phone" db:"emergency_contact_phone"`
	DateOfBirth           *string   `json:"date_of_birth" db:"date_of_birth"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" db:"updated_at"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// SuccessResponse represents a success response wrapper
type DeliveryProfileDetailsSuccessResponse struct {
	Success bool                   `json:"success"`
	Data    DeliveryPartnerProfile `json:"data"`
	Message string                 `json:"message"`
}
