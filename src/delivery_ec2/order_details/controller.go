package order_details

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"qvickly/database/postgres"
	"qvickly/models/delivery"
	"strconv"
)

// GetDeliveryOrderDetail godoc
// @Summary Get detailed order information for delivery partner
// @Description Retrieve comprehensive order details for a delivery partner including items, customer info, and earnings
// @Tags Delivery Partners
// @Accept json
// @Produce json
// @Param order_id query string true "Order ID" example("1")
// @Param delivery_partner_id query string true "Delivery Partner UUID" format(uuid) example("de111111-2222-3333-4444-555555555555")
// @Success 200 {object} delivery.OrderDetailResponse "Detailed order information"
// @Failure 400 {object} delivery.ErrorResponse "Invalid order ID, missing parameters, or invalid UUID format"
// @Failure 403 {object} delivery.ErrorResponse "Order not assigned to this delivery partner"
// @Failure 404 {object} delivery.ErrorResponse "Order not found or delivery partner not found"
// @Failure 500 {object} delivery.ErrorResponse "Internal server error"
// @Router /delivery/order/detail [get]
func GetDeliveryOrderDetail(c *gin.Context) {
	// Get parameters from query
	orderIDParam := c.Query("order_id")
	partnerIDParam := c.Query("delivery_partner_id")

	if orderIDParam == "" || partnerIDParam == "" {
		c.JSON(http.StatusBadRequest, delivery.ErrorResponse{
			Error:   "missing_parameters",
			Message: "Both order_id and delivery_partner_id parameters are required",
			Code:    400,
		})
		return
	}

	// Validate and parse order ID
	orderID, err := strconv.Atoi(orderIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, delivery.ErrorResponse{
			Error:   "invalid_order_id",
			Message: "Invalid order ID format",
			Code:    400,
		})
		return
	}

	// Validate UUID format for delivery partner
	partnerID, err := uuid.Parse(partnerIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, delivery.ErrorResponse{
			Error:   "invalid_uuid",
			Message: "Invalid UUID format for delivery partner ID",
			Code:    400,
		})
		return
	}

	// Get order details
	orderDetail, err := postgres.GetOrderDetail(orderID, partnerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, delivery.ErrorResponse{
				Error:   "not_found",
				Message: "Order not found or not assigned to this delivery partner",
				Code:    404,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, delivery.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to retrieve order details",
			Code:    500,
		})
		return
	}

	c.JSON(http.StatusOK, orderDetail)
}
