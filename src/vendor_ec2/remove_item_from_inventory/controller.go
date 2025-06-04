package remove_item_from_inventory

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qvickly/database/postgres"
	"strconv"
)

// RemoveItemFromInventoryHandler godoc
// @Summary Remove Item from Vendor Inventory
// @Description Permanently remove an item from the vendor's inventory. This action cannot be undone. The item will no longer be available for sale.
// @Tags inventory
// @Accept json
// @Produce json
// @Param vendor_id path string true "Vendor ID" format:"uuid" example:"123e4567-e89b-12d3-a456-426614174000"
// @Param item_id path int true "Item ID to remove from inventory" example:"456"
// @Router /vendor/{vendor_id}/inventory/{item_id} [delete]
func RemoveItemFromInventoryHandler(c *gin.Context) {
	vendorID := c.Param("vendor_id")
	itemIDStr := c.Param("item_id")
	itemID, err := strconv.Atoi(itemIDStr)

	if vendorID == "" || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid vendor_id or item_id"})
		return
	}

	err = postgres.DeleteInventoryItem(vendorID, itemID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to remove item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Item removed from inventory"})
}
