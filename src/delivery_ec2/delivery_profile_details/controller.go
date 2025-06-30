package delivery_profile_details

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"qvickly/database/postgres"
	"qvickly/models/delivery"
	"qvickly/models/vendors"
)

// GetDeliveryPartnerProfile godoc
// @Summary Get delivery partner profile details
// @Description Retrieve detailed profile information for a delivery partner by ID
// @Tags Delivery Partners
// @Accept json
// @Produce json
// @Param id query string true "Delivery Partner UUID" format(uuid) example("de111111-2222-3333-4444-555555555555")
// @Success 200 {object} delivery.DeliveryProfileDetailsSuccessResponse "Successfully retrieved delivery partner profile"
// @Success 200 {object} delivery.DeliveryProfileDetailsSuccessResponse "Delivery partner profile details"
// @Failure 400 {object} delivery.ErrorResponse "Invalid UUID format or missing ID parameter"
// @Failure 404 {object} delivery.ErrorResponse "Delivery partner not found"
// @Failure 500 {object} delivery.ErrorResponse "Internal server error"
// @Router /delivery/profile/details [post]
func GetDeliveryPartnerProfile(c *gin.Context) {
	var json vendors.ProfileRequestBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, dob, err := postgres.GetDeliveryPartnerProfileDetails(json.Phone, json.Password)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, delivery.ErrorResponse{
				Error:   "not_found",
				Message: "Delivery partner not found",
				Code:    404,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, delivery.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to retrieve delivery partner profile",
			Code:    500,
		})
		return
	}

	// Handle nullable date of birth
	if dob.Valid {
		data.DateOfBirth = &dob.String
	}

	// Return success response
	c.JSON(http.StatusOK, delivery.DeliveryProfileDetailsSuccessResponse{
		Success: true,
		Data:    data,
		Message: "Delivery partner profile retrieved successfully",
	})
}
