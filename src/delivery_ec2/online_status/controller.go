package online_status

import (
	"github.com/gin-gonic/gin"
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
// @Param location body delivery.UpdateOnlineStatusRequest true "GPS coordinates"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /delivery/update_location [post]
func UpdateDeliveryPartnerOnlineStatus(c *gin.Context) {
	// Parse request body
	var request delivery.UpdateOnlineStatusRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	// Update location
	err := postgres.UpdateDeliveryPartnerOnlineStatus(request.DeliveryId, request.OnlineStatus)
	if err != nil {
		switch err.Error() {
		case "partner_not_found":
			c.JSON(http.StatusNotFound, nil)
		default:
			c.JSON(http.StatusInternalServerError, nil)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
