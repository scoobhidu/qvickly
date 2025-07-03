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
func GetDeliveryDetails(deliveryBoyID uuid.UUID) (*delivery.DeliveryDetailsResponse, error) {
	var vendor_assignment_id, order_id, customer_id, delivery_id, vendor_id uuid.UUID
	var order_time, pack_by_time, pick_up_time, deliver_by_time, paid_time, delivery_time sql.NullTime
	var instructions, status, c_phone, full_name, title, address_line1, address_line2,
		city, state, postal_code, country, account_type, business_name, owner_name, v_phone, address string

	var amount, c_latitude, c_longitude, v_latitude, v_longitude float64

	completedQuery := `
     select
    opa.vendor_assignment_id,
    opa.vendor_id,
    opa.order_id,
    ca.customer_id,
    o.order_time,
    o.pack_by_time,
    o.pick_up_time,
    o.deliver_by_time,
    o.paid_time,
    o.delivery_time,
    o.instructions,
    o.amount,
    o.status,
    ot.delivery_id,
    c.phone,
    c.full_name,
    ca.title,
    ca.address_line1,
    ca.address_line2,
    ca.city,
    ca.state,
    ca.postal_code,
    ca.country,
    ca.latitude,
    ca.longitude,
    v.account_type,
    v.business_name,
    v.owner_name,
    v.phone,
    v.address,
    v.latitude,
    v.longitude
    from vendor.order_pickup_assignments opa
    left join customer.orders o on opa.order_id = o.order_id
    left join delivery.order_tracker ot on o.order_id = ot.order_id
    left join delivery.vendor_pickup_tracker vpt on ot.delivery_id = vpt.delivery_id
    left join profile.customer c on o.customer_id = c.id
    left join profile.customer_addresses ca on c.id = ca.customer_id
    left join profile.vendors v on opa.vendor_id = v.vendor_id

    where ot.delivery_id = $1::uuid; 
    `

	err := pgPool.QueryRow(context.Background(), completedQuery, deliveryBoyID).Scan(
		&vendor_assignment_id, &vendor_id, &order_id, &customer_id, &order_time, &pack_by_time, &pick_up_time,
		&deliver_by_time, &paid_time, &delivery_time, &instructions, &amount, &status, &delivery_id, &c_phone,
		&full_name, &title, &address_line1, &address_line2, &city, &state, &postal_code, &country, &c_latitude,
		&c_longitude, &account_type, &business_name, &owner_name, &v_phone, &address, &v_latitude,
		&v_longitude,
	)
	if err != nil {
		return nil, err
	}
	var orT, paT, piT, deliverByT, paidT, deliveryT time.Time

	if order_time.Valid {
		orT = order_time.Time
	}
	if pack_by_time.Valid {
		paT = pack_by_time.Time
	}
	if pick_up_time.Valid {
		piT = pick_up_time.Time
	}
	if deliver_by_time.Valid {
		deliverByT = deliver_by_time.Time
	}
	if paid_time.Valid {
		paT = paid_time.Time
	}
	if delivery_time.Valid {
		deliveryT = delivery_time.Time
	}

	return &delivery.DeliveryDetailsResponse{
		VendorAssignmentId: vendor_assignment_id,
		VendorId:           vendor_id,
		OrderId:            order_id,
		CustomerId:         customer_id,
		DeliveryId:         delivery_id,
		OrderTime:          &orT,
		PackByTime:         &paT,
		PickUpTime:         &piT,
		DeliverByTime:      &deliverByT,
		PaidTime:           &paidT,
		DeliveryTime:       &deliveryT,
		Instructions:       instructions,
		Amount:             amount,
		Status:             status,
		Phone:              c_phone,
		CustomerName:       full_name,
		Title:              title,
		Address1:           address_line1,
		Address2:           address_line2,
		City:               city,
		State:              state,
		PostalCode:         postal_code,
		Country:            country,
		CLatitude:          c_latitude,
		CLongitude:         c_longitude,
		VendorType:         account_type,
		BusinessName:       business_name,
		OwnerName:          owner_name,
		VPhone:             v_phone,
		VAddress:           address,
		VLatitude:          v_latitude,
		VLongitude:         v_longitude,
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

// Helper function to get order items (also updated for your schema)
func GetDeliveryVendorItems(assignmentId uuid.UUID) ([]delivery.OrderItemSummary, error) {
	itemsQuery := `
       select ot.item_id, ot.qty, gi.title, gi.price_retail, gi.image_url_1 from vendor.order_items ot
         left join master.grocery_items gi on ot.item_id = gi.item_id
         where vendor_assignment_id = $1
    `

	rows, err := pgPool.Query(context.Background(), itemsQuery, assignmentId)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// Helper function to get order items (also updated for your schema)
func GetDeliveryCustomerItems(orderid uuid.UUID) ([]delivery.OrderItemSummary, error) {
	itemsQuery := `   
		select ot.item_id, ot.qty, gi.title, gi.price_retail, gi.image_url_1 from customer.order_items ot
			left join master.grocery_items gi on ot.item_id = gi.item_id
			where order_id = $1
    `

	rows, err := pgPool.Query(context.Background(), itemsQuery, orderid)
	if err != nil {
		return nil, err
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
			return nil, err
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
