package delivery

import "time"

// DeliveryPartnerProfile represents the delivery partner profile response
type DeliveryPartnerProfile struct {
	ID                    string    `json:"id" db:"id"`
	Name                  string    `json:"name" db:"name"`
	PhoneNumber           string    `json:"phone_number" db:"phone_number"`
	ImageURL              *string   `json:"image_url" db:"image_url"`
	Rating                float64   `json:"rating" db:"rating"`
	Online                bool      `json:"online" db:"online_status"`
	IsActive              bool      `json:"is_active" db:"is_active"`
	Latitude              *float64  `json:"latitude" db:"latitude"`
	Longitude             *float64  `json:"longitude" db:"longitude"`
	AadharCardNumber      *string   `json:"aadhar_card_number" db:"aadhar_card_number"`
	DrivingLicenseNumber  *string   `json:"driving_license_number" db:"driving_license_number"`
	VehicleType           *string   `json:"vehicle_type" db:"vehicle_type"`
	VehicleNumber         *string   `json:"vehicle_number" db:"vehicle_number"`
	UPIID                 *string   `json:"upi_id" db:"upi_id"`
	EmergencyContactName  *string   `json:"emergency_contact_name" db:"emergency_contact_name"`
	EmergencyContactPhone *string   `json:"emergency_contact_phone" db:"emergency_contact_phone"`
	DateOfBirth           *string   `json:"date_of_birth" db:"date_of_birth"`
	AadharCardImageURL    *string   `json:"aadhar_card_image_url" db:"aadhar_card_image_url"`
	TotalDeliveries       int       `json:"total_deliveries" db:"total_deliveries"`
	LastLocationUpdate    time.Time `json:"last_location_update" db:"last_location_update"`
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
