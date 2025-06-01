package vendors

import "time"

// OrderListItem represents a single order in the list view
type OrderListItem struct {
	OrderID         string     `json:"order_id" example:"123e4567-e89b-12d3-a456-426614174000" description:"Unique identifier for the order"`
	OrderStatus     string     `json:"order_status" example:"packed" enums:"pending,accepted,packed,ready,completed,cancelled,rejected" description:"Current status of the order"`
	CustomerName    string     `json:"customer_name" example:"John Doe" description:"Full name of the customer who placed the order"`
	OrderTimePlaced time.Time  `json:"order_time_placed" example:"2025-05-30T13:30:00Z" description:"Timestamp when the order was placed"`
	TotalAmount     float64    `json:"total_amount" example:"45.99" description:"Total amount for the order including taxes and delivery fees"`
	PackByTime      *time.Time `json:"pack_by_time" example:"2025-05-30T14:30:00Z" description:"Expected time when order should be packed (null if not set)"`
	DeliveredByTime *time.Time `json:"delivered_by_time" example:"2025-05-30T15:00:00Z" description:"Time when order was delivered (null if not delivered yet)"`
	PickByTime      *time.Time `json:"pick_by_time" example:"2025-05-30T14:45:00Z" description:"Expected time when order should be picked up by delivery partner (null if not set)"`
}

// OrdersListResponse represents paginated list of orders
type OrdersListResponse struct {
	Success    bool            `json:"success" example:"true" description:"Indicates if the request was successful"`
	Data       []OrderListItem `json:"data" description:"Array of orders for the current page"`
	Page       int             `json:"page" example:"1" minimum:"1" description:"Current page number"`
	Limit      int             `json:"limit" example:"10" minimum:"1" maximum:"10" description:"Number of items per page (max 10)"`
	TotalCount int             `json:"total_count" example:"125" description:"Total number of orders across all pages"`
	TotalPages int             `json:"total_pages" example:"13" description:"Total number of pages available"`
	HasNext    bool            `json:"has_next" example:"true" description:"Whether there are more pages after the current page"`
	HasPrev    bool            `json:"has_prev" example:"false" description:"Whether there are pages before the current page"`
}
