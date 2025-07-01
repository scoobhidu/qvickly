package delivery

// UpdateLocationRequest represents the location update request
type UpdateLocationRequest struct {
	Lat  float64 `json:"lat" binding:"required" example:"28.6139391"`
	Long float64 `json:"long" binding:"required" example:"77.2090212"`
}

// UpdateLocationRequest represents the location update request
type UpdateOnlineStatusRequest struct {
	DeliveryId   string `json:"id" binding:"required" example:"f958b442-531d-4b79-b216-3ed11311aef9"`
	OnlineStatus bool   `json:"status" binding:"required" example:"false"`
}
