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
// @Param request body vendors.GetVendorProfileRequestBody true "Complete vendor profile information"
// @Router /vendor/profile/details [post]
func GetVendorProfileDetails(context *gin.Context) {
	var json vendors.GetVendorProfileRequestBody
	if err := context.ShouldBindJSON(&json); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	vendorDetails, err := postgres.GetVendorProfile(json.Phone, json.Password)
	if err != nil {
		context.JSON(500, gin.H{"error": "Failed to get vendor profile" + err.Error()})
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
