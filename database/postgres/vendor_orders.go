package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"qvickly/models/vendors"
	"time"
)

func GetVendorTodaysOrderSummary(vendorId string) (orderSummary *vendors.TodayOrderSummary, err error) {
	orderSummary = new(vendors.TodayOrderSummary)
	rows, err := pgClient.Query(context.Background(), `select distinct(status), count(*) from postgres.orders.orders where Date(order_time) = CURRENT_DATE and account_id=$1::uuid group by status`, vendorId)
	defer rows.Close()
	if errors.Is(err, pgx.ErrNoRows) {
		orderSummary = nil
	} else if err == nil {
		for rows.Next() {
			var status string
			var count int
			if err := rows.Scan(&status, &count); err != nil {
				return nil, err
			} else {
				switch status {
				case "pending":
					orderSummary.Pending = count
					break
				case "accepted":
					orderSummary.Accepted = count
					break
				case "packed":
					orderSummary.Packed = count
					break
				case "ready":
					orderSummary.Ready = count
					break
				case "completed":
					orderSummary.Completed = count
					break
				case "cancelled":
					orderSummary.Cancelled = count
					break
				case "rejected":
					orderSummary.Rejected = count
					break
				default:
					break
				}
			}
		}
		if err := rows.Err(); err != nil { // Add this check!
			return nil, err
		}
	}

	return
}

func GetVendorOrderDetails(orderID string) (*vendors.OrderDetailsResponse, error) {
	ctx := context.Background()

	// Main query to get order details with delivery partner and customer info
	mainQuery := `
		SELECT 
			o.pickup_pin as delivery_partner_pin,
			dp.name as delivery_partner_name,
			dp.phone_number as delivery_partner_phone,
			o.pack_by_time,
			o.paid_by_time,
			o.delivered_by_time,
			o.id as order_id,
			o.customer_name,
			CONCAT(a.address_line1, ', ', a.address_line2, ', ', a.city, ', ', a.state, ' ', a.postal_code) as customer_address,
			o.order_time as order_created_time,
			o.total_amount as order_total_amount
		FROM orders.orders o
		LEFT JOIN delivery_partners.delivery_partners dp ON o.delivery_partner_id = dp.id
		LEFT JOIN user_profile.addresses a ON o.customer_address_id = a.id
		WHERE o.id = $1`

	var response vendors.OrderDetailsResponse
	var deliveryPartnerPin, deliveryPartnerName, deliveryPartnerPhone *string
	var packByTime, paidByTime, deliveredByTime *time.Time

	err := pgClient.QueryRow(ctx, mainQuery, orderID).Scan(
		&deliveryPartnerPin,
		&deliveryPartnerName,
		&deliveryPartnerPhone,
		&packByTime,
		&paidByTime,
		&deliveredByTime,
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
	response.PaidByTime = paidByTime
	response.DeliveredByTime = deliveredByTime

	// Get current order status (latest status from logs)
	statusQuery := `
		SELECT status 
		FROM orders.order_status_logs 
		WHERE order_id = $1
		ORDER BY changed_at DESC 
		LIMIT 1`

	err = pgClient.QueryRow(ctx, statusQuery, orderID).Scan(&response.OrderStatus)
	if err != nil {
		// If no status logs exist, use order.status as fallback
		fallbackQuery := `SELECT status FROM orders.orders WHERE id = $1::uuid`
		err = pgClient.QueryRow(ctx, fallbackQuery, orderID).Scan(&response.OrderStatus)
		if err != nil {
			return nil, err
		}
	}

	// Get order items
	itemsQuery := `
		SELECT 
			i.id as item_id,
			ii.image_url as item_image_url,
			i.name as item_name,
			oi.quantity as qty_ordered
		FROM orders.order_items oi
		JOIN vendor_items.items i ON oi.item_id = i.id
		JOIN vendor_items.item_images ii ON oi.item_id = ii.item_id
		WHERE oi.order_id = $1`

	rows, err := pgClient.Query(ctx, itemsQuery, orderID)
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

	// Update orders table
	_, err := pgClient.Exec(ctx,
		`UPDATE orders.orders SET status = $1, updated_at = $2 WHERE id = $3`,
		newStatus, updateTime, orderID)

	if err != nil {
		return err
	}

	// Insert into status logs
	_, err = pgClient.Exec(ctx,
		`INSERT INTO orders.order_status_logs (order_id, status, changed_at) 
			 VALUES ($1, $2, $3)
			 ON CONFLICT (order_id) 
			 DO UPDATE SET 
				status = EXCLUDED.status,
				changed_at = EXCLUDED.changed_at;`,
		orderID, newStatus, updateTime)

	return err
}

func GetVendorOrders(vendorID string, page, limit int) (*vendors.OrdersListResponse, error) {
	ctx := context.Background()

	// Calculate offset
	offset := (page - 1) * limit

	// Get total count first
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM orders.orders WHERE account_id = $1::uuid`
	err := pgClient.QueryRow(ctx, countQuery, vendorID).Scan(&totalCount)
	if err != nil {
		return nil, err
	}

	// Get orders with latest status from order_status_logs
	ordersQuery := `
		WITH latest_status AS (
			SELECT 
				order_id, 
				status,
				ROW_NUMBER() OVER (PARTITION BY order_id ORDER BY changed_at DESC) as rn
			FROM orders.order_status_logs
		)
		SELECT 
			o.id::text as order_id,
			o.status as order_status,
			o.customer_name,
			o.order_time as order_time_placed,
			o.total_amount,
			o.pack_by_time,
			o.delivered_by_time,
			o.pack_by_time as pick_by_time
		FROM orders.orders o
		WHERE o.account_id = $1::uuid
		ORDER BY o.order_time DESC
		LIMIT $2 OFFSET $3`

	rows, err := pgClient.Query(ctx, ordersQuery, vendorID, limit, offset)
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
