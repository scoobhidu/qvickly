package delivery

// UpdateLocationRequest represents the location update request
type UpdateLocationRequest struct {
	Lat  float64 `json:"lat" binding:"required" example:"28.6139391"`
	Long float64 `json:"long" binding:"required" example:"77.2090212"`
}
