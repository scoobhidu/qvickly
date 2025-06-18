package postgres

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"qvickly/models/delivery"
	"strconv"
	"time"
)

// getBasicOrdersSummary retrieves basic summary information
func GetBasicOrdersSummary(partnerID uuid.UUID) (*delivery.OrdersSummaryResponse, error) {
	// Query for completed orders count and total earnings
	completedQuery := `
		SELECT 
			COUNT(*) as completed_count,
			COALESCE(SUM(delivery_fee), 0) as total_earnings
		FROM orders.orders 
		WHERE delivery_partner_id = $1 
		AND status = 'completed'
	`

	var completedCount int
	var totalEarnings float64

	err := pgPool.QueryRow(context.Background(), completedQuery, partnerID).Scan(&completedCount, &totalEarnings)
	if err != nil {
		return nil, err
	}

	// Query for active orders count
	activeQuery := `
		SELECT COUNT(*) 
		FROM orders.orders 
		WHERE delivery_partner_id = $1 
		AND status IN ('accepted', 'packed', 'ready')
	`

	var activeCount int
	err = pgPool.QueryRow(context.Background(), activeQuery, partnerID).Scan(&activeCount)
	if err != nil {
		return nil, err
	}

	return &delivery.OrdersSummaryResponse{
		Completed:    completedCount,
		Earnings:     totalEarnings,
		ActiveOrders: activeCount,
	}, nil
}

// GetBasicRecentOrders retrieves basic recent orders information
func GetBasicRecentOrders(partnerID uuid.UUID, limit int, statusFilter string) ([]delivery.RecentOrderResponse, error) {
	// Build the query with optional status filter
	baseQuery := `
		SELECT 
			o.id,
			o.status,
			o.delivery_fee,
			COALESCE(
				(SELECT osl.changed_at 
				 FROM orders.order_status_logs osl 
				 WHERE osl.order_id = o.id), 
				o.updated_at
			) as last_status_updated_time,
			(SELECT COUNT(*) 
			 FROM orders.order_items oi 
			 WHERE oi.order_id = o.id) as items_count
		FROM orders.orders o
		WHERE o.delivery_partner_id = $1 
	`

	var args []interface{}
	argIndex := 2

	args = append(args, partnerID)

	// Add status filter if provided
	if statusFilter != "" {
		baseQuery += " AND o.status = $" + strconv.Itoa(argIndex)
		args = append(args, statusFilter)
		argIndex++
	}

	// Add ordering and limit
	baseQuery += ` 
		ORDER BY o.order_time DESC 
		LIMIT $` + strconv.Itoa(argIndex)
	args = append(args, limit)

	rows, err := pgPool.Query(context.Background(), baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []delivery.RecentOrderResponse

	for rows.Next() {
		var order delivery.RecentOrderResponse

		err := rows.Scan(
			&order.ID,
			&order.Status,
			&order.Earnings,
			&order.LastStatusUpdatedTime,
			&order.Items,
		)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

// GetOrderDetail retrieves comprehensive order details for delivery partner
func GetOrderDetail(orderID int, partnerID uuid.UUID) (*delivery.OrderDetailResponse, error) {
	// Main order query with all required details
	orderQuery := `
		SELECT 
			o.id,
			o.status,
			o.delivery_fee,
			o.customer_name,
			o.location as delivery_address,
			o.delivery_instructions,
			o.total_amount,
			va.business_name as store_name,
			addr.latitude as delivery_latitude,
			addr.longitude as delivery_longitude,
			-- Get accepted_at from status logs
			(SELECT changed_at FROM orders.order_status_logs 
			 WHERE order_id = o.id AND status = 'accepted' 
			 ORDER BY changed_at ASC LIMIT 1) as accepted_at,
			-- Get delivered_at from delivered_by_time or status logs
			COALESCE(
				o.delivered_by_time,
				(SELECT changed_at FROM orders.order_status_logs 
				 WHERE order_id = o.id AND status = 'completed' 
				 ORDER BY changed_at DESC LIMIT 1)
			) as delivered_at
		FROM orders.orders o
		JOIN vendor_accounts.vendor_accounts va ON o.account_id = va.id
		LEFT JOIN user_profile.addresses addr ON o.customer_address_id = addr.id
		WHERE o.id = $1 
		AND o.delivery_partner_id = $2
	`

	var detail delivery.OrderDetailResponse
	var acceptedAt sql.NullTime
	var deliveredAt sql.NullTime
	var deliveryInstructions sql.NullString
	var deliveryLatitude sql.NullFloat64
	var deliveryLongitude sql.NullFloat64

	err := pgPool.QueryRow(context.Background(), orderQuery, orderID, partnerID).Scan(
		&detail.ID,
		&detail.Status,
		&detail.DeliveryFee,
		&detail.CustomerName,
		&detail.DeliveryAddress,
		&deliveryInstructions,
		&detail.ItemsValue,
		&detail.StoreName,
		&deliveryLatitude,
		&deliveryLongitude,
		&acceptedAt,
		&deliveredAt,
	)

	if err != nil {
		return nil, err
	}

	// Handle nullable fields
	if acceptedAt.Valid {
		detail.AcceptedAt = &acceptedAt.Time
	}

	if deliveredAt.Valid {
		detail.DeliveredAt = &deliveredAt.Time
	}

	if deliveryInstructions.Valid {
		detail.DeliveryInstruction = &deliveryInstructions.String
	}

	if deliveryLatitude.Valid {
		detail.DeliveryLatitude = &deliveryLatitude.Float64
	}

	if deliveryLongitude.Valid {
		detail.DeliveryLongitude = &deliveryLongitude.Float64
	}

	// Calculate bonus (example: 10% of delivery fee for completed orders during peak hours)
	detail.Bonus = calculateBonus(detail.DeliveryFee, detail.Status, acceptedAt.Time)

	// Calculate total earning (delivery fee + bonus)
	detail.Earning = detail.DeliveryFee + detail.Bonus

	// Get order items
	items, err := getOrderItems(orderID)
	if err != nil {
		return nil, err
	}
	detail.Items = items

	return &detail, nil
}

// getOrderItems retrieves all items for the order
func getOrderItems(orderID int) ([]delivery.OrderItemDetail, error) {
	itemsQuery := `
		SELECT 
			oi.item_id,
			vi.name as item_name,
			oi.quantity,
			-- Try to get image from item_images first, then from vendor_items
			COALESCE(
				(SELECT image_url FROM vendor_items.item_images 
				 WHERE item_id = oi.item_id AND position = 1 LIMIT 1),
				-- If no vendor item images, try qvickly products
				''
-- 				(SELECT image_url FROM qvickly_grocery_products.items 
-- 				 WHERE id = oi.item_id LIMIT 1)
			) as image_url
		FROM orders.order_items oi
		JOIN vendor_items.items vi ON oi.item_id = vi.id
		WHERE oi.order_id = $1
		ORDER BY oi.id
	`

	rows, err := pgPool.Query(context.Background(), itemsQuery, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []delivery.OrderItemDetail

	for rows.Next() {
		var item delivery.OrderItemDetail
		var imageURL sql.NullString

		err := rows.Scan(
			&item.ID,
			&item.Label,
			&item.Qty,
			&imageURL,
		)
		if err != nil {
			return nil, err
		}

		if imageURL.Valid {
			item.ImageURL = &imageURL.String
		}

		items = append(items, item)
	}

	return items, nil
}

// calculateBonus calculates bonus based on various factors
func calculateBonus(deliveryFee float64, status string, acceptedAt time.Time) float64 {
	bonus := 0.0

	// Only calculate bonus for completed orders
	if status != "completed" {
		return bonus
	}

	// Peak hour bonus (6-9 AM and 6-9 PM)
	hour := acceptedAt.Hour()
	if (hour >= 6 && hour <= 9) || (hour >= 18 && hour <= 21) {
		bonus += deliveryFee * 0.1 // 10% peak hour bonus
	}

	// Weekend bonus (Saturday and Sunday)
	if acceptedAt.Weekday() == time.Saturday || acceptedAt.Weekday() == time.Sunday {
		bonus += deliveryFee * 0.05 // 5% weekend bonus
	}

	// High value order bonus (delivery fee > ₹30)
	if deliveryFee > 30.0 {
		bonus += 5.0 // ₹5 high value bonus
	}

	return bonus
}
