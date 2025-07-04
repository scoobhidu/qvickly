package delivery

import (
	"github.com/google/uuid"
	"time"
)

// OrdersSummaryResponse represents the delivery partner orders summary
type OrdersSummaryResponse struct {
	Completed    int     `json:"completed"`
	Earnings     float64 `json:"earnings"`
	ActiveOrders int     `json:"active_orders"`
}

type PickupDetail struct {
	VendorAssignmentId uuid.UUID `json:"vendor_assignment_id"`
	VendorId           uuid.UUID `json:"vendor_id"`
	OrderId            uuid.UUID `json:"order_id"`
	Status             string    `json:"status"`
	PickupTime         time.Time `json:"pickup_time"`
	Address            string    `json:"address"`
	Name               string    `json:"name"`
	Items              int       `json:"items"`
	Amount             string    `json:"amount"`
	PickedUp           bool      `json:"picked_up"`
	OrderTime          time.Time `json:"order_time"`
}

type DeliveryDetail struct {
	OrderId       uuid.UUID `json:"order_id"`
	PickUpTime    time.Time `json:"pickup_time"`
	DeliverByTime time.Time `json:"deliver_by_time"`
	Status        string    `json:"status"`
	Address       string    `json:"address"`
	Name          string    `json:"name"`
}

// RecentOrderResponse represents a recent order for delivery partner
type RecentOrderResponse struct {
	ID                    uuid.UUID `json:"id"`
	Status                string    `json:"status"`
	Earnings              float64   `json:"earnings"`
	LastStatusUpdatedTime time.Time `json:"last_status_updated_time"`
	Items                 int       `json:"items"`
}

// OrderDetailResponse represents the detailed order information for delivery partner
type OrderDetailResponse struct {
	ID                  string      `json:"id"`
	Status              string      `json:"status"`
	AcceptedAt          *time.Time  `json:"accepted_at"`
	DeliveredAt         *time.Time  `json:"delivered_at"`
	DeliveryFee         float64     `json:"delivery_fee"`
	Bonus               float64     `json:"bonus"`
	Earning             float64     `json:"earning"`
	CustomerName        string      `json:"customer_name"`
	DeliveryAddress     string      `json:"delivery_address"`
	DeliveryLatitude    *float64    `json:"delivery_latitude"`
	DeliveryLongitude   *float64    `json:"delivery_longitude"`
	DeliveryInstruction *string     `json:"delivery_instruction"`
	ItemsValue          float64     `json:"items_value"`
	StoreName           string      `json:"store_name"`
	Items               []OrderItem `json:"items"`
}

// OrderItemDetail represents individual item in the order
type OrderItemDetail struct {
	ID       int     `json:"id"`
	Label    string  `json:"label"`
	Qty      int     `json:"qty"`
	ImageURL *string `json:"image_url"`
}

// Extended version with more grocery-specific fields
type OrderItem struct {
	// Basic item information
	ID          uuid.UUID `json:"id" db:"item_id"`
	Name        string    `json:"name" db:"title"`
	Description string    `json:"description,omitempty" db:"description"`

	// Quantity and pricing
	Quantity       int     `json:"quantity" db:"qty"`
	UnitPrice      float64 `json:"unit_price" db:"price_retail"`
	WholesalePrice float64 `json:"wholesale_price,omitempty" db:"price_wholesale"`
	TotalPrice     float64 `json:"total_price" db:"total_price"`

	// Images (multiple image support)
	ImageURL1 string `json:"image_url_1,omitempty" db:"image_url_1"`
	ImageURL2 string `json:"image_url_2,omitempty" db:"image_url_2"`
	ImageURL3 string `json:"image_url_3,omitempty" db:"image_url_3"`
	ImageURL4 string `json:"image_url_4,omitempty" db:"image_url_4"`

	// Category information
	CategoryID    int `json:"category_id" db:"category_id"`
	SubcategoryID int `json:"subcategory_id" db:"subcategory_id"`

	// Additional metadata
	SearchKeywords string `json:"search_keywords,omitempty" db:"search_keywords"`

	// Timestamps
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Extended version with more grocery-specific fields
type OrderItemSummary struct {
	ID        int     `json:"id" db:"item_id"`
	Name      string  `json:"name" db:"title"`
	Quantity  int     `json:"quantity" db:"qty"`
	Price     float64 `json:"total_price" db:"total_price"`
	ImageURL1 string  `json:"image_url_1,omitempty" db:"image_url_1"`
}

// Method to get the primary image URL
func (oi *OrderItem) GetPrimaryImageURL() string {
	if oi.ImageURL1 != "" {
		return oi.ImageURL1
	}
	if oi.ImageURL2 != "" {
		return oi.ImageURL2
	}
	if oi.ImageURL3 != "" {
		return oi.ImageURL3
	}
	if oi.ImageURL4 != "" {
		return oi.ImageURL4
	}
	return ""
}

// Method to get all available image URLs
func (oi *OrderItem) GetAllImageURLs() []string {
	var urls []string
	if oi.ImageURL1 != "" {
		urls = append(urls, oi.ImageURL1)
	}
	if oi.ImageURL2 != "" {
		urls = append(urls, oi.ImageURL2)
	}
	if oi.ImageURL3 != "" {
		urls = append(urls, oi.ImageURL3)
	}
	if oi.ImageURL4 != "" {
		urls = append(urls, oi.ImageURL4)
	}
	return urls
}
