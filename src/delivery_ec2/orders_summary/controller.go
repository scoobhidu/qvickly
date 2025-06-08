package orders_summary

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"qvickly/database/postgres"
	"qvickly/models/delivery"
)

// GetDeliveryPartnerOrdersSummary godoc
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
func GetDeliveryPartnerOrdersSummary(c *gin.Context) {
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

	partnerID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, delivery.ErrorResponse{
			Error:   "invalid_uuid",
			Message: "Invalid UUID format for ID parameter",
			Code:    400,
		})
		return
	}

	response, err := postgres.GetBasicOrdersSummary(partnerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, delivery.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to retrieve orders summary",
			Code:    500,
		})
		return
	}
	c.JSON(http.StatusOK, response)
}
