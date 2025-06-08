package update_location

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"qvickly/database/postgres"
	"qvickly/models/delivery"
)

// UpdateDeliveryPartnerLocation godoc
// @Summary Update delivery partner location
// @Description Update the current GPS location of a delivery partner and set them online
// @Tags Delivery Partners
// @Accept json
// @Produce json
// @Param delivery_partner_id query string true "Delivery Partner UUID" format(uuid) example("de111111-2222-3333-4444-555555555555")
// @Param location body delivery.UpdateLocationRequest true "GPS coordinates"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /delivery/update_location [post]
func UpdateDeliveryPartnerLocation(c *gin.Context) {
	// Get delivery partner ID from query parameter
	partnerIDParam := c.Query("delivery_partner_id")
	if partnerIDParam == "" {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	// Validate UUID format
	partnerID, err := uuid.Parse(partnerIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	// Parse request body
	var request delivery.UpdateLocationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	// Validate coordinates
	if !isValidCoordinates(request.Lat, request.Long) {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	// Update location
	err = postgres.ProcessLocationUpdate(partnerID, request.Lat, request.Long)
	if err != nil {
		switch err.Error() {
		case "partner_not_found":
			c.JSON(http.StatusNotFound, nil)
		default:
			c.JSON(http.StatusInternalServerError, nil)
		}
		return
	}

	c.JSON(http.StatusOK, nil)
}

// isValidCoordinates validates latitude and longitude ranges
func isValidCoordinates(lat, long float64) bool {
	// Validate latitude range (-90 to 90)
	if lat < -90 || lat > 90 {
		return false
	}

	// Validate longitude range (-180 to 180)
	if long < -180 || long > 180 {
		return false
	}

	// Check for obviously invalid coordinates (0,0 is unlikely for real deliveries)
	if lat == 0 && long == 0 {
		return false
	}

	return true
}
