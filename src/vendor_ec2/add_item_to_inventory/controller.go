package add_item_to_inventory

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qvickly/database/postgres"
	"qvickly/models/vendors"
)

// AddItemToInventoryHandler godoc
// @Summary Add Item to Vendor Inventory
// @Description Add a new item from the master catalog to the vendor's inventory with initial stock quantity and optional price override
// @Tags inventory
// @Accept json
// @Produce json
// @Param vendor_id path string true "Vendor ID" format:"uuid" example:"123e4567-e89b-12d3-a456-426614174000"
// @Param request body vendors.AddItemToInventoryRequest true "Add item to inventory request"
// @Router /api/vendors/{vendor_id}/inventory [post]
func AddItemToInventoryHandler(c *gin.Context) {
	vendorID := c.Param("vendor_id")
	if vendorID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "vendor_id is required"})
		return
	}

	var req vendors.AddItemToInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid JSON"})
		return
	}

	if req.ItemID <= 0 || req.StockQuantity < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid item_id or stock_quantity"})
		return
	}

	if err := postgres.AddItemsToInventory(vendorID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to add item to inventory"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Item added to inventory"})
	}
	return
}
