package delivery_details

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"qvickly/database/postgres"
	"qvickly/models/delivery"
)

// GetDeliveryDetails godoc
// @Summary Get delivery partner orders summary
// @Description Retrieve orders summary including completed orders, earnings, and active orders for a delivery partner
// @Tags Delivery Partners
// @Accept json
// @Produce json
// @Param id query string true "Delivery Partner UUID" format(uuid) example("de111111-2222-3333-4444-555555555555")
// @Param detailed query boolean false "Include detailed order lists" example(false)
// @Success 200 {object} delivery.OrdersSummaryResponse "Basic orders summary"
// @Failure 400 {object} delivery.ErrorResponse "Invalid UUID format or missing ID parameter"
// @Failure 404 {object} delivery.ErrorResponse "Delivery partner not found"
// @Failure 500 {object} delivery.ErrorResponse "Internal server error"
// @Router /delivery/profile/orders/summary [get]
func GetDeliveryDetails(c *gin.Context) {
	// Get ID from query parameter
	idParam := c.Query("delivery_id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, delivery.ErrorResponse{
			Error:   "missing_parameter",
			Message: "ID parameter is required",
			Code:    400,
		})
		return
	}

	deliveryID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, delivery.ErrorResponse{
			Error:   "invalid_uuid",
			Message: "Invalid UUID format for ID parameter",
			Code:    400,
		})
		return
	}

	response, err := postgres.GetDeliveryDetails(deliveryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, delivery.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to retrieve delivery details",
			Code:    500,
		})
		return
	}
	c.JSON(http.StatusOK, response)
}

// GetDeliveryVendorItems godoc
// @Summary Get delivery partner orders summary
// @Description Retrieve orders summary including completed orders, earnings, and active orders for a delivery partner
// @Tags Delivery Partners
// @Accept json
// @Produce json
// @Param id query string true "Delivery Partner UUID" format(uuid) example("de111111-2222-3333-4444-555555555555")
// @Param detailed query boolean false "Include detailed order lists" example(false)
// @Success 200 {object} delivery.OrdersSummaryResponse "Basic orders summary"
// @Failure 400 {object} delivery.ErrorResponse "Invalid UUID format or missing ID parameter"
// @Failure 404 {object} delivery.ErrorResponse "Delivery partner not found"
// @Failure 500 {object} delivery.ErrorResponse "Internal server error"
// @Router /delivery/profile/orders/summary [get]
func GetDeliveryVendorItems(c *gin.Context) {
	// Get ID from query parameter
	idParam := c.Query("vendor_assignment_id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, delivery.ErrorResponse{
			Error:   "missing_parameter",
			Message: "ID parameter is required",
			Code:    400,
		})
		return
	}

	vendorAssignmentId, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, delivery.ErrorResponse{
			Error:   "invalid_uuid",
			Message: "Invalid UUID format for ID parameter",
			Code:    400,
		})
		return
	}

	response, err := postgres.GetDeliveryVendorItems(vendorAssignmentId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, delivery.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to retrieve delivery details",
			Code:    500,
		})
		return
	}
	c.JSON(http.StatusOK, response)
}

// GetDeliveryCustomerItems godoc
// @Summary Get delivery partner orders summary
// @Description Retrieve orders summary including completed orders, earnings, and active orders for a delivery partner
// @Tags Delivery Partners
// @Accept json
// @Produce json
// @Param id query string true "Delivery Partner UUID" format(uuid) example("de111111-2222-3333-4444-555555555555")
// @Param detailed query boolean false "Include detailed order lists" example(false)
// @Success 200 {object} delivery.OrdersSummaryResponse "Basic orders summary"
// @Failure 400 {object} delivery.ErrorResponse "Invalid UUID format or missing ID parameter"
// @Failure 404 {object} delivery.ErrorResponse "Delivery partner not found"
// @Failure 500 {object} delivery.ErrorResponse "Internal server error"
// @Router /delivery/profile/orders/summary [get]
func GetDeliveryCustomerItems(c *gin.Context) {
	// Get ID from query parameter
	idParam := c.Query("order_id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, delivery.ErrorResponse{
			Error:   "missing_parameter",
			Message: "ID parameter is required",
			Code:    400,
		})
		return
	}

	orderId, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, delivery.ErrorResponse{
			Error:   "invalid_uuid",
			Message: "Invalid UUID format for ID parameter",
			Code:    400,
		})
		return
	}

	response, err := postgres.GetDeliveryCustomerItems(orderId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, delivery.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to retrieve delivery details" + err.Error(),
			Code:    500,
		})
		return
	}
	c.JSON(http.StatusOK, response)
}
