package update_location

import (
	"github.com/gin-gonic/gin"
)

// UpdateDeliveryPartnerLocation godoc
// @Summary Update delivery partner location and online status
// @Description Update the current location and online status of a delivery partner
// @Tags Delivery Partners
// @Accept json
// @Produce json
// @Param id query string true "Delivery Partner UUID" format(uuid)
// @Param latitude query number true "Latitude coordinate" example(28.6139391)
// @Param longitude query number true "Longitude coordinate" example(77.2090212)
// @Param online query boolean false "Online status" example(true)
// @Success 200 {object} map[string]interface{} "Location updated successfully"
// @Failure 400 {object} ErrorResponse "Invalid parameters"
// @Failure 404 {object} ErrorResponse "Delivery partner not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /delivery/profile/location [put]
func UpdateDeliveryPartnerLocation(c *gin.Context) {
	// Implementation for updating location
	// This is a bonus endpoint for location updates
}
