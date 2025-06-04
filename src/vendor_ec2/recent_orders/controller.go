package recent_orders

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qvickly/database/postgres"
	"strconv"
)

// GetVendorOrdersHandler godoc
// @Summary Get Vendor Orders List
// @Description Retrieve a paginated list of orders for a specific vendor. Maximum 10 orders per page to ensure optimal performance.
// @Tags orders
// @Accept json
// @Produce json
// @Param vendor_id query string true "Vendor ID" format:"uuid" example:"123e4567-e89b-12d3-a456-426614174000"
// @Param page query int false "Page number (default: 1)" minimum:"1" default:"1" example:"1"
// @Param limit query int false "Items per page (default: 10, max: 10)" minimum:"1" maximum:"10" default:"10" example:"5"
// @Success 200 {object} vendors.OrdersListResponse "Orders retrieved successfully"
// @Router /vendor/orders/recent_orders [get]
func GetVendorOrdersHandler(c *gin.Context) {
	// Get vendor_id from query parameter
	vendorID := c.Query("vendor_id")
	if vendorID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "vendor_id is required",
		})
		return
	}

	// Get page parameter (default: 1)
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid page number",
		})
		return
	}

	// Get limit parameter (default: 10, max: 10)
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}
	if limit > 10 {
		limit = 10 // Enforce max 10 orders
	}

	// Get orders
	response, err := postgres.GetVendorOrders(vendorID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch orders",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}
