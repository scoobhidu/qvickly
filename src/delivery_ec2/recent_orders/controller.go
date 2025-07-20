package recent_orders

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"qvickly/database/postgres"
	"qvickly/models/delivery"
	"strconv"
)

// GetDeliveryPartnerRecentOrders godoc
// @Summary Get recent orders for delivery partner
// @Description Retrieve recent orders assigned to a delivery partner with status and earnings information
// @Tags Delivery Partners
// @Accept json
// @Produce json
// @Param id query string true "Delivery Partner UUID" format(uuid) example("de111111-2222-3333-4444-555555555555")
// @Param limit query int false "Number of recent orders to fetch (default: 20, max: 100)" example(20)
// @Param status query string false "Filter by order status (pending, accepted, packed, ready, completed, cancelled, rejected)" example("completed")
// @Param detailed query boolean false "Include detailed order information" example(false)
// @Success 200 {array} delivery.RecentOrderResponse "List of recent orders"
// @Failure 400 {object} delivery.ErrorResponse "Invalid UUID format, missing ID parameter, or invalid limit"
// @Failure 404 {object} delivery.ErrorResponse "Delivery partner not found"
// @Failure 500 {object} delivery.ErrorResponse "Internal server error"
// @Router /delivery/orders/recent [get]
func GetDeliveryPartnerRecentOrders(c *gin.Context) {
	// Get ID from query parameter
	idParam := c.Query("id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, delivery.ErrorResponse{
			Error:   "missing_parameter",
			Message: "ID parameter is required",
			Code:    400,
		})
		return
	}

	// Validate UUID format
	partnerID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, delivery.ErrorResponse{
			Error:   "invalid_uuid",
			Message: "Invalid UUID format for ID parameter",
			Code:    400,
		})
		return
	}

	// Get limit parameter (default: 20, max: 100)
	limit := 20
	if limitParam := c.Query("limit"); limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil {
			if parsedLimit > 0 && parsedLimit <= 100 {
				limit = parsedLimit
			} else if parsedLimit > 100 {
				limit = 100
			}
		}
	}

	// Get status filter (optional)
	statusFilter := c.Query("status")

	orders, err := postgres.GetBasicRecentOrders(partnerID, limit, statusFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, delivery.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to retrieve recent orders",
			Code:    500,
		})
		return
	}
	c.JSON(http.StatusOK, orders)
}

// GetDeliveryPartnerAllOrders godoc
// @Summary Get recent orders for delivery partner
// @Description Retrieve recent orders assigned to a delivery partner with status and earnings information
// @Tags Delivery Partners
// @Accept json
// @Produce json
// @Param id query string true "Delivery Partner UUID" format(uuid) example("de111111-2222-3333-4444-555555555555")
// @Param limit query int false "Number of recent orders to fetch (default: 20, max: 100)" example(20)
// @Param status query string false "Filter by order status (pending, accepted, packed, ready, completed, cancelled, rejected)" example("completed")
// @Param detailed query boolean false "Include detailed order information" example(false)
// @Success 200 {array} delivery.RecentOrderResponse "List of recent orders"
// @Failure 400 {object} delivery.ErrorResponse "Invalid UUID format, missing ID parameter, or invalid limit"
// @Failure 404 {object} delivery.ErrorResponse "Delivery partner not found"
// @Failure 500 {object} delivery.ErrorResponse "Internal server error"
// @Router /delivery/orders/all [get]
func GetDeliveryPartnerAllOrders(c *gin.Context) {
	// Get ID from query parameter
	idParam := c.Query("id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, delivery.ErrorResponse{
			Error:   "missing_parameter",
			Message: "ID parameter is required",
			Code:    400,
		})
		return
	}

	// Validate UUID format
	partnerID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, delivery.ErrorResponse{
			Error:   "invalid_uuid",
			Message: "Invalid UUID format for ID parameter",
			Code:    400,
		})
		return
	}

	// Get status filter (optional)
	statusFilter := c.Query("status")

	orders, err := postgres.GetBasicAllOrders(partnerID, statusFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, delivery.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to retrieve recent orders",
			Code:    500,
		})
		return
	}
	c.JSON(http.StatusOK, orders)
}
