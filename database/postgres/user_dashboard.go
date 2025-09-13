package postgres

import (
	"context"
	"fmt"
	"net/http"
	"qvickly/models/user"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetDashboardItems() []user.CategoryWithItems {
	categoryRows, err := pgPool.Query(context.Background(), "SELECT grocery_category_id, title FROM quickkart.master.grocery_categories ORDER BY title")
	if err != nil {
		return make([]user.CategoryWithItems, 0)
	}
	defer categoryRows.Close()

	var categories []user.CategoryWithItems

	for categoryRows.Next() {
		var category user.CategoryWithItems
		err := categoryRows.Scan(&category.ID, &category.Name)
		if err != nil {
			return make([]user.CategoryWithItems, 0)
		}

		// Get up to 5 items for each category
		itemRows, err := pgPool.Query(context.Background(), `
			SELECT item_id, title, description, price_wholesale, price_retail, COALESCE(image_url_1, '') as image_url_1 
			FROM quickkart.master.grocery_items 
			WHERE category_id = $1 
			ORDER BY title 
			LIMIT 5`, category.ID)
		if err != nil {
			return make([]user.CategoryWithItems, 0)
		}

		var items []user.GroceryItem
		for itemRows.Next() {
			var item user.GroceryItem
			err := itemRows.Scan(&item.ID, &item.Title, &item.Description, &item.PriceWholesale, &item.PriceRetail, &item.ImageURL1)
			if err != nil {
				itemRows.Close()
				return make([]user.CategoryWithItems, 0)
			}
			items = append(items, item)
		}
		itemRows.Close()

		category.Items = items
		categories = append(categories, category)
	}

	return categories
}

func GetCategoriesWithSubcategories() []user.CategoryWithSubcategories {
	categoryRows, err := pgPool.Query(context.Background(), "SELECT grocery_category_id, title FROM quickkart.master.grocery_categories ORDER BY title")
	if err != nil {
		return make([]user.CategoryWithSubcategories, 0)
	}
	defer categoryRows.Close()

	var categories []user.CategoryWithSubcategories

	for categoryRows.Next() {
		var category user.CategoryWithSubcategories
		err := categoryRows.Scan(&category.ID, &category.Name)
		if err != nil {
			return make([]user.CategoryWithSubcategories, 0)
		}

		// Get subcategories for each category
		// Note: You'll need to add an image_url field to your subcategories table
		// For now, I'm using a placeholder
		subcategoryRows, err := pgPool.Query(context.Background(), `
			SELECT grocery_subcategory_id, grocery_category_id, title, 
			       'https://example.com/subcategory_placeholder.jpg' as image_url
			FROM quickkart.master.grocery_subcategories 
			WHERE grocery_category_id = $1 
			ORDER BY title`, category.ID)
		if err != nil {
			return make([]user.CategoryWithSubcategories, 0)
		}

		var subcategories []user.GrocerySubcategory
		for subcategoryRows.Next() {
			var subcategory user.GrocerySubcategory
			err := subcategoryRows.Scan(&subcategory.ID, &subcategory.CategoryID, &subcategory.Title, &subcategory.ImageURL)
			if err != nil {
				subcategoryRows.Close()
				return make([]user.CategoryWithSubcategories, 0)
			}
			subcategories = append(subcategories, subcategory)
		}
		subcategoryRows.Close()

		category.Subcategories = subcategories
		categories = append(categories, category)
	}

	return categories
}

func GetNudges() []user.DashboardNudge {
	rows, err := pgPool.Query(context.Background(), "SELECT nudge_title, nudge_body, nudge_image, nudge_navigation_uri FROM quickkart.master.dashboard_nudges")
	if err != nil {
		return make([]user.DashboardNudge, 0)
	}
	defer rows.Close()

	var nudges []user.DashboardNudge
	for rows.Next() {
		var nudge user.DashboardNudge
		err := rows.Scan(&nudge.NudgeTitle, &nudge.NudgeBody, &nudge.NudgeImage, &nudge.NudgeNavigationURI)
		if err != nil {
			return make([]user.DashboardNudge, 0)
		}
		nudges = append(nudges, nudge)
	}
	return nudges
}

func GetItemsBySubCategory(c *gin.Context, err error, categoryID int, limit int, offset int) ([]user.GroceryItem, bool) {
	rows, err := pgPool.Query(context.Background(), `
		SELECT item_id, title, description, price_wholesale, price_retail, mrp, COALESCE(image_url_1, '') as image_url_1 
		FROM quickkart.master.grocery_items 
		WHERE subcategory_id = $1 
		ORDER BY title 
		LIMIT $2 OFFSET $3`, categoryID, limit, offset)
	if err != nil {
		return make([]user.GroceryItem, 0), true
	}
	defer rows.Close()

	var items []user.GroceryItem
	for rows.Next() {
		var item user.GroceryItem
		err := rows.Scan(&item.ID, &item.Title, &item.Description, &item.PriceWholesale, &item.PriceRetail, &item.Mrp, &item.ImageURL1)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan item"})
			return make([]user.GroceryItem, 0), true
		}
		items = append(items, item)
	}
	return items, false
}

func GetDailyEssentialItems(c *gin.Context) ([]user.GroceryItem, bool) {
	rows, err := pgPool.Query(context.Background(), `
		SELECT item_id, title, description, price_wholesale, price_retail, mrp, COALESCE(image_url_1, '') as image_url_1 
		FROM quickkart.master.grocery_items 
		WHERE daily_essential = true 
		ORDER BY title 
		LIMIT 30`)
	if err != nil {
		return make([]user.GroceryItem, 0), true
	}
	defer rows.Close()

	var items []user.GroceryItem
	for rows.Next() {
		var item user.GroceryItem
		err := rows.Scan(&item.ID, &item.Title, &item.Description, &item.PriceWholesale, &item.PriceRetail, &item.Mrp, &item.ImageURL1)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan item"})
			return make([]user.GroceryItem, 0), true
		}
		items = append(items, item)
	}
	return items, false
}

func GetHotItems(c *gin.Context) ([]user.GroceryItem, bool) {
	rows, err := pgPool.Query(context.Background(), `
		SELECT item_id, title, description, price_wholesale, price_retail, mrp, COALESCE(image_url_1, '') as image_url_1 
		FROM quickkart.master.grocery_items 
		WHERE hot_products = true 
		ORDER BY title 
		LIMIT 30`)
	if err != nil {
		return make([]user.GroceryItem, 0), true
	}
	defer rows.Close()

	var items []user.GroceryItem
	for rows.Next() {
		var item user.GroceryItem
		err := rows.Scan(&item.ID, &item.Title, &item.Description, &item.PriceWholesale, &item.PriceRetail, &item.Mrp, &item.ImageURL1)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan item"})
			return make([]user.GroceryItem, 0), true
		}
		items = append(items, item)
	}
	return items, false
}

func GetRecentSearches(c *gin.Context, cID string) ([]user.RecentSearch, bool) {
	rows, err := pgPool.Query(context.Background(), `
		SELECT search_term, created_at
		FROM quickkart.profile.user_recent_searches
		where customer_id = $1::uuid
		LIMIT 3`, cID)
	if err != nil {
		return make([]user.RecentSearch, 0), true
	}
	defer rows.Close()

	var items []user.RecentSearch
	for rows.Next() {
		var item user.RecentSearch
		err := rows.Scan(&item.Search, &item.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan item"})
			return make([]user.RecentSearch, 0), true
		}
		items = append(items, item)
	}
	return items, false
}

func GetItemsByFilter(categoryID int, minPrice, maxPrice int, searchQuery string, limit int, offset int) ([]user.GroceryItem, bool) {
	// Build dynamic query based on filters
	query := `SELECT item_id, title, description, price_wholesale, price_retail, mrp, COALESCE(image_url_1, '') as image_url_1 
             FROM quickkart.master.grocery_items 
             WHERE subcategory_id = $1`

	args := []interface{}{categoryID}
	argIndex := 2

	// Add price range filters
	if minPrice != -1 {
		query += fmt.Sprintf(" AND price_retail >= $%d", argIndex)
		args = append(args, minPrice)
		argIndex++
	}

	if maxPrice != -1 {
		query += fmt.Sprintf(" AND price_retail <= $%d", argIndex)
		args = append(args, maxPrice)
		argIndex++
	}

	// Add search query filter
	if searchQuery != "" {
		query += fmt.Sprintf(` AND (
           LOWER(title) ILIKE $%d OR 
           LOWER(description) ILIKE $%d OR 
           LOWER(search_keywords) ILIKE $%d
       )`, argIndex, argIndex, argIndex)
		searchTerm := "%" + strings.ToLower(searchQuery) + "%"
		args = append(args, searchTerm)
		argIndex++
	}

	// Add ordering and pagination
	query += fmt.Sprintf(" ORDER BY title LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	// Execute query
	rows, err := pgPool.Query(context.Background(), query, args...)
	if err != nil {
		return make([]user.GroceryItem, 0), true
	}
	defer rows.Close()

	var items []user.GroceryItem
	for rows.Next() {
		var item user.GroceryItem
		err := rows.Scan(&item.ID, &item.Title, &item.Description, &item.PriceWholesale, &item.PriceRetail, &item.Mrp, &item.ImageURL1)
		if err != nil {
			return make([]user.GroceryItem, 0), true
		}
		items = append(items, item)
	}
	return items, false
}
