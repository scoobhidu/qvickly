package vendors

import "time"

// OrderDetailsResponse represents detailed information about an order
type OrderDetailsResponse struct {
	DeliveryPartnerPin   *string     `json:"delivery_partner_pin" example:"1234" description:"4-digit PIN for delivery partner verification"`
	DeliveryPartnerName  *string     `json:"delivery_partner_name" example:"John Delivery" description:"Full name of the assigned delivery partner"`
	DeliveryPartnerPhone *string     `json:"delivery_partner_phone" example:"+1234567890" description:"Contact phone number of delivery partner"`
	PackByTime           *time.Time  `json:"pack_by_time" example:"2025-05-30T14:30:00Z" description:"Expected time when order should be packed"`
	PaidByTime           *time.Time  `json:"paid_by_time" example:"2025-05-30T13:45:00Z" description:"Time when payment was completed"`
	DeliveredByTime      *time.Time  `json:"delivered_by_time" example:"2025-05-30T15:00:00Z" description:"Time when order was delivered (null if not delivered yet)"`
	OrderID              string      `json:"order_id" example:"123e4567-e89b-12d3-a456-426614174000" description:"Unique identifier for the order"`
	OrderStatus          string      `json:"order_status" example:"packed" enums:"pending,accepted,packed,ready,completed,cancelled,rejected" description:"Current status of the order"`
	CustomerName         string      `json:"customer_name" example:"Jane Smith" description:"Full name of the customer who placed the order"`
	CustomerAddress      string      `json:"customer_address" example:"123 Main St, Apt 4B, New York, NY 10001" description:"Complete delivery address"`
	OrderCreatedTime     time.Time   `json:"order_created_time" example:"2025-05-30T13:30:00Z" description:"Timestamp when the order was initially created"`
	OrderTotalAmount     float64     `json:"order_total_amount" example:"45.99" description:"Total amount for the order including taxes and fees"`
	Items                []OrderItem `json:"items" description:"List of items included in the order"`
}

// OrderItem represents an individual item within an order
type OrderItem struct {
	ItemID       int    `json:"item_id" example:"1" description:"Unique identifier for the item"`
	ItemImageURL string `json:"item_image_url" example:"https://example.com/images/margherita-pizza.jpg" description:"URL to the item's display image"`
	ItemName     string `json:"item_name" example:"Margherita Pizza" description:"Display name of the item"`
	QtyOrdered   int    `json:"qty_ordered" example:"2" minimum:"1" description:"Quantity of this item ordered"`
}
