package order_details

import (
	"github.com/gin-gonic/gin"
	"qvickly/database/postgres"
	"strconv"
)

// GetVendorOrderDetail godoc
// @Summary Get Vendor Order Details
// @Description Retrieve comprehensive details about a specific order including delivery information, customer details, and all items
// @Tags orders
// @Accept json
// @Produce json
// @Param order_id query string true "Order ID" format:"uuid" example:"123e4567-e89b-12d3-a456-426614174000"
// @Success 200 {object} vendors.OrderDetailsResponse "Order details retrieved successfully"
// @Router /vendor/orders/order_details [get]
func GetVendorOrderDetail(ctx *gin.Context) {
	//vars := mux.Vars(r)
	orderID := ctx.Query("order_id")

	if orderID == "" {
		ctx.JSON(500, gin.H{"error": "Order ID not provided"})
		return
	}

	id, err := strconv.Atoi(orderID)
	orderDetails, err := postgres.GetVendorOrderDetails(id)
	if err != nil {
		if err.Error() == "order not found" {
			ctx.JSON(500, gin.H{"error": "order not found"})
			return
		}
		ctx.JSON(500, gin.H{"error": "Internal server error" + err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"order_details": orderDetails})
}
