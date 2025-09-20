package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"qvickly/models/vendors"
	"strconv"
	"strings"
)

func GetInventorySummaryData(vendorID string) (summary *vendors.InventorySummary, err error) {
	ctx := context.Background()
	query := `SELECT total_items, in_stock_items, out_of_stock_items FROM vendor_inventory_summary WHERE vendor_id = $1::uuid`

	summary = new(vendors.InventorySummary)
	err = pgPool.QueryRow(ctx, query, vendorID).Scan(
		&summary.TotalItems, &summary.InStockItems, &summary.OutOfStockItems)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			*summary = vendors.InventorySummary{TotalItems: 0, InStockItems: 0, OutOfStockItems: 0}
		} else {
			*summary = vendors.InventorySummary{TotalItems: 0, InStockItems: 0, OutOfStockItems: 0}
			//c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to fetch summary"})
			return
		}
	}
	return
}

func GetInventoryItemsPagination(vendorID string, categoryID string, search string, filter string, limit int, offset int) (totalCount int, items []vendors.InventoryItem, err error) {
	totalCount = 0

	// Build query with filters
	//	JOIN vendor_items.item_images ii ON vi.item_id = ii.item_id
	baseQuery := `
		FROM quickkart.vendor.inventory vi
		JOIN quickkart.master.grocery_items i ON vi.item_id = i.item_id
		JOIN quickkart.master.grocery_subcategories sc ON i.subcategory_id = sc.grocery_subcategory_id
		WHERE vi.vendor_id = $1::uuid`

	args := []interface{}{vendorID}
	argIndex := 2

	if categoryID != "" {
		baseQuery += ` AND i.grocery_category_id = $` + strconv.Itoa(argIndex) + `::integer`
		args = append(args, categoryID)
		argIndex++
	}

	if search != "" {
		baseQuery += ` AND (i.title ILIKE $` + strconv.Itoa(argIndex) + ` OR i.description ILIKE $` + strconv.Itoa(argIndex) + `)`
		args = append(args, "%"+search+"%")
		argIndex++
	}

	switch filter {
	case "in_stock":
		baseQuery += ` AND vi.qty > 0`
	case "out_of_stock":
		baseQuery += ` AND vi.qty = 0`
	}

	// Get total count
	countQuery := `SELECT COUNT(*) ` + baseQuery

	err = pgPool.QueryRow(context.Background(), countQuery, args...).Scan(&totalCount)
	if err != nil {
		items = nil
		totalCount = 0
		return
	}

	//COALESCE(ii.image_url, '') as image_url,
	// Get items
	itemsQuery := `
		SELECT 
			'', vi.item_id, i.title, i.description, i.subcategory_id, sc.title as category_name,
			i.image_url_1,
			vi.qty, vi.qty > 0 as is_available,
			COALESCE(vi.wholesale_price_override, i.price_retail) as price, i.price_wholesale,
			CASE WHEN vi.qty = 0 THEN true ELSE false END as out_of_stock
		` + baseQuery + `
		ORDER BY i.title ASC
		LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)

	args = append(args, limit, offset)

	rows, err := pgPool.Query(context.Background(), itemsQuery, args...)
	if err != nil {
		totalCount = 0
		items = nil
		return
	}
	defer rows.Close()

	//var items []vendors.InventoryItem
	for rows.Next() {
		var item vendors.InventoryItem
		err = rows.Scan(
			&item.ID,
			&item.ItemID, &item.Name, &item.Description, &item.CategoryID,
			&item.CategoryName,
			&item.ImageURL,
			&item.StockQuantity, &item.IsAvailable,
			&item.Price, &item.PriceOverride, &item.OutOfStock)
		if err != nil {
			totalCount = 0
			items = nil
			return
		}
		items = append(items, item)
	}
	return
}

func AddItemsToInventory(vendorID string, req vendors.AddItemToInventoryRequest) (err error) {
	var tx pgx.Tx
	ctx := context.Background()
	tx, err = pgPool.Begin(ctx)
	if err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Transaction failed"})
		return
	}
	defer tx.Rollback(ctx)

	insertQuery := `
		INSERT INTO quickkart.vendor.inventory (vendor_id, item_id, qty, created_at, updated_at)
		VALUES ($1::uuid, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING vendor_id`

	_, err = tx.Exec(ctx, insertQuery, vendorID, req.ItemID, req.StockQuantity)
	if err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to add item"})
		return
	}

	// Log inventory movement
	//if req.StockQuantity > 0 {
	//	_, err = tx.Exec(ctx, `
	//		INSERT INTO inventory_movements (vendor_inventory_id, movement_type, quantity_change, previous_quantity, new_quantity, reason, created_at)
	//		VALUES ($1::uuid, 'add', $2, 0, $2, 'Initial stock', CURRENT_TIMESTAMP)`,
	//		inventoryID, req.StockQuantity)
	//	if err != nil {
	//		//c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to log movement"})
	//		return
	//	}
	//}

	tx.Commit(ctx)
	return
}
func UpdateInventoryItem(vendorID string, itemID int, req vendors.UpdateInventoryRequest) (err error) {
	ctx := context.Background()

	// Get current inventory record
	var currentQty int
	var inventoryItemID string
	err = pgPool.QueryRow(ctx, `
       SELECT item_id, qty FROM quickkart.vendor.inventory 
       WHERE vendor_id = $1::uuid AND item_id = $2`,
		vendorID, itemID).Scan(&inventoryItemID, &currentQty)

	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("item not found in inventory")
		}
		return err
	}

	// Build update query dynamically
	updateParts := []string{}
	args := []interface{}{vendorID, itemID}
	argIndex := 3

	if req.StockQuantity != nil {
		updateParts = append(updateParts, "qty = $"+strconv.Itoa(argIndex))
		args = append(args, *req.StockQuantity)
		argIndex++
	}

	// Note: Based on your schema, vendor.inventory doesn't have is_available or price_override columns
	// If you need these features, you might need to add these columns to the table or handle them differently

	if req.PriceOverride != nil {
		updateParts = append(updateParts, "wholesale_price_override = $"+strconv.Itoa(argIndex))
		args = append(args, *req.PriceOverride)
		argIndex++
	}

	if len(updateParts) == 0 {
		return fmt.Errorf("no fields to update")
	}

	updateParts = append(updateParts, "updated_at = CURRENT_TIMESTAMP")
	updateQuery := `UPDATE quickkart.vendor.inventory SET ` + strings.Join(updateParts, ", ") +
		` WHERE vendor_id = $1::uuid AND item_id = $2`

	result, err := pgPool.Exec(ctx, updateQuery, args...)
	if err != nil {
		return fmt.Errorf("failed to update inventory: %v", err)
	}

	// Check if any rows were affected
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no inventory record updated")
	}

	// Log stock movement if quantity changed
	if req.StockQuantity != nil && *req.StockQuantity != currentQty {
		quantityChange := *req.StockQuantity - currentQty

		// Since there's no inventory_movements table in your schema, you might want to create one
		// or handle stock movements differently. For now, I'll comment this out:
		/*
		   _, err = pgPool.Exec(ctx, `
		      INSERT INTO inventory_movements (vendor_inventory_id, movement_type, quantity_change, previous_quantity, new_quantity, reason, created_at)
		      VALUES ($1::uuid, $2, $3, $4, $5, 'Manual adjustment', CURRENT_TIMESTAMP)`,
		      inventoryItemID, movementType, quantityChange, currentQty, *req.StockQuantity)

		   if err != nil {
		      // Log the error but don't fail the main operation
		      log.Printf("Failed to log inventory movement: %v", err)
		   }
		*/

		// Alternative: You could log to a general audit table or handle this differently
		log.Printf("Inventory updated for vendor %s, item %d: %d -> %d (change: %d)",
			vendorID, itemID, currentQty, *req.StockQuantity, quantityChange)
	}

	return nil
}

func DeleteInventoryItem(vendorID string, itemID int) (err error) {
	ctx := context.Background()
	_, err = pgPool.Exec(ctx, `
		DELETE FROM quickkart.vendor.inventory 
		WHERE vendor_id = $1::uuid AND item_id = $2`,
		vendorID, itemID)
	return
}

func GetItemCategories() (categories []vendors.Category, err error) {
	ctx := context.Background()
	rows, err := pgPool.Query(ctx, `select s.title, sc.title, sc.subcategory_id from quickkart.master.grocery_items sc left join quickkart.master.grocery_subcategories s on sc.subcategory_id = s.grocery_subcategory_id`)
	if err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to fetch categories"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var cat vendors.Category
		rows.Scan(&cat.SuperCategory, &cat.Name, &cat.ID)
		categories = append(categories, cat)
	}

	return
}
