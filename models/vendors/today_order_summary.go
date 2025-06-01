package vendors

// TodayOrderSummary represents order count statistics for the current day grouped by status
type TodayOrderSummary struct {
	Pending   int `json:"pending" example:"12" description:"Number of orders that are pending vendor acceptance"`
	Accepted  int `json:"accepted" example:"8" description:"Number of orders accepted by vendor but not yet packed"`
	Packed    int `json:"packed" example:"15" description:"Number of orders that have been packed and ready for pickup"`
	Ready     int `json:"ready" example:"6" description:"Number of orders ready for delivery or customer pickup"`
	Completed int `json:"completed" example:"45" description:"Number of orders that have been successfully delivered/completed"`
	Cancelled int `json:"cancelled" example:"3" description:"Number of orders cancelled by customer or vendor"`
	Rejected  int `json:"rejected" example:"2" description:"Number of orders rejected by vendor"`
}
