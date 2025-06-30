package postgres

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"qvickly/models/delivery"
	"strconv"
	"time"
)

// GetBasicOrdersSummary retrieves basic summary information
func GetBasicOrdersSummary(deliveryBoyID uuid.UUID) (*delivery.OrdersSummaryResponse, error) {
	// Query for completed orders count and total earnings
	// Using 'delivered' status as that's the equivalent of 'completed' in your schema
	completedQuery := `
       SELECT 
          COUNT(*) as completed_count,
          COALESCE(SUM(o.amount), 0) as total_earnings
       FROM customer.orders o
       INNER JOIN delivery.order_tracker ot ON o.order_id = ot.order_id
       WHERE ot.delivery_boy_id = $1 
       AND LOWER(o.status) like 'delivered'
    `

	var completedCount int
	var totalEarnings float64

	err := pgPool.QueryRow(context.Background(), completedQuery, deliveryBoyID).Scan(&completedCount, &totalEarnings)
	if err != nil {
		return nil, err
	}

	// Query for active orders count
	// Active orders are those that are accepted, packed, or picked (in progress)
	activeQuery := `
       SELECT COUNT(*) 
       FROM customer.orders o
       INNER JOIN delivery.order_tracker ot ON o.order_id = ot.order_id
       WHERE ot.delivery_boy_id = $1 
       AND LOWER(o.status) IN ('accepted', 'packed', 'picked')
    `

	var activeCount int
	err = pgPool.QueryRow(context.Background(), activeQuery, deliveryBoyID).Scan(&activeCount)
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
func GetBasicRecentOrders(deliveryBoyID uuid.UUID, limit int, statusFilter string) ([]delivery.RecentOrderResponse, error) {
	// Build the query with optional status filter
	baseQuery := `
	  SELECT
		o.order_id,
		o.status,
		o.amount,
		o.updated_at as last_status_updated_time,
		(SELECT COUNT(*)
		 FROM customer.order_items oi
		 WHERE oi.order_id = o.order_id) as items_count
		FROM customer.orders o
				 left JOIN delivery.order_tracker ot ON o.order_id = ot.order_id
		WHERE ot.delivery_boy_id = $1 
    `

	var args []interface{}
	argIndex := 2

	args = append(args, deliveryBoyID)

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
func GetOrderDetail(orderID uuid.UUID, deliveryBoyID uuid.UUID) (*delivery.OrderDetailResponse, error) {
	// Main order query with all required details
	orderQuery := `
       SELECT 
          o.order_id,
          o.status,
          o.amount as delivery_fee,
          c.full_name as customer_name,
          CONCAT(ca.address_line1, CASE WHEN ca.address_line2 IS NOT NULL THEN ', ' || ca.address_line2 ELSE '' END, 
                 ', ', ca.city, ', ', ca.state, ' ', ca.postal_code) as delivery_address,
          o.instructions as delivery_instructions,
          o.amount as total_amount,
          v.business_name as store_name,
          ca.latitude as delivery_latitude,
          ca.longitude as delivery_longitude,
          -- Use order_time as accepted_at since we don't have status logs
          o.order_time as accepted_at,
          -- Use delivery_time as delivered_at
          o.delivery_time as delivered_at,
          -- Additional fields that might be useful
          o.pack_by_time,
          o.pick_up_time,
          o.paid_time
       FROM customer.orders o
       INNER JOIN delivery.order_tracker ot ON o.order_id = ot.order_id
       INNER JOIN profile.customer c ON o.customer_id = c.id
       LEFT JOIN profile.customer_addresses ca ON c.id = ca.customer_id AND ca.is_default = true
       LEFT JOIN vendor.order_pickup_assignments opa ON o.order_id = opa.order_id
       LEFT JOIN profile.vendors v ON opa.vendor_id = v.vendor_id
       WHERE o.order_id = $1 
       AND ot.delivery_boy_id = $2
    `

	var detail delivery.OrderDetailResponse
	var acceptedAt sql.NullTime
	var deliveredAt sql.NullTime
	var packByTime sql.NullTime
	var pickUpTime sql.NullTime
	var paidTime sql.NullTime
	var deliveryInstructions sql.NullString
	var deliveryLatitude sql.NullFloat64
	var deliveryLongitude sql.NullFloat64
	var storeName sql.NullString

	err := pgPool.QueryRow(context.Background(), orderQuery, orderID, deliveryBoyID).Scan(
		&detail.ID,
		&detail.Status,
		&detail.DeliveryFee,
		&detail.CustomerName,
		&detail.DeliveryAddress,
		&deliveryInstructions,
		&detail.ItemsValue,
		&storeName,
		&deliveryLatitude,
		&deliveryLongitude,
		&acceptedAt,
		&deliveredAt,
		&packByTime,
		&pickUpTime,
		&paidTime,
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

	if storeName.Valid {
		detail.StoreName = storeName.String
	}

	// Calculate bonus (example: 10% of delivery fee for completed orders during peak hours)
	var bonusTime time.Time
	if acceptedAt.Valid {
		bonusTime = acceptedAt.Time
	}
	detail.Bonus = calculateBonus(detail.DeliveryFee, detail.Status, bonusTime)

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

// Helper function to get order items (also updated for your schema)
func getOrderItems(orderID uuid.UUID) ([]delivery.OrderItem, error) {
	itemsQuery := `
       SELECT 
          gi.title,
          gi.description,
          oi.qty,
          gi.price_retail as unit_price,
          (oi.qty * gi.price_retail) as total_price,
          gi.image_url_1 as image_url
       FROM customer.order_items oi
       INNER JOIN master.grocery_items gi ON oi.item_id = gi.item_id
       WHERE oi.order_id = $1
       ORDER BY gi.title
    `

	rows, err := pgPool.Query(context.Background(), itemsQuery, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []delivery.OrderItem

	for rows.Next() {
		var item delivery.OrderItem
		var description sql.NullString
		var imageURL sql.NullString

		err := rows.Scan(
			&item.Name,
			&description,
			&item.Quantity,
			&item.UnitPrice,
			&item.TotalPrice,
			&imageURL,
		)
		if err != nil {
			return nil, err
		}

		if description.Valid {
			item.Description = description.String
		}

		if imageURL.Valid {
			item.ImageURL1 = imageURL.String
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
