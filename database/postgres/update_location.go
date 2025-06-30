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
func ProcessLocationUpdate(deliveryBoyID uuid.UUID, latitude, longitude float64) error {
	now := time.Now()

	// Start transaction
	tx, err := pgPool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	updateQuery := `
       UPDATE profile.delivery_boy 
       SET latitude = $1,
           longitude = $2,
           is_active = true,
           updated_at = $3
       WHERE id = $4 AND is_active = true
       RETURNING full_name
    `

	var deliveryBoyName string
	err = tx.QueryRow(context.Background(), updateQuery, latitude, longitude, now, deliveryBoyID).Scan(&deliveryBoyName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("delivery_boy_not_found")
		}
		return err
	}

	// Commit transaction
	if err = tx.Commit(context.Background()); err != nil {
		return err
	}

	return nil
}
