package postgres

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"qvickly/models/delivery"
)

func GetDeliveryPartnerProfileDetails(partnerID uuid.UUID) (profile delivery.DeliveryPartnerProfile, dateOfBirth sql.NullString, err error) {
	query := `
		SELECT 
			id,
			name,
			phone_number,
			image_url,
			rating,
			online_status,
			is_active,
			latitude,
			longitude,
			aadhar_card_number,
			driving_license_number,
			vehicle_type,
			vehicle_number,
			upi_id,
			emergency_contact_name,
			emergency_contact_phone,
			CASE 
				WHEN date_of_birth IS NOT NULL THEN date_of_birth::text 
				ELSE NULL 
			END as date_of_birth,
			aadhar_card_image_url,
			total_deliveries,
			last_location_update,
			created_at,
			updated_at
		FROM delivery_partners.delivery_partners 
		WHERE id = $1
	`

	err = pgClient.QueryRow(context.Background(), query, partnerID).Scan(
		&profile.ID,
		&profile.Name,
		&profile.PhoneNumber,
		&profile.ImageURL,
		&profile.Rating,
		&profile.Online,
		&profile.IsActive,
		&profile.Latitude,
		&profile.Longitude,
		&profile.AadharCardNumber,
		&profile.DrivingLicenseNumber,
		&profile.VehicleType,
		&profile.VehicleNumber,
		&profile.UPIID,
		&profile.EmergencyContactName,
		&profile.EmergencyContactPhone,
		&dateOfBirth,
		&profile.AadharCardImageURL,
		&profile.TotalDeliveries,
		&profile.LastLocationUpdate,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	return
}
