package update_order_status

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"qvickly/database/postgres"
	"strings"
)

// Gin HTTP Handler

// UpdateOrderStatusHandler godoc
// @Summary Update Order Status
// @Description Update the status of an order to track its progress through the fulfillment pipeline. Status changes help coordinate between vendors, delivery partners, and customers.
// @Tags orders
// @Accept json
// @Produce json
// @Param request body UpdateOrderStatusRequest true "Order status update request"
// @Success 200 {object} UpdateOrderStatusResponse "Order status updated successfully"
// @Router /vendor/orders/update_order_status [put]
func UpdateOrderStatusHandler(c *gin.Context) {
	var req UpdateOrderStatusRequest

	// Parse JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid JSON",
		})
		return
	}

	// Validate inputs
	if req.OrderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "order_id is required",
		})
		return
	}

	req.Status = strings.ToLower(strings.TrimSpace(req.Status))
	if req.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "status is required",
		})
		return
	}

	if !ValidOrderStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid status",
		})
		return
	}

	// Update status
	err := postgres.UpdateOrderStatus(req.OrderID, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to update order",
		})
		return
	}

	// Send success response
	c.JSON(http.StatusOK, UpdateOrderStatusResponse{
		Success: true,
		Message: fmt.Sprintf("Order %s updated to %s", req.OrderID, req.Status),
	})
}
