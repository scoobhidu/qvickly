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

// OrderDetailResponse represents the detailed order information for delivery partner
type OrderDetailResponse struct {
	ID                  string            `json:"id"`
	Status              string            `json:"status"`
	AcceptedAt          *time.Time        `json:"accepted_at"`
	DeliveredAt         *time.Time        `json:"delivered_at"`
	DeliveryFee         float64           `json:"delivery_fee"`
	Bonus               float64           `json:"bonus"`
	Earning             float64           `json:"earning"`
	CustomerName        string            `json:"customer_name"`
	DeliveryAddress     string            `json:"delivery_address"`
	DeliveryLatitude    *float64          `json:"delivery_latitude"`
	DeliveryLongitude   *float64          `json:"delivery_longitude"`
	DeliveryInstruction *string           `json:"delivery_instruction"`
	ItemsValue          float64           `json:"items_value"`
	StoreName           string            `json:"store_name"`
	Items               []OrderItemDetail `json:"items"`
}

// OrderItemDetail represents individual item in the order
type OrderItemDetail struct {
	ID       int     `json:"id"`
	Label    string  `json:"label"`
	Qty      int     `json:"qty"`
	ImageURL *string `json:"image_url"`
}
