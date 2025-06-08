package vendor_profile_details

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qvickly/database/postgres"
	"qvickly/models/vendors"
)

// GetVendorProfileDetails godoc
// @Summary Get Vendor Profile Details
// @Description Retrieve complete profile information for a specific vendor including business details, location, and operating hours
// @Tags vendor-profile
// @Accept json
// @Produce json
// @Param id query string true "Vendor ID" format:"uuid" example:"123e4567-e89b-12d3-a456-426614174000"
// @Router /vendor/profile/details [get]
func GetVendorProfileDetails(context *gin.Context) {
	vendorID := context.Query("id")
	vendorDetails, err := postgres.GetVendorProfile(vendorID)
	if err != nil {
		context.JSON(500, gin.H{"error": "Failed to get vendor profile"})
		return
	} else {
		context.JSON(200, gin.H{"vendor_details": vendorDetails})
	}
}

// TODO | add an API to push vendor's image in S3 and return the public image URL

// CreateVendorProfileDetails godoc
// @Summary Create Vendor Profile
// @Description Create a new vendor profile with complete business information, location details, and operating hours
// @Tags vendor-profile
// @Accept json
// @Produce json
// @Param request body vendors.CompleteVendorProfile true "Complete vendor profile information"
// @Router /vendor/profile/details/create [post]
func CreateVendorProfileDetails(context *gin.Context) {
	var json vendors.CompleteVendorProfile
	if err := context.ShouldBindJSON(&json); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := postgres.AddVendorProfile(json)
	if err != nil {
		context.JSON(500, gin.H{"error": "Failed to get vendor profile"})
		return
	} else {
		context.JSON(200, gin.H{"message": "Successfully created vendor profile"})
	}
}
