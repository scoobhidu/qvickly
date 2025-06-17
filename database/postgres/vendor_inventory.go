package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"qvickly/models/vendors"
	"strconv"
	"strings"
)

func GetInventorySummaryData(vendorID string) (summary *vendors.InventorySummary, err error) {
	ctx := context.Background()
	query := `SELECT total_items, in_stock_items, out_of_stock_items FROM vendor_inventory_summary WHERE vendor_id = $1::uuid`

	summary = new(vendors.InventorySummary)
	err = pgClient.QueryRow(ctx, query, vendorID).Scan(
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
		FROM public.vendor_inventory vi
		JOIN vendor_items.items i ON vi.item_id = i.id
		JOIN vendor_items.categories c ON i.category_id = c.id
		WHERE vi.vendor_id = $1::uuid AND i.is_active = true`

	args := []interface{}{vendorID}
	argIndex := 2

	if categoryID != "" {
		baseQuery += ` AND i.category_id = $` + strconv.Itoa(argIndex) + `::integer`
		args = append(args, categoryID)
		argIndex++
	}

	if search != "" {
		baseQuery += ` AND (i.name ILIKE $` + strconv.Itoa(argIndex) + ` OR i.description ILIKE $` + strconv.Itoa(argIndex) + `)`
		args = append(args, "%"+search+"%")
		argIndex++
	}

	switch filter {
	case "in_stock":
		baseQuery += ` AND vi.stock_quantity > 0 AND vi.is_available = true`
	case "out_of_stock":
		baseQuery += ` AND vi.stock_quantity = 0`
	}

	// Get total count
	countQuery := `SELECT COUNT(*) ` + baseQuery

	err = pgClient.QueryRow(context.Background(), countQuery, args...).Scan(&totalCount)
	if err != nil {
		items = nil
		totalCount = 0
		return
	}

	//COALESCE(ii.image_url, '') as image_url,
	// Get items
	itemsQuery := `
		SELECT 
			vi.id, vi.item_id, i.name, i.description, i.category_id, c.name as category_name,
			vi.stock_quantity, vi.is_available,
			COALESCE(vi.price_override, i.price_retail) as price, i.price_wholesale,
			CASE WHEN vi.stock_quantity = 0 THEN true ELSE false END as out_of_stock
		` + baseQuery + `
		ORDER BY i.name ASC
		LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)

	args = append(args, limit, offset)

	rows, err := pgClient.Query(context.Background(), itemsQuery, args...)
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
			//&item.ImageURL,
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
	tx, err = pgClient.Begin(ctx)
	if err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Transaction failed"})
		return
	}
	defer tx.Rollback(ctx)

	// Insert into vendor_inventory
	var inventoryID string
	insertQuery := `
		INSERT INTO vendor_inventory (vendor_id, item_id, stock_quantity, is_available, created_at, updated_at)
		VALUES ($1::uuid, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id`

	err = tx.QueryRow(ctx, insertQuery, vendorID, req.ItemID, req.StockQuantity, true).Scan(&inventoryID)
	if err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to add item"})
		return
	}

	// Log inventory movement
	if req.StockQuantity > 0 {
		_, err = tx.Exec(ctx, `
			INSERT INTO inventory_movements (vendor_inventory_id, movement_type, quantity_change, previous_quantity, new_quantity, reason, created_at)
			VALUES ($1::uuid, 'add', $2, 0, $2, 'Initial stock', CURRENT_TIMESTAMP)`,
			inventoryID, req.StockQuantity)
		if err != nil {
			//c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to log movement"})
			return
		}
	}

	tx.Commit(ctx)
	return
}

func UpdateInventoryItem(vendorID string, itemID int, req vendors.UpdateInventoryRequest) (err error) {
	ctx := context.Background()

	// Get current inventory record
	var currentStock int
	var inventoryID string
	err = pgClient.QueryRow(ctx, `
		SELECT id, stock_quantity FROM vendor_inventory 
		WHERE vendor_id = $1::uuid AND item_id = $2`,
		vendorID, itemID).Scan(&inventoryID, &currentStock)

	if err != nil {
		//c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Item not found in inventory"})
		return
	}

	// Build update query dynamically
	updateParts := []string{}
	args := []interface{}{vendorID, itemID}
	argIndex := 3

	if req.StockQuantity != nil {
		updateParts = append(updateParts, "stock_quantity = $"+strconv.Itoa(argIndex))
		args = append(args, *req.StockQuantity)
		argIndex++
	}
	if req.IsAvailable != nil {
		updateParts = append(updateParts, "is_available = $"+strconv.Itoa(argIndex))
		args = append(args, *req.IsAvailable)
		argIndex++
	}
	if req.PriceOverride != nil {
		updateParts = append(updateParts, "price_override = $"+strconv.Itoa(argIndex))
		args = append(args, *req.PriceOverride)
		argIndex++
	}

	if len(updateParts) == 0 {
		//c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "No fields to update"})
		return
	}

	updateParts = append(updateParts, "updated_at = CURRENT_TIMESTAMP")
	updateQuery := `UPDATE vendor_inventory SET ` + strings.Join(updateParts, ", ") +
		` WHERE vendor_id = $1::uuid AND item_id = $2`

	_, err = pgClient.Exec(ctx, updateQuery, args...)
	if err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to update inventory"})
		return
	}

	// Log stock movement if quantity changed
	if req.StockQuantity != nil && *req.StockQuantity != currentStock {
		movementType := "adjustment"
		quantityChange := *req.StockQuantity - currentStock

		_, err = pgClient.Exec(ctx, `
			INSERT INTO inventory_movements (vendor_inventory_id, movement_type, quantity_change, previous_quantity, new_quantity, reason, created_at)
			VALUES ($1::uuid, $2, $3, $4, $5, 'Manual adjustment', CURRENT_TIMESTAMP)`,
			inventoryID, movementType, quantityChange, currentStock, *req.StockQuantity)
	}

	return
}

func DeleteInventoryItem(vendorID string, itemID int) (err error) {
	ctx := context.Background()
	_, err = pgClient.Exec(ctx, `
		DELETE FROM vendor_inventory 
		WHERE vendor_id = $1::uuid AND item_id = $2`,
		vendorID, itemID)
	return
}

func GetItemCategories() (categories []vendors.Category, err error) {
	ctx := context.Background()
	rows, err := pgClient.Query(ctx, `select sc.category_name, c.name, c.id from vendor_items.super_categories sc right join vendor_items.categories c on sc.sub_categories = c.name`)
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
