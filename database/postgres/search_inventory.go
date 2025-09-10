package postgres

import (
	"context"
	"fmt"
	"log"
	"qvickly/models/vendors"
	"strings"
)

func ExecuteItemSearch(filters vendors.SearchFilters) ([]vendors.Item, error) {
	offset := (filters.Page - 1) * filters.Limit

	// Escape single quotes in search query to prevent SQL injection
	filters.Query = strings.ReplaceAll(filters.Query, "'", "''")

	// Base query with proper parameterization to prevent SQL injection
	baseQuery := `
        SELECT 
            gi.item_id,
            gi.category_id,
            gi.subcategory_id,
            gi.title as name,
            gi.description,
            gi.price_retail,
            gi.price_wholesale,
            gi.search_keywords,
            gi.created_at,
            gi.image_url_1,
            gi.image_url_2,
            gi.image_url_3,
            gi.mrp,
            gc.title as category_name,
            gs.title as subcategory_name
        FROM quickkart.master.grocery_items gi 
        LEFT JOIN quickkart.master.grocery_categories gc ON gi.category_id = gc.grocery_category_id
        LEFT JOIN quickkart.master.grocery_subcategories gs ON gi.subcategory_id = gs.grocery_subcategory_id
        WHERE gi.title ILIKE $1`

	var args []interface{}
	args = append(args, "%"+filters.Query+"%")
	argCount := 2

	// Add category filter if specified
	if filters.CategoryID != 0 {
		baseQuery += fmt.Sprintf(" AND gi.category_id = $%d", argCount)
		args = append(args, filters.CategoryID)
		argCount++
	}

	// Add subcategory filter if specified (assuming you might want this)
	if filters.CategoryID != 0 {
		baseQuery += fmt.Sprintf(" AND gi.subcategory_id = $%d", argCount)
		args = append(args, filters.CategoryID)
		argCount++
	}

	// Add ordering and pagination
	baseQuery += fmt.Sprintf(" ORDER BY gi.created_at DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, filters.Limit, offset)

	rows, err := pgPool.Query(context.Background(), baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("search query failed: %v", err)
	}
	defer rows.Close()

	var items []vendors.Item
	for rows.Next() {
		var item vendors.Item
		var subcategoryName *string

		err := rows.Scan(
			&item.ID,
			&item.CategoryID,
			&item.CategoryID,
			&item.Name,
			&item.Description,
			&item.PriceRetail,
			&item.PriceWholesale,
			&item.SearchKeywords,
			&item.CreatedAt,
			&item.ImageURL1,
			&item.ImageURL2,
			&item.ImageURL3,
			&item.Mrp,
			&item.CategoryName,
			&subcategoryName,
		)
		if err != nil {
			log.Printf("Row scan error: %v", err)
			continue
		}

		// Handle nullable subcategory name
		if subcategoryName != nil {
			*item.CategoryName = *subcategoryName
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	return items, nil
}
