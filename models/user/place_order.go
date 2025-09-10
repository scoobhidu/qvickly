package user

import (
	"github.com/google/uuid"
)

// Request/Response structs
type OrderItem struct {
	ItemID   int64 `json:"item_id" binding:"required"`
	Quantity int   `json:"quantity" binding:"required,min=1"`
}

type PlaceOrderRequest struct {
	CustomerID uuid.UUID   `json:"customer_id" binding:"required"`
	AddressID  int64       `json:"address_id" binding:"required"`
	Items      []OrderItem `json:"items" binding:"required,min=1"`
}

type PlaceOrderResponse struct {
	OrderID          uuid.UUID `json:"order_id"`
	AssignedVendorID uuid.UUID `json:"assigned_vendor_id"`
	DeliveryBoyID    uuid.UUID `json:"delivery_boy_id"`
	EstimatedTime    string    `json:"estimated_delivery_time"`
	TotalAmount      float64   `json:"total_amount"`
	Message          string    `json:"message"`
}

// Database models
type Vendor struct {
	ID        uuid.UUID `db:"vendor_id"`
	Name      string    `db:"business_name"`
	Latitude  float64   `db:"latitude"`
	Longitude float64   `db:"longitude"`
	IsActive  bool      `db:"is_active"`
	IsLive    bool      `db:"is_live"`
}

type DeliveryBoy struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"full_name"`
	Latitude  float64   `db:"latitude"`
	Longitude float64   `db:"longitude"`
	IsActive  bool      `db:"is_active"`
}

type CustomerAddress struct {
	ID        int64   `db:"address_id"`
	Latitude  float64 `db:"latitude"`
	Longitude float64 `db:"longitude"`
}

type InventoryItem struct {
	VendorID          uuid.UUID `db:"vendor_id"`
	ItemID            int64     `db:"item_id"`
	Quantity          int       `db:"qty"`
	WholesalePrice    float64   `db:"wholesale_price_override"`
	RetailPrice       float64   `db:"price_retail"`
	WholesaleFallback float64   `db:"price_wholesale"`
}
