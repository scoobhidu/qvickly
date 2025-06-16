package search_items_for_inventory_addition

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"qvickly/database/postgres"
	"qvickly/models/vendors"
)

// SearchInventoryItems searchItems godoc
// @Summary Search items with advanced filters
// @Description Search for items that vendors can add to their inventory with various filtering options
// @Tags Items
// @Accept json
// @Produce json
// @Param request body vendors.SearchFilters true "Search filters, if 0 then that means all category IDs"
// @Success 200 {object} vendors.SearchResponse "Successful search results"
// @Failure 400 {object} vendors.SearchResponse "Invalid query parameters"
// @Failure 500 {object} vendors.SearchResponse "Internal server error"
// @Router /vendor/inventory/search [post]
func SearchInventoryItems(c *gin.Context) {
	var filters vendors.SearchFilters

	// Bind query parameters
	if err := c.ShouldBindJSON(&filters); err != nil {
		c.JSON(http.StatusBadRequest, vendors.SearchResponse{
			Success: false,
			Message: "Invalid query parameters",
			Error: &vendors.ErrorDetail{
				Code:    "INVALID_PARAMETERS",
				Message: err.Error(),
			},
		})
		return
	}

	// Set defaults
	if filters.Page <= 0 {
		filters.Page = 1
	}
	if filters.Limit <= 0 {
		filters.Limit = 20
	}
	if filters.Limit > 100 {
		filters.Limit = 100 // Max limit
	}

	// Build the search query
	items, err := postgres.ExecuteItemSearch(filters)
	if err != nil {
		log.Printf("Search error: %v", err)
		c.JSON(http.StatusInternalServerError, vendors.SearchResponse{
			Success: false,
			Message: "Search failed",
			Error: &vendors.ErrorDetail{
				Code:    "SEARCH_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	// Get available categories with counts
	//categories, err := getAvailableCategories(filters)
	//if err != nil {
	//	log.Printf("Categories error: %v", err)
	//	// Continue without categories if there's an error
	//	categories = []CategorySummary{}
	//}

	// Calculate pagination
	pagination := vendors.PaginationInfo{
		CurrentPage:  filters.Page,
		ItemsPerPage: filters.Limit,
		HasPrevious:  filters.Page > 1,
	}

	c.JSON(http.StatusOK, vendors.SearchResponse{
		Success: true,
		Message: "Items retrieved successfully",
		Data: &vendors.SearchData{
			Items:      items,
			Pagination: pagination,
			Filters:    filters,
			//Categories: categories,
		},
	})
}

//// Get available categories with item counts
//func getAvailableCategories(filters SearchFilters) ([]CategorySummary, error) {
//	query := `
//		SELECT
//			c.id,
//			c.name,
//			COUNT(i.id) as item_count
//		FROM vendor_items.categories c
//		LEFT JOIN vendor_items.items i ON c.id = i.category_id AND i.is_active = true
//	`
//
//	// Apply same filters for category counts (except category_id filter)
//	conditions := []string{}
//	args := []interface{}{}
//	argCount := 1
//
//	// Text search filter for categories
//	if filters.Query != "" {
//		conditions = append(conditions, fmt.Sprintf(`
//			(to_tsvector('english', COALESCE(i.name, '') || ' ' || COALESCE(i.description, '') || ' ' || COALESCE(i.search_keywords, ''))
//			@@ plainto_tsquery('english', $%d)
//			OR i.name ILIKE $%d
//			OR i.description ILIKE $%d)
//		`, argCount, argCount+1, argCount+2))
//
//		searchTerm := strings.TrimSpace(filters.Query)
//		args = append(args, searchTerm, "%"+searchTerm+"%", "%"+searchTerm+"%")
//		argCount += 3
//	}
//
//	if len(conditions) > 0 {
//		query += " WHERE " + strings.Join(conditions, " AND ")
//	}
//
//	query += " GROUP BY c.id, c.name ORDER BY c.name"
//
//	rows, err := db.Query(query, args...)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//
//	var categories []CategorySummary
//	for rows.Next() {
//		var category CategorySummary
//		err := rows.Scan(&category.ID, &category.Name, &category.ItemCount)
//		if err != nil {
//			log.Printf("Category scan error: %v", err)
//			continue
//		}
//		categories = append(categories, category)
//	}
//
//	return categories, nil
//}
