package delivery

import "time"

// OrdersSummaryResponse represents the delivery partner orders summary
type OrdersSummaryResponse struct {
	Completed    int     `json:"completed"`
	Earnings     float64 `json:"earnings"`
	ActiveOrders int     `json:"active_orders"`
}

// RecentOrderResponse represents a recent order for delivery partner
type RecentOrderResponse struct {
	ID                    int       `json:"id"`
	Status                string    `json:"status"`
	Earnings              float64   `json:"earnings"`
	LastStatusUpdatedTime time.Time `json:"last_status_updated_time"`
	Items                 int       `json:"items"`
}
