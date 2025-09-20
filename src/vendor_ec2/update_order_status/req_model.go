package update_order_status

// Valid order statuses
var ValidOrderStatuses = map[string]bool{
	"pending":   true,
	"accepted":  true,
	"packed":    true,
	"picked":    true,
	"delivered": true,
	"cancelled": true,
	"rejected":  true,
}

// Request/Response structures
type UpdateOrderStatusRequest struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

type UpdateOrderStatusResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
