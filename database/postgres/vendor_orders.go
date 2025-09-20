package postgres

import (
	"context"
	"fmt"
	"qvickly/models/vendors"
	"time"

	"github.com/jackc/pgx/v5"
)

func GetVendorTodaysOrderSummary(vendorId string) (orderSummary *vendors.TodayOrderSummary, err error) {
	orderSummary = new(vendors.TodayOrderSummary)

	// Query to get order counts by status for today's orders assigned to this vendor
	query := `
        SELECT o.status, COUNT(*) 
        FROM quickkart.customer.orders o
        INNER JOIN quickkart.vendor.order_pickup_assignments opa ON o.order_id = opa.order_id
        WHERE DATE(o.order_time) = CURRENT_DATE 
        AND opa.vendor_id = $1::uuid
        GROUP BY o.status`

	rows, err := pgPool.Query(context.Background(), query, vendorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Process the results
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}

		// Map status counts to the summary struct
		switch status {
		case "pending":
			orderSummary.Pending = count
		case "accepted":
			orderSummary.Accepted = count
		case "packed":
			orderSummary.Packed = count
		case "ready":
			orderSummary.Ready = count
		case "completed":
			orderSummary.Completed = count
		case "cancelled":
			orderSummary.Cancelled = count
		case "rejected":
			orderSummary.Rejected = count
		default:
			// Handle any unexpected status values
			break
		}
	}

	// Check for any row iteration errors
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orderSummary, nil
}

func GetVendorOrderDetails(orderID string) (*vendors.OrderDetailsResponse, error) {
	ctx := context.Background()

	// Main query to get order details with delivery partner and customer info
	mainQuery := `
		SELECT 
			db.pin as delivery_partner_pin,
			db.full_name as delivery_partner_name,
			db.phone as delivery_partner_phone,
			o.pack_by_time,
			o.paid_time,
			o.delivery_time,
			o.order_id,
			c.full_name as customer_name,
			CONCAT(ca.address_line1, ', ', ca.address_line2, ', ', ca.city, ', ', ca.state, ' ', ca.postal_code) as customer_address,
			o.order_time as order_created_time,
			o.amount as order_total_amount
		FROM quickkart.customer.orders o
		LEFT JOIN quickkart.delivery.order_tracker dt ON o.order_id = dt.order_id
		LEFT JOIN quickkart.profile.delivery_boy db ON dt.delivery_boy_id = db.id
		LEFT JOIN quickkart.profile.customer c ON o.customer_id = c.id
		LEFT JOIN quickkart.profile.customer_addresses ca ON c.id = ca.customer_id AND ca.is_default = true
		WHERE o.order_id = $1`

	var response vendors.OrderDetailsResponse
	var deliveryPartnerPin, deliveryPartnerName, deliveryPartnerPhone *string
	var packByTime, paidTime, deliveryTime *time.Time

	err := pgPool.QueryRow(ctx, mainQuery, orderID).Scan(
		&deliveryPartnerPin,
		&deliveryPartnerName,
		&deliveryPartnerPhone,
		&packByTime,
		&paidTime,
		&deliveryTime,
		&response.OrderID,
		&response.CustomerName,
		&response.CustomerAddress,
		&response.OrderCreatedTime,
		&response.OrderTotalAmount,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, err
	}

	// Assign nullable fields
	response.DeliveryPartnerPin = deliveryPartnerPin
	response.DeliveryPartnerName = deliveryPartnerName
	response.DeliveryPartnerPhone = deliveryPartnerPhone
	response.PackByTime = packByTime
	response.PaidByTime = paidTime
	response.DeliveredByTime = deliveryTime

	// Get current order status from the orders table
	statusQuery := `
		SELECT status 
		FROM quickkart.customer.orders 
		WHERE order_id = $1`

	err = pgPool.QueryRow(ctx, statusQuery, orderID).Scan(&response.OrderStatus)
	if err != nil {
		return nil, err
	}

	// Get order items
	itemsQuery := `
		SELECT 
			gi.item_id,
			gi.image_url_1 as item_image_url,
			gi.title as item_name,
			coi.qty as qty_ordered
		FROM quickkart.customer.order_items coi
		JOIN quickkart.master.grocery_items gi ON coi.item_id = gi.item_id
		WHERE coi.order_id = $1`

	rows, err := pgPool.Query(ctx, itemsQuery, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []vendors.OrderItem
	for rows.Next() {
		var item vendors.OrderItem
		err := rows.Scan(
			&item.ItemID,
			&item.ItemImageURL,
			&item.ItemName,
			&item.QtyOrdered,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	response.Items = items
	return &response, nil
}

// Update order status function
func UpdateOrderStatus(orderID, newStatus string) error {
	ctx := context.Background()
	updateTime := time.Now()

	var err error

	if newStatus == "pending" {
		_, err = pgPool.Exec(ctx,
			`UPDATE quickkart.customer.orders SET status = $1, updated_at = $2, pack_by_time = $3 WHERE order_id = $4`,
			newStatus, updateTime, updateTime.Add(time.Minute*12), orderID)
	} else {
		_, err = pgPool.Exec(ctx,
			`UPDATE quickkart.customer.orders SET status = $1, updated_at = $2 WHERE order_id = $3`,
			newStatus, updateTime, orderID)
	}

	if newStatus == "picked" {
		_, err = pgPool.Exec(ctx,
			`UPDATE quickkart.vendor.order_pickup_assignments SET picked_up = true WHERE order_id = $1`,
			orderID)
	}

	// Insert into status logs
	//_, err = pgPool.Exec(ctx,
	//	`INSERT INTO orders.order_status_logs (order_id, status, changed_at)
	//		 VALUES ($1, $2, $3)
	//		 ON CONFLICT (order_id)
	//		 DO UPDATE SET
	//			status = EXCLUDED.status,
	//			changed_at = EXCLUDED.changed_at;`,
	//	orderID, newStatus, updateTime)

	return err
}

func GetVendorOrders(vendorID string, page, limit int) (*vendors.OrdersListResponse, error) {
	ctx := context.Background()

	// Calculate offset
	offset := (page - 1) * limit

	// Get total count first - count orders assigned to this vendor
	var totalCount int
	countQuery := `
        SELECT COUNT(DISTINCT o.order_id) 
        FROM quickkart.customer.orders o
        INNER JOIN quickkart.vendor.order_pickup_assignments opa ON o.order_id = opa.order_id
        WHERE opa.vendor_id = $1::uuid`
	err := pgPool.QueryRow(ctx, countQuery, vendorID).Scan(&totalCount)
	if err != nil {
		return nil, err
	}

	// Get orders assigned to this vendor
	ordersQuery := `
       SELECT 
          o.order_id::text as order_id,
          o.status as order_status,
          c.full_name as customer_name,
          o.order_time as order_time_placed,
          o.amount as total_amount,
          o.pack_by_time,
          o.delivery_time as delivered_by_time,
          o.pack_by_time as pick_by_time
       FROM quickkart.customer.orders o
       INNER JOIN quickkart.vendor.order_pickup_assignments opa ON o.order_id = opa.order_id
       LEFT JOIN quickkart.profile.customer c ON o.customer_id = c.id
       WHERE opa.vendor_id = $1::uuid
       ORDER BY o.order_time DESC
       LIMIT $2 OFFSET $3`

	rows, err := pgPool.Query(ctx, ordersQuery, vendorID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []vendors.OrderListItem
	for rows.Next() {
		var order vendors.OrderListItem
		err := rows.Scan(
			&order.OrderID,
			&order.OrderStatus,
			&order.CustomerName,
			&order.OrderTimePlaced,
			&order.TotalAmount,
			&order.PackByTime,
			&order.DeliveredByTime,
			&order.PickByTime,
		)

		if order.PackByTime == nil {
			newTime := order.OrderTimePlaced.Add(time.Minute * 12)
			order.PackByTime = &(newTime)
		}

		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Calculate pagination info
	totalPages := (totalCount + limit - 1) / limit
	hasNext := page < totalPages
	hasPrev := page > 1

	response := &vendors.OrdersListResponse{
		Success:    true,
		Data:       orders,
		Page:       page,
		Limit:      limit,
		TotalCount: totalCount,
		TotalPages: totalPages,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
	}

	return response, nil
}
