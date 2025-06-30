package postgres

import (
	"context"
	"database/sql"
	"qvickly/models/delivery"
)

func GetDeliveryPartnerProfileDetails(phone, password string) (profile delivery.DeliveryPartnerProfile, dateOfBirth sql.NullString, err error) {
	query := `
		SELECT 
			id,
			full_name,
			phone,
			image_url,
			rating,
			is_active,
			aadhar,
			upi_id,
			emergency_contact_name,
			emergency_contact_phone,
			date_of_birth,
			created_at,
			updated_at
		FROM quickkart.profile.delivery_boy 
		WHERE phone = $1 and password = $2
	`

	err = pgPool.QueryRow(context.Background(), query, phone, password).Scan(
		&profile.ID,
		&profile.Name,
		&profile.PhoneNumber,
		&profile.ImageURL,
		&profile.Rating,
		&profile.IsActive,
		&profile.AadharCardNumber,
		&profile.UPIID,
		&profile.EmergencyContactName,
		&profile.EmergencyContactPhone,
		&dateOfBirth,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	return
}
