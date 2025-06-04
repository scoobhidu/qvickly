package inventory_summary

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qvickly/database/postgres"
)

// 1. Get Vendor Inventory Summary

// GetVendorInventorySummaryHandler godoc
// @Summary Get Vendor Inventory Summary
// @Description Retrieve summary statistics of vendor's inventory including total items, in-stock items, and out-of-stock items
// @Tags inventory
// @Accept json
// @Produce json
// @Param vendor_id path string true "Vendor ID" format:"uuid" example:"123e4567-e89b-12d3-a456-426614174000"
// @Success 200 {object} vendors.InventorySummary "Inventory summary retrieved successfully"
// @Router /vendor/{vendor_id}/inventory/summary [get]
func GetVendorInventorySummaryHandler(c *gin.Context) {
	vendorID := c.Param("vendor_id")
	if vendorID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "vendor_id is required"})
		return
	}
	summary, err := postgres.GetInventorySummaryData(vendorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to fetch summary | " + err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": true, "data": summary})
	}
}
