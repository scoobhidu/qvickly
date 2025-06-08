package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

// ProcessLocationUpdate handles the main location update logic
func ProcessLocationUpdate(partnerID uuid.UUID, latitude, longitude float64) error {
	now := time.Now()

	// Start transaction
	tx, err := pgClient.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	// Update delivery partner location and set online
	updateQuery := `
		UPDATE delivery_partners.delivery_partners 
		SET latitude = $1,
		    longitude = $2,
		    online_status = true,
		    last_location_update = $3,
		    updated_at = $3
		WHERE id = $4 AND is_active = true
		RETURNING name
	`

	var partnerName string
	err = tx.QueryRow(context.Background(), updateQuery, latitude, longitude, now, partnerID).Scan(&partnerName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("partner_not_found")
		}
		return err
	}

	// Commit transaction
	if err = tx.Commit(context.Background()); err != nil {
		return err
	}

	return nil
}
