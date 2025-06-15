package vendors

import "github.com/google/uuid"

// InventoryItem represents an item in vendor's inventory with stock and pricing details
type InventoryItem struct {
	ID            uuid.UUID `json:"id" example:"1001" description:"Unique inventory record ID"`
	ItemID        int       `json:"item_id" example:"456" description:"Reference to the master item catalog"`
	Name          string    `json:"name" example:"Margherita Pizza" description:"Display name of the item"`
	Description   string    `json:"description" example:"Classic pizza with fresh tomatoes, mozzarella cheese, and basil" description:"Detailed description of the item"`
	CategoryID    int       `json:"category_id" example:"5" description:"ID of the category this item belongs to"`
	CategoryName  string    `json:"category_name" example:"Pizza" description:"Name of the category this item belongs to"`
	ImageURL      string    `json:"image_url" example:"https://my-bucket.s3.amazonaws.com/items/margherita-pizza.jpg" description:"URL to the item's display image"`
	StockQuantity int       `json:"stock_quantity" example:"25" minimum:"0" description:"Current available quantity in stock"`
	IsAvailable   bool      `json:"is_available" example:"true" description:"Whether the item is currently available for ordering"`
	Price         float64   `json:"price" example:"12.99" minimum:"0" description:"Current selling price (includes any vendor override)"`
	PriceOverride *float64  `json:"price_override" example:"11.99" minimum:"0" description:"Vendor-specific price override (null if using default price)"`
	OutOfStock    bool      `json:"out_of_stock" example:"false" description:"Calculated field: true if stock_quantity is 0"`
}

// InventorySummary provides overview statistics of vendor's inventory
type InventorySummary struct {
	TotalItems      int `json:"total_items" example:"150" description:"Total number of items in vendor's inventory"`
	InStockItems    int `json:"in_stock_items" example:"135" description:"Number of items currently available (stock > 0 and is_available = true)"`
	OutOfStockItems int `json:"out_of_stock_items" example:"15" description:"Number of items currently out of stock (stock_quantity = 0)"`
}

// AddItemToInventoryRequest for adding new items to vendor inventory
type AddItemToInventoryRequest struct {
	ItemID        int `json:"item_id" binding:"required" example:"456" description:"ID of the item from master catalog to add to inventory"`
	StockQuantity int `json:"stock_quantity" binding:"required" example:"50" minimum:"0" description:"Initial stock quantity to add"`
}

// UpdateInventoryRequest for updating existing inventory items
type UpdateInventoryRequest struct {
	StockQuantity *int     `json:"stock_quantity" example:"30" minimum:"0" description:"Update stock quantity (null to keep current value)"`
	IsAvailable   *bool    `json:"is_available" example:"true" description:"Update availability status (null to keep current value)"`
	PriceOverride *float64 `json:"price_override" example:"10.99" minimum:"0" description:"Update price override (null to remove override and use default price)"`
}

// InventoryMovementRequest for recording stock movements
type InventoryMovementRequest struct {
	ItemID       int    `json:"item_id" binding:"required" example:"456" description:"ID of the inventory item"`
	MovementType string `json:"movement_type" binding:"required" example:"add" enums:"add,remove,adjustment" description:"Type of stock movement"`
	Quantity     int    `json:"quantity" binding:"required" example:"20" minimum:"1" description:"Quantity to add/remove or new quantity for adjustment"`
	Reason       string `json:"reason" example:"Received new shipment" description:"Reason for the stock movement"`
}

// SearchItem represents an item available to be added to inventory
type SearchItem struct {
	ID           int     `json:"id" example:"456" description:"Unique item ID from master catalog"`
	Name         string  `json:"name" example:"Margherita Pizza" description:"Display name of the item"`
	Description  string  `json:"description" example:"Classic pizza with fresh tomatoes, mozzarella cheese, and basil" description:"Detailed description of the item"`
	CategoryID   int     `json:"category_id" example:"5" description:"ID of the category this item belongs to"`
	CategoryName string  `json:"category_name" example:"Pizza" description:"Name of the category this item belongs to"`
	ImageURL     string  `json:"image_url" example:"https://my-bucket.s3.amazonaws.com/items/margherita-pizza.jpg" description:"URL to the item's display image"`
	Price        float64 `json:"price" example:"12.99" minimum:"0" description:"Default price for this item"`
}

// Category represents a product category
type Category struct {
	SuperCategory string `json:"super_category" example:"Margherita Pizza" description:"Super category"`
	ID            int    `json:"id" example:"5" description:"Unique category ID"`
	Name          string `json:"name" example:"Pizza" description:"Category display name"`
}
