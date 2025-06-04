package update_inventory

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qvickly/database/postgres"
	"qvickly/models/vendors"
	"strconv"
)

// 5. Update Inventory Item

// UpdateInventoryItemHandler godoc
// @Summary Update Inventory Item
// @Description Update stock quantity, availability status, or price override for an existing inventory item. All fields are optional - only provide fields you want to update.
// @Tags inventory
// @Accept json
// @Produce json
// @Param vendor_id path string true "Vendor ID" format:"uuid" example:"123e4567-e89b-12d3-a456-426614174000"
// @Param item_id path int true "Item ID to update" example:"456"
// @Param request body vendors.UpdateInventoryRequest true "Inventory update request (all fields optional)"
// @Router /vendor/{vendor_id}/inventory/{item_id} [put]
func UpdateInventoryItemHandler(c *gin.Context) {
	vendorID := c.Param("vendor_id")
	itemIDStr := c.Param("item_id")
	itemID, err := strconv.Atoi(itemIDStr)

	if vendorID == "" || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid vendor_id or item_id"})
		return
	}

	var req vendors.UpdateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid JSON"})
		return
	}

	if err := postgres.UpdateInventoryItem(vendorID, itemID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to update inventory | " + err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Inventory updated successfully"})
	}
}
