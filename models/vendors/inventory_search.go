package vendors

import (
	"github.com/google/uuid"
	"time"
)

// Item represents an item that vendors can add to their inventory
type Item struct {
	ID             int        `json:"id" db:"id"`
	AccountID      uuid.UUID  `json:"account_id" db:"account_id"`
	VendorID       *uuid.UUID `json:"vendor_id" db:"vendor_id"`
	CategoryID     *int       `json:"category_id" db:"category_id"`
	CategoryName   *string    `json:"category_name,omitempty"`
	Name           string     `json:"name" db:"name"`
	Description    *string    `json:"description" db:"description"`
	PriceRetail    *float64   `json:"price_retail" db:"price_retail"`
	PriceWholesale *float64   `json:"price_wholesale" db:"price_wholesale"`
	IsAvailable    bool       `json:"is_available" db:"is_available"`
	Stock          int        `json:"stock" db:"stock"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
	SearchKeywords *string    `json:"search_keywords" db:"search_keywords"`
	IsActive       bool       `json:"is_active" db:"is_active"`
	VendorName     *string    `json:"vendor_name,omitempty"`
}

// SearchFilters represents the available search filters
type SearchFilters struct {
	Query      string `json:"query" example:"Lay's'"`  // Search term for name/description/keywords
	CategoryID int    `json:"category_id" example:"0"` // Filter by category
	Page       int    `json:"page" example:"1"`        // Pagination page number
	Limit      int    `json:"limit" example:"10"`      // Items per page
}

// SearchResponse represents the search API response
type SearchResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    *SearchData  `json:"data,omitempty"`
	Error   *ErrorDetail `json:"error,omitempty"`
}

// SearchData contains the search results and metadata
type SearchData struct {
	Items      []Item            `json:"items"`
	Pagination PaginationInfo    `json:"pagination"`
	Filters    SearchFilters     `json:"applied_filters"`
	Categories []CategorySummary `json:"available_categories"`
}

// PaginationInfo contains pagination details
type PaginationInfo struct {
	CurrentPage  int   `json:"current_page"`
	TotalPages   int   `json:"total_pages"`
	TotalItems   int64 `json:"total_items"`
	ItemsPerPage int   `json:"items_per_page"`
	HasNext      bool  `json:"has_next"`
	HasPrevious  bool  `json:"has_previous"`
}

// CategorySummary provides category info with item counts
type CategorySummary struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ItemCount int    `json:"item_count"`
}

// ErrorDetail represents error information
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
