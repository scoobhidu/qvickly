package postgres

import (
	"context"
	"github.com/google/uuid"
	"qvickly/models/delivery"
	"strconv"
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

	err := pgClient.QueryRow(context.Background(), completedQuery, partnerID).Scan(&completedCount, &totalEarnings)
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
	err = pgClient.QueryRow(context.Background(), activeQuery, partnerID).Scan(&activeCount)
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

	rows, err := pgClient.Query(context.Background(), baseQuery, args...)
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
