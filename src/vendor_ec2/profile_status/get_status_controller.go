package profile_status

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qvickly/database/postgres"
)

// GetProfileStatus godoc
// @Summary Get profile live or not status
// @Description ”
// @Tags profile
// @Accept json
// @Produce json
// @Param vendor_id path string true "Vendor ID" format:"uuid" example:"123e4567-e89b-12d3-a456-426614174000"
// @Router /vendor/{vendor_id}/profile/status [get]
func GetProfileStatus(c *gin.Context) {
	vendorID := c.Param("vendor_id")

	if vendorID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "status": false, "error": "Invalid vendor_id or item_id"})
		return
	}

	status, err := postgres.GetProfileVendorStatus(vendorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "status": status, "error": "Failed to update inventory | " + err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": true, "status": status})
	}
}
