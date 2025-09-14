package user_ec2

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"net/http"
	"qvickly/database/postgres"
	"qvickly/models/user"
	"time"

	"github.com/gin-gonic/gin"
)

// Main order placement handler
func PlaceOrder(c *gin.Context) {
	var req user.PlaceOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customerAddr, vendor, deliveryBoy, orderID, totalAmount, err := postgres.PlaceOrderInDB(req)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Vendor not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error | " + err.Error()})
			return
		}
	}
	// Calculate estimated delivery time (vendor distance + delivery distance)
	vendorDistance := calculateDistance(customerAddr.Latitude, customerAddr.Longitude, vendor.Latitude, vendor.Longitude)
	deliveryDistance := calculateDistance(vendor.Latitude, vendor.Longitude, deliveryBoy.Latitude, deliveryBoy.Longitude)
	estimatedMinutes := int((vendorDistance + deliveryDistance) * 3) // Assume 3 minutes per km
	estimatedTime := time.Now().Add(time.Duration(estimatedMinutes) * time.Minute).Format("15:04")

	response := user.PlaceOrderResponse{
		OrderID:          orderID,
		AssignedVendorID: vendor.ID,
		DeliveryBoyID:    deliveryBoy.ID,
		EstimatedTime:    estimatedTime,
		TotalAmount:      totalAmount,
		Message:          fmt.Sprintf("Order placed successfully! Assigned to %s, delivery by %s", vendor.Name, deliveryBoy.Name),
	}

	c.JSON(http.StatusOK, response)
}

// Utility functions
func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth's radius in kilometers

	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}
