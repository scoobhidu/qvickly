package vendors

import "time"

// VendorProfileDetails represents basic vendor profile information for quick display
type VendorProfileDetails struct {
	VendorId        string `json:"vendor_id" db:"vendor_id" example:"88888888-9999-aaaa-bbbb-cccccccccccc"`
	ImageS3URL      string `json:"image_s3_url" db:"image_url" example:"https://my-bucket.s3.amazonaws.com/vendors/store-123.jpg" description:"S3 URL for the vendor's store image/logo"`
	OwnerName       string `json:"owner_name" db:"owner_name" example:"John Smith" description:"Full name of the store owner"`
	StoreName       string `json:"store_name" db:"business_name" example:"Smith's Fresh Market" description:"Display name of the store/business"`
	StoreLiveStatus bool   `json:"live_status" db:"live_status" example:"true" description:"Whether the store is currently open and accepting orders"`
}

// CompleteVendorProfile represents comprehensive vendor profile information
type GetVendorProfileRequestBody struct {
	Phone    string `json:"phone" db:"phone" example:"9876543211" description:"Primary contact phone number for the vendor"`
	Password string `json:"password" db:"password" example:"c4a538ea019b7a" description:"Password for the vendor"`
}

// CompleteVendorProfile represents comprehensive vendor profile information
type CompleteVendorProfile struct {
	Phone           string    `json:"phone" db:"phone" example:"+1234567890" description:"Primary contact phone number for the vendor"`
	AccountType     string    `json:"account_type" db:"account_type" example:"premium" enums:"basic,premium,enterprise" description:"Type of vendor account subscription"`
	BusinessName    string    `json:"business_name" db:"business_name" example:"Smith's Fresh Market LLC" description:"Official registered business name"`
	OwnerName       string    `json:"owner_name" db:"owner_name" example:"John Smith" description:"Full name of the business owner"`
	Email           string    `json:"email" db:"email" example:"john.smith@freshmarket.com" format:"email" description:"Primary email address for business communications"`
	Address         string    `json:"address" db:"address" example:"123 Main Street, Downtown, New York, NY 10001" description:"Complete physical address of the store"`
	Latitude        float64   `json:"latitude" db:"latitude" example:"40.7128" minimum:"-90" maximum:"90" description:"Geographic latitude coordinate of the store location"`
	Longitude       float64   `json:"longitude" db:"longitude" example:"-74.0060" minimum:"-180" maximum:"180" description:"Geographic longitude coordinate of the store location"`
	GSTIN           string    `json:"gstin" db:"gstin" example:"22AAAAA0000A1Z5" description:"Goods and Services Tax Identification Number (for Indian businesses)"`
	OpeningTime     time.Time `json:"opening_time" db:"opening_time" example:"2000-01-01T09:00:00Z" description:"Daily store opening time (time component only, date is ignored)"`
	ClosingTime     time.Time `json:"closing_time" db:"closing_time" example:"2000-01-01T21:00:00Z" description:"Daily store closing time (time component only, date is ignored)"`
	ImageS3URL      string    `json:"image_s3_url" db:"image_url" example:"https://my-bucket.s3.amazonaws.com/vendors/store-123.jpg" description:"S3 URL for the vendor's store image/logo"`
	StoreLiveStatus bool      `json:"live_status" db:"live_status" example:"true" description:"Whether the store is currently open and accepting orders"`
}

// Request models for vendor profile operations
type UpdateVendorProfileRequest struct {
	BusinessName *string    `json:"business_name,omitempty" example:"Smith's Fresh Market LLC" description:"Updated business name"`
	OwnerName    *string    `json:"owner_name,omitempty" example:"John Smith" description:"Updated owner name"`
	Email        *string    `json:"email,omitempty" format:"email" example:"john.smith@freshmarket.com" description:"Updated email address"`
	Address      *string    `json:"address,omitempty" example:"123 Main Street, Downtown, New York, NY 10001" description:"Updated store address"`
	Latitude     *float64   `json:"latitude,omitempty" example:"40.7128" minimum:"-90" maximum:"90" description:"Updated latitude coordinate"`
	Longitude    *float64   `json:"longitude,omitempty" example:"-74.0060" minimum:"-180" maximum:"180" description:"Updated longitude coordinate"`
	GSTIN        *string    `json:"gstin,omitempty" example:"22AAAAA0000A1Z5" description:"Updated GSTIN number"`
	OpeningTime  *time.Time `json:"opening_time,omitempty" example:"2000-01-01T09:00:00Z" description:"Updated opening time"`
	ClosingTime  *time.Time `json:"closing_time,omitempty" example:"2000-01-01T21:00:00Z" description:"Updated closing time"`
}

type UpdateLiveStatusRequest struct {
	LiveStatus bool `json:"live_status" example:"true" description:"Set store live status (true = open, false = closed)"`
}

// Response models
type VendorProfileResponse struct {
	Success bool                  `json:"success" example:"true"`
	Data    CompleteVendorProfile `json:"data"`
}

type VendorProfileDetailsResponse struct {
	Success bool                 `json:"success" example:"true"`
	Data    VendorProfileDetails `json:"data"`
}
