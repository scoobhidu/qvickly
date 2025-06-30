package verify_pickup

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"qvickly/database/postgres"
	"qvickly/models/delivery"
	"strconv"
)

// VerifyPickup godoc
// @Summary Verify order pickup from vendor using PIN
// @Description Verify that delivery partner has picked up order from vendor using vendor-provided PIN and update order status
// @Tags Delivery Partners
// @Accept json
// @Produce json
// @Param order_id query string true "Order ID" example("123")
// @Param delivery_boy_id query string true "Delivery Partner UUID" format(uuid) example("de111111-2222-3333-4444-555555555555")
// @Param pin body delivery.VerifyPickupRequest true "Pickup PIN provided by vendor"
// @Success 200 {object} delivery.VerifyPickupResponse "Pickup verified successfully and order status updated"
// @Failure 400 {object} delivery.PickupErrorResponse "Invalid parameters, missing PIN, or wrong PIN"
// @Failure 403 {object} delivery.PickupErrorResponse "Order not assigned to this delivery partner or invalid status transition"
// @Failure 404 {object} delivery.PickupErrorResponse "Order not found"
// @Failure 409 {object} delivery.PickupErrorResponse "Order already picked up or invalid status for pickup"
// @Failure 500 {object} delivery.PickupErrorResponse "Internal server error"
// @Router /delivery/verify_pickup [post]
func VerifyPickup(c *gin.Context) {
	// Get query parameters
	vendorAssignmentId := c.Query("vendor_assignment_id")
	deliveryBoyIDParam := c.Query("delivery_boy_id")

	if vendorAssignmentId == "" || deliveryBoyIDParam == "" {
		c.JSON(http.StatusBadRequest, delivery.PickupErrorResponse{
			Success: false,
			Error:   "missing_parameters",
			Message: "Both order_id and delivery_boy_id parameters are required",
			Code:    400,
		})
		return
	}

	// Validate and parse order ID
	orderID, err := uuid.Parse(vendorAssignmentId)
	if err != nil {
		c.JSON(http.StatusBadRequest, delivery.PickupErrorResponse{
			Success: false,
			Error:   "invalid_order_id",
			Message: "Invalid order ID format",
			Code:    400,
		})
		return
	}

	// Validate UUID format for delivery partner
	deliveryPartnerID, err := uuid.Parse(deliveryBoyIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, delivery.PickupErrorResponse{
			Success: false,
			Error:   "invalid_uuid",
			Message: "Invalid UUID format for delivery partner ID",
			Code:    400,
		})
		return
	}

	// Parse request body for PIN
	var request delivery.VerifyPickupRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, delivery.PickupErrorResponse{
			Success: false,
			Error:   "invalid_request",
			Message: "Invalid request body. PIN is required: " + err.Error(),
			Code:    400,
		})
		return
	}

	// Validate PIN format (4 digits)
	if request.Pin < 1000 || request.Pin > 9999 {
		c.JSON(http.StatusBadRequest, delivery.PickupErrorResponse{
			Success: false,
			Error:   "invalid_pin",
			Message: "PIN must be a 4-digit number",
			Code:    400,
		})
		return
	}

	// Verify pickup
	response, err := postgres.ProcessPickupVerification(orderID, deliveryPartnerID, request.Pin)
	if err != nil {
		// Handle different types of errors
		switch err.Error() {
		case "order_not_found":
			c.JSON(http.StatusNotFound, delivery.PickupErrorResponse{
				Success: false,
				Error:   "order_not_found",
				Message: "Order not found or not assigned to this delivery partner",
				Code:    404,
				OrderID: &orderID,
			})
		case "wrong_pin":
			c.JSON(http.StatusBadRequest, delivery.PickupErrorResponse{
				Success: false,
				Error:   "wrong_pin",
				Message: "Incorrect pickup PIN. Please verify with vendor",
				Code:    400,
				OrderID: &orderID,
			})
		case "invalid_status":
			// Get current status for error response
			currentStatus := postgres.GetCurrentOrderStatus(orderID)
			c.JSON(http.StatusConflict, delivery.PickupErrorResponse{
				Success:       false,
				Error:         "invalid_status",
				Message:       "Order is not ready for pickup. Current status: " + currentStatus,
				Code:          409,
				OrderID:       &orderID,
				CurrentStatus: &currentStatus,
			})
		case "already_picked_up":
			c.JSON(http.StatusConflict, delivery.PickupErrorResponse{
				Success: false,
				Error:   "already_picked_up",
				Message: "Order has already been picked up",
				Code:    409,
				OrderID: &orderID,
			})
		default:
			c.JSON(http.StatusInternalServerError, delivery.PickupErrorResponse{
				Success: false,
				Error:   "internal_error",
				Message: "Failed to verify pickup",
				Code:    500,
				OrderID: &orderID,
			})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}
