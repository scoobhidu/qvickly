package profile_status

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qvickly/database/postgres"
)

type RequestBody struct {
	Status bool `json:"status"`
}

// SetProfileStatus godoc
// @Summary Get profile live or not status
// @Description ”"
// @Tags profile
// @Accept json
// @Produce json
// @Param vendor_id path string true "Vendor ID" format:"uuid" example:"123e4567-e89b-12d3-a456-426614174000"
// @Param request body RequestBody true "Status to be set"
// @Router /vendor/{vendor_id}/profile/status [post]
func SetProfileStatus(c *gin.Context) {
	vendorID := c.Param("vendor_id")
	var req RequestBody
	if vendorID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid vendor_id or item_id"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid JSON"})
	}

	if err := postgres.SetProfileVendorStatus(vendorID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to update profile status" + err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Status updated successfully"})
	}
}
