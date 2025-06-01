package profile_details

import (
	"github.com/gin-gonic/gin"
	"qvickly/database/postgres"
)

// GetVendorOrderSummary godoc
// @Summary Get Vendor Today's Order Summary
// @Description Retrieve order count statistics for the current day, grouped by order status. Provides a quick overview of vendor's daily order performance.
// @Tags order-analytics
// @Accept json
// @Produce json
// @Param vendor_id query string true "Vendor ID" format:"uuid" example:"123e4567-e89b-12d3-a456-426614174000"
// @Success 200 {object} vendors.TodayOrderSummary "Today's order summary retrieved successfully"
// @Router /api/vendor/order-summary [get]
func GetVendorOrderSummary(context *gin.Context) {
	vendorId := context.Query("vendor_id")

	vendorDetails, err := postgres.GetVendorTodaysOrderSummary(vendorId)
	if err != nil {
		context.JSON(500, gin.H{"error": "Failed to get vendor profile"})
		return
	} else {
		context.JSON(200, gin.H{"vendor_details": vendorDetails})
	}
}
