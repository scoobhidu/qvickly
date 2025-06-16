package postgres

import (
	"context"
	"fmt"
	"log"
	"qvickly/models/vendors"
	"strconv"
	"strings"
)

// Execute the item search with filters
func ExecuteItemSearch(filters vendors.SearchFilters) ([]vendors.Item, error) {
	offset := (filters.Page - 1) * filters.Limit

	filters.Query = strings.ReplaceAll(filters.Query, "'", "''")

	// Base query with joins
	baseQuery := "SELECT i.id, i.account_id, i.account_id, i.category_id, " +
		"i.name, i.description, i.price_retail, i.price_wholesale, i.is_available, " +
		"i.stock, i.created_at, i.updated_at, i.search_keywords, i.is_active, " +
		"c.name as category_name " +
		"FROM vendor_items.items i LEFT JOIN vendor_items.categories c ON i.category_id = c.id " +
		"where i.name Like '%" +
		filters.Query + "%'"
	if filters.CategoryID == 0 {
		baseQuery += " LIMIT " + strconv.Itoa(filters.Limit) +
			" OFFSET " + strconv.Itoa(offset)
	} else {
		baseQuery += " AND i.category_id = " + strconv.Itoa(filters.CategoryID) +
			"LIMIT " + strconv.Itoa(filters.Limit) +
			" OFFSET " + strconv.Itoa(offset)
	}

	rows, err := pgClient.Query(context.Background(), baseQuery)
	if err != nil {
		return nil, fmt.Errorf("search query failed: %v", err)
	}
	defer rows.Close()

	var items []vendors.Item
	for rows.Next() {
		var item vendors.Item
		err := rows.Scan(
			&item.ID,
			&item.AccountID,
			&item.VendorID,
			&item.CategoryID,
			&item.Name,
			&item.Description,
			&item.PriceRetail,
			&item.PriceWholesale,
			&item.IsAvailable,
			&item.Stock,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.SearchKeywords,
			&item.IsActive,
			&item.CategoryName,
		)
		if err != nil {
			log.Printf("Row scan error: %v", err)
			continue
		}
		items = append(items, item)
	}

	return items, nil
}
