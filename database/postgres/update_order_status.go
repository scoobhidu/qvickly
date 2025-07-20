package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"qvickly/models/delivery"
	"time"
)

// ProcessPickupVerification handles the main pickup verification logic
func ProcessPickupVerification(orderID uuid.UUID, deliveryPartnerID uuid.UUID, providedPin int) (*delivery.VerifyPickupResponse, error) {
	// Start database transaction
	tx, err := pgPool.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())

	// Get order details and verify assignment
	orderQuery := `
		SELECT 
			o.id,
			o.status,
			o.pickup_pin,
			o.customer_name,
			dp.name as delivery_partner_name,
			va.business_name as vendor_name,
			(SELECT COUNT(*) FROM orders.order_items WHERE order_id = o.id) as items_count
		FROM orders.orders o
		JOIN delivery_partners.delivery_partners dp ON o.delivery_partner_id = dp.id
		JOIN vendor_accounts.vendor_accounts va ON o.account_id = va.id
		WHERE o.id = $1 AND o.delivery_partner_id = $2
	`

	var orderInfo struct {
		ID                  int
		Status              string
		PickupPin           sql.NullInt32
		CustomerName        string
		DeliveryPartnerName string
		VendorName          string
		ItemsCount          int
	}

	err = tx.QueryRow(context.Background(), orderQuery, orderID, deliveryPartnerID).Scan(
		&orderInfo.ID,
		&orderInfo.Status,
		&orderInfo.PickupPin,
		&orderInfo.CustomerName,
		&orderInfo.DeliveryPartnerName,
		&orderInfo.VendorName,
		&orderInfo.ItemsCount,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order_not_found")
		}
		return nil, err
	}

	// Check if order status is valid for pickup (should be 'ready')
	if orderInfo.Status != "ready" {
		if orderInfo.Status == "completed" {
			return nil, fmt.Errorf("already_picked_up")
		}
		return nil, fmt.Errorf("invalid_status")
	}

	// Verify PIN
	if !orderInfo.PickupPin.Valid {
		return nil, fmt.Errorf("no_pin_set")
	}

	if int(orderInfo.PickupPin.Int32) != providedPin {
		return nil, fmt.Errorf("wrong_pin")
	}

	// Update order status to indicate pickup
	now := time.Now()
	newStatus := "picked_up" // You might want to add this status to your enum, or use existing status

	// If "picked_up" is not in your status enum, use an existing status like "ready" -> "completed"
	// For this example, I'll assume we transition from "ready" to "completed" after pickup
	newStatus = "completed" // Or create a new intermediate status

	updateOrderQuery := `
		UPDATE orders.orders 
		SET status = $1, 
		    delivered_by_time = $2,
		    updated_at = $2
		WHERE id = $3
	`

	_, err = tx.Exec(context.Background(), updateOrderQuery, newStatus, now, orderID)
	if err != nil {
		return nil, err
	}

	// Insert status log entry
	insertStatusLogQuery := `
		INSERT INTO orders.order_status_logs (order_id, status, changed_at)
		VALUES ($1, $2, $3)
	`

	_, err = tx.Exec(context.Background(), insertStatusLogQuery, orderID, newStatus, now)
	if err != nil {
		return nil, err
	}

	// Optional: Update delivery partner stats
	updatePartnerStatsQuery := `
		UPDATE delivery_partners.delivery_partners  
		SET total_deliveries = total_deliveries + 1,
		    last_location_update = $2,
		    updated_at = $2
		WHERE id = $1
	`

	_, err = tx.Exec(context.Background(), updatePartnerStatsQuery, deliveryPartnerID, now)
	if err != nil {
		// Log error but don't fail the transaction
		// This is optional data
	}

	// Commit transaction
	if err = tx.Commit(context.Background()); err != nil {
		return nil, err
	}

	// Prepare response
	response := &delivery.VerifyPickupResponse{
		Success:         true,
		Message:         "Pickup verified successfully. Order status updated.",
		OrderID:         orderID,
		NewStatus:       newStatus,
		VerifiedAt:      now,
		DeliveryPartner: orderInfo.DeliveryPartnerName,
		VendorName:      orderInfo.VendorName,
		CustomerName:    orderInfo.CustomerName,
		ItemsCount:      orderInfo.ItemsCount,
	}

	return response, nil
}

// GetCurrentOrderStatus gets the current status of an order for error responses
func GetCurrentOrderStatus(orderID uuid.UUID) string {
	var status string
	query := `SELECT status FROM orders.orders WHERE id = $1`

	err := pgPool.QueryRow(context.Background(), query, orderID).Scan(&status)
	if err != nil {
		return "unknown"
	}

	return status
}
