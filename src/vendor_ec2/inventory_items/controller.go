package get_vendor_items

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qvickly/database/postgres"
	"strconv"
)

// 2. Get Vendor Inventory Items (with pagination and filtering)

// GetVendorInventoryHandler godoc
// @Summary Get Vendor Inventory Items
// @Description Retrieve paginated inventory items for a specific vendor with filtering and search capabilities
// @Tags inventory
// @Accept json
// @Produce json
// @Param vendor_id path string true "Vendor ID" format:"uuid" example:"123e4567-e89b-12d3-a456-426614174000"
// @Param page query int false "Page number (default: 1)" minimum:"1" default:"1" example:"1"
// @Param limit query int false "Items per page (default: 20, max: 50)" minimum:"1" maximum:"50" default:"20" example:"20"
// @Param category_id query string false "Filter by category ID" example:"5"
// @Param search query string false "Search items by name or description" example:"pizza"
// @Param filter query string false "Filter items by stock status" enums:"all,in_stock,out_of_stock" default:"all" example:"in_stock"
// @Router /vendor/{vendor_id}/inventory [get]
func GetVendorInventoryHandler(c *gin.Context) {
	vendorID := c.Param("vendor_id")
	if vendorID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "vendor_id is required"})
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	categoryID := c.Query("category_id")
	search := c.Query("search")
	filter := c.Query("filter") // "all", "in_stock", "out_of_stock"

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 20
	}
	offset := (page - 1) * limit

	totalCount, items, err := postgres.GetInventoryItemsPagination(vendorID, categoryID, search, filter, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to fetch items | " + err.Error()})
	}
	totalPages := (totalCount + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    items,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total_count": totalCount,
			"total_pages": totalPages,
			"has_next":    page < totalPages,
			"has_prev":    page > 1,
		},
	})
}
