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

// GetDeliveryDetails retrieves basic summary information
func GetDeliveryDetails(deliveryBoyID uuid.UUID) ([]delivery.PickupDetail, []delivery.DeliveryDetail, error) {
	completedQuery := `
		SELECT
			v.business_name AS vendor_name,
			o.pick_up_time AS pickup_time,
			v.address,
			COUNT(oi.item_id) AS total_items,
			CONCAT('₹', SUM(gi.price_retail * oi.qty)) AS total_amount,
			o.order_id,
			vpt.vendor_assignment_id,
			v.vendor_id,
			opa.picked_up,
			o.order_time
		FROM
			delivery.vendor_pickup_tracker vpt
				JOIN vendor.order_pickup_assignments opa ON vpt.vendor_assignment_id = opa.vendor_assignment_id
				JOIN profile.vendors v ON opa.vendor_id = v.vendor_id
				JOIN customer.orders o ON opa.order_id = o.order_id
				JOIN customer.order_items oi ON o.order_id = oi.order_id
				JOIN master.grocery_items gi ON gi.item_id = oi.item_id
		WHERE
			vpt.delivery_id = $1
		GROUP BY
			v.vendor_id,
			vpt.vendor_assignment_id,
			vpt.delivery_id,
			v.business_name,
			o.pick_up_time,
			o.order_id,
			vpt.delivery_id,
			opa.picked_up
		ORDER BY
			o.pick_up_time ASC; 
    `

	rows, err := pgPool.Query(context.Background(), completedQuery, deliveryBoyID)
	if err != nil {
		return nil, nil, err
	}

	var details []delivery.PickupDetail

	for rows.Next() {
		var pickupDetail delivery.PickupDetail
		err := rows.Scan(
			&pickupDetail.Name,
			&pickupDetail.PickupTime,
			&pickupDetail.Address,
			&pickupDetail.Items,
			&pickupDetail.Amount,
			&pickupDetail.OrderId,
			&pickupDetail.VendorAssignmentId,
			&pickupDetail.VendorId,
			&pickupDetail.PickedUp,
			&pickupDetail.OrderTime,
		)
		if err != nil {
			return nil, nil, err
		}
		details = append(details, pickupDetail)
	}

	completedQuery = `
		SELECT 
			o.order_id,
			o.pick_up_time,
			o.deliver_by_time AS delivery_time,
			concat(ca.city, ', ', ca.state) as address,
			c.full_name AS customer_name,
			o.status
		FROM 
			delivery.order_tracker ot
			JOIN customer.orders o ON ot.order_id = o.order_id
			JOIN profile.customer c ON o.customer_id = c.id
			JOIN profile.customer_addresses ca ON o.customer_id = ca.customer_id AND ca.is_default = true
		WHERE 
			ot.delivery_id = $1
		ORDER BY 
			o.deliver_by_time ASC, 
			o.order_id ASC;
	`

	rows, err = pgPool.Query(context.Background(), completedQuery, deliveryBoyID)
	if err != nil {
		return nil, nil, err
	}

	var dDetails []delivery.DeliveryDetail

	for rows.Next() {
		var dDetail delivery.DeliveryDetail
		err := rows.Scan(
			&dDetail.OrderId,
			&dDetail.PickUpTime,
			&dDetail.DeliverByTime,
			&dDetail.Address,
			&dDetail.Name,
			&dDetail.Status,
		)
		if err != nil {
			return nil, nil, err
		}
		dDetails = append(dDetails, dDetail)
	}

	return details, dDetails, nil
}

// GetBasicRecentOrders retrieves basic recent orders information
func GetBasicRecentOrders(deliveryBoyID uuid.UUID, limit int, statusFilter string) ([]delivery.RecentOrderResponse, error) {
	// Build the query with optional status filter
	baseQuery := `
	  SELECT
		ot.delivery_id,
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

// GetBasicRecentOrders retrieves basic recent orders information
func GetBasicAllOrders(deliveryBoyID uuid.UUID, statusFilter string) ([]delivery.RecentOrderResponse, error) {
	// Build the query with optional status filter
	baseQuery := `
	  SELECT
		ot.delivery_id,
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
	baseQuery += `ORDER BY o.order_time DESC `

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

// Helper function to get order items (also updated for your schema)
func GetDeliveryVendorItems(assignmentId uuid.UUID) ([]delivery.OrderItemSummary, delivery.PickupSummary, error) {
	itemsQuery := `
       select ot.item_id, ot.qty, gi.title, gi.price_retail, gi.image_url_1 from vendor.order_items ot
         left join master.grocery_items gi on ot.item_id = gi.item_id
         where vendor_assignment_id = $1
    `
	profileQuery := `
       select v.latitude, v.longitude, v.address, v.phone, v.owner_name, v.business_name, opa.created_at from vendor.order_pickup_assignments opa
		left join profile.vendors v on opa.vendor_id = v.vendor_id
        where opa.vendor_assignment_id = $1
    `

	var summary delivery.PickupSummary
	err := pgPool.QueryRow(context.Background(), profileQuery, assignmentId).Scan(
		&summary.Latitude,
		&summary.Longitude,
		&summary.Address,
		&summary.Phone,
		&summary.OwnerName,
		&summary.BusinessName,
		&summary.CreatedAt,
	)
	if err != nil {
		return nil, delivery.PickupSummary{}, err
	}

	rows, err := pgPool.Query(context.Background(), itemsQuery, assignmentId)
	if err != nil {
		return nil, summary, err
	}
	defer rows.Close()

	var items []delivery.OrderItemSummary

	for rows.Next() {
		var item delivery.OrderItemSummary

		err := rows.Scan(
			&item.ID,
			&item.Quantity,
			&item.Name,
			&item.Price,
			&item.ImageURL1,
		)
		if err != nil {
			return nil, summary, err
		}
		items = append(items, item)
	}

	return items, summary, nil
}

// Helper function to get order items (also updated for your schema)
func GetDeliveryCustomerItems(orderid uuid.UUID) ([]delivery.OrderItemSummary, delivery.OrderSummary, error) {
	itemsQuery := `   
		select ot.item_id, ot.qty, gi.title, gi.price_retail, gi.image_url_1 from customer.order_items ot
			left join master.grocery_items gi on ot.item_id = gi.item_id
			where order_id = $1
    `

	profileQuery := `
		select order_time, deliver_by_time, delivery_time, instructions,
			   amount, status, phone, full_name, concat(address_line1, address_line2, postal_code), 
			   latitude, longitude from customer.orders
			left join profile.customer c on orders.customer_id = c.id
			left join profile.customer_addresses ca on c.id = ca.customer_id
        	where order_id = $1
    `

	var summary delivery.OrderSummary
	err := pgPool.QueryRow(context.Background(), profileQuery, orderid).Scan(
		&summary.OrderTime,
		&summary.DeliverByTime,
		&summary.DeliveryTime,
		&summary.Instructions,
		&summary.Amount,
		&summary.Status,
		&summary.Phone,
		&summary.FullName,
		&summary.Address,
		&summary.Latitude,
		&summary.Longitude,
	)
	if err != nil {
		return nil, delivery.OrderSummary{}, err
	}

	rows, err := pgPool.Query(context.Background(), itemsQuery, orderid)
	if err != nil {
		return nil, summary, err
	}
	defer rows.Close()

	var items []delivery.OrderItemSummary

	for rows.Next() {
		var item delivery.OrderItemSummary

		err := rows.Scan(
			&item.ID,
			&item.Quantity,
			&item.Name,
			&item.Price,
			&item.ImageURL1,
		)
		if err != nil {
			return nil, summary, err
		}
		items = append(items, item)
	}

	return items, summary, nil
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
