package postgres

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"qvickly/models/user"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func PlaceOrderInDB(req user.PlaceOrderRequest) (customerAddr *user.CustomerAddress, vendor *user.Vendor, deliveryBoy *user.DeliveryBoy, orderID uuid.UUID, totalAmount float64, err error) {
	tx, err := pgPool.Begin(context.Background())
	// 1. Get customer's delivery address
	customerAddr, err = getCustomerAddress(tx, req.CustomerID, req.AddressID)
	if err != nil {
		defer tx.Rollback(context.Background())
		err = errors.New("invalid customer address")
		return
	}

	// 2. Find the best vendor
	vendor, totalAmount, err = findBestVendor(tx, req.Items, customerAddr.Latitude, customerAddr.Longitude)
	if err != nil {
		defer tx.Rollback(context.Background())
		return
	}

	// 3. Create the main order
	orderID = uuid.New()
	err = createOrder(tx, orderID, req.CustomerID, totalAmount)
	if err != nil {
		return
	}

	// 4. Create order items
	err = createOrderItems(tx, orderID, req.Items)
	if err != nil {
		return
	}

	// 5. Create vendor order items
	err = createVendorOrderItems(tx, orderID, vendor.ID, req.Items)
	if err != nil {
		return
	}

	// 6. Create vendor pickup assignment
	vendorAssignmentID := uuid.New()
	err = createVendorPickupAssignment(tx, vendorAssignmentID, vendor.ID, orderID)
	if err != nil {
		return
	}

	// 7. Find and assign closest delivery boy
	deliveryBoy, err = findClosestDeliveryBoy(tx, vendor.Latitude, vendor.Longitude)
	if err != nil {
		return
	}

	// 8. Create delivery tracking record
	deliveryID := uuid.New()
	err = createDeliveryTracking(tx, deliveryID, orderID, deliveryBoy.ID, vendorAssignmentID)
	if err != nil {
		return
	}

	// 9. Update inventory
	err = updateInventory(tx, vendor.ID, req.Items)
	if err != nil {
		return
	}

	// Commit transaction
	if err = tx.Commit(context.Background()); err != nil {
		return
	}
	return
}

func findClosestDeliveryBoy(tx pgx.Tx, vendorLat, vendorLon float64) (*user.DeliveryBoy, error) {
	query := `
		SELECT id, full_name, latitude, longitude, is_active
		FROM quickkart.profile.delivery_boy
		WHERE is_active = true`

	rows, err := tx.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var closestBoy *user.DeliveryBoy
	var minDistance = math.MaxFloat64

	for rows.Next() {
		var boy user.DeliveryBoy
		err := rows.Scan(&boy.ID, &boy.Name, &boy.Latitude, &boy.Longitude, &boy.IsActive)
		if err != nil {
			continue
		}

		distance := calculateDistance(vendorLat, vendorLon, boy.Latitude, boy.Longitude)
		if distance < minDistance {
			minDistance = distance
			closestBoy = &boy
		}
	}

	if closestBoy == nil {
		return nil, fmt.Errorf("no active delivery boys available")
	}

	return closestBoy, nil
}

func createOrder(tx pgx.Tx, orderID, customerID uuid.UUID, amount float64) error {
	query := `
        INSERT INTO customer.orders (
            order_id, customer_id, order_time, amount, status
        ) VALUES ($1, $2, $3, $4, $5)`

	_, err := tx.Exec(context.Background(), query, orderID, customerID, time.Now(), amount, "placed")
	return err
}

func createOrderItems(tx pgx.Tx, orderID uuid.UUID, items []user.OrderItem) error {
	query := `
		INSERT INTO quickkart.customer.order_items (order_id, item_id, qty)
		VALUES ($1::uuid, $2, $3)`

	for _, item := range items {
		_, err := tx.Exec(context.Background(), query, orderID.String(), item.ItemID, item.Quantity)
		if err != nil {
			return err
		}
	}
	return nil
}

func createVendorOrderItems(tx pgx.Tx, orderID, vendorID uuid.UUID, items []user.OrderItem) error {
	// Insert items into vendor.order_items table
	query := `
        INSERT INTO quickkart.vendor.order_pickup_assignments(order_id, vendor_id)
        VALUES ($1, $2)`

	for _, item := range items {
		_, err := tx.Exec(context.Background(), query, orderID, vendorID)
		if err != nil {
			return fmt.Errorf("failed to insert vendor order item %d: %v", item.ItemID, err)
		}
	}
	return nil
}

func createVendorPickupAssignment(tx pgx.Tx, assignmentID, vendorID, orderID uuid.UUID) error {
	query := `
		INSERT INTO quickkart.vendor.order_pickup_assignments (
			vendor_assignment_id, vendor_id, order_id, picked_up, created_at
		) VALUES ($1, $2, $3, $4, $5)`

	_, err := tx.Exec(context.Background(), query, assignmentID, vendorID, orderID, false, time.Now())
	return err
}

func createDeliveryTracking(tx pgx.Tx, deliveryID, orderID, deliveryBoyID, vendorAssignmentID uuid.UUID) error {
	query := `
		INSERT INTO quickkart.delivery.order_tracker (
			delivery_id, order_id, delivery_boy_id, vendor_assignment_id
		) VALUES ($1, $2, $3, $4)`

	_, err := tx.Exec(context.Background(), query, deliveryID, orderID, deliveryBoyID, vendorAssignmentID)
	return err
}

func updateInventory(tx pgx.Tx, vendorID uuid.UUID, items []user.OrderItem) error {
	query := `
		UPDATE quickkart.vendor.inventory
		SET qty = qty - $1, updated_at = $2
		WHERE vendor_id = $3 AND item_id = $4`

	for _, item := range items {
		_, err := tx.Exec(context.Background(), query, item.Quantity, time.Now(), vendorID, item.ItemID)
		if err != nil {
			return err
		}
	}
	return nil
}

// Additional helper endpoints
//func getOrderStatus(c *gin.Context) {
//	orderID := c.Param("order_id")
//
//	query := `
//		SELECT
//			o.order_id,
//			o.status,
//			o.amount,
//			v.business_name as vendor_name,
//			db.full_name as delivery_boy_name,
//			ot.delivery_id
//		FROM quickkart.customer.orders o
//		LEFT JOIN quickkart.profile.vendors v ON o.vendor_id = v.vendor_id
//		LEFT JOIN quickkart.delivery.order_tracker ot ON o.order_id = ot.order_id
//		LEFT JOIN quickkart.profile.delivery_boy db ON ot.delivery_boy_id = db.id
//		WHERE o.order_id = $1`
//
//	var result struct {
//		OrderID         string  `json:"order_id"`
//		Status          string  `json:"status"`
//		Amount          float64 `json:"amount"`
//		VendorName      string  `json:"vendor_name"`
//		DeliveryBoyName string  `json:"delivery_boy_name"`
//		DeliveryID      string  `json:"delivery_id"`
//	}
//
//	err := pgPool.QueryRow(context.Background(), query, orderID).Scan(
//		&result.OrderID, &result.Status, &result.Amount,
//		&result.VendorName, &result.DeliveryBoyName, &result.DeliveryID,
//	)
//
//	if err != nil {
//		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
//		return
//	}
//
//	c.JSON(http.StatusOK, result)
//}

func updateDeliveryStatus(c *gin.Context) {
	orderID := c.Param("order_id")

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the appropriate timestamp based on status
	var query string
	switch req.Status {
	case "packed":
		query = "UPDATE quickkart.customer.orders SET pack_by_time = $1 WHERE order_id = $2"
	case "picked_up":
		query = "UPDATE quickkart.customer.orders SET pick_up_time = $1 WHERE order_id = $2"
	case "delivered":
		query = "UPDATE quickkart.customer.orders SET delivery_time = $1 WHERE order_id = $2"
	case "paid":
		query = "UPDATE quickkart.customer.orders SET paid_time = $1 WHERE order_id = $2"
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	_, err := pgPool.Exec(context.Background(), query, time.Now(), orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
}

func getCustomerAddress(tx pgx.Tx, customerID uuid.UUID, addressID int64) (*user.CustomerAddress, error) {
	var addr user.CustomerAddress
	query := `
		SELECT address_id, latitude, longitude
		FROM quickkart.profile.customer_addresses
		WHERE customer_id = $1 AND address_id = $2`

	err := tx.QueryRow(context.Background(), query, customerID, addressID).Scan(&addr.ID, &addr.Latitude, &addr.Longitude)
	return &addr, err
}

func findBestVendor(tx pgx.Tx, items []user.OrderItem, customerLat, customerLon float64) (*user.Vendor, float64, error) {
	// Get all active and live vendors
	vendorQuery := `
		SELECT vendor_id, business_name, latitude, longitude, is_active, is_active
		FROM quickkart.profile.vendors
		WHERE is_active = true AND is_active = true`

	rows, err := tx.Query(context.Background(), vendorQuery)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var vendors []user.Vendor
	for rows.Next() {
		var v user.Vendor
		err := rows.Scan(&v.ID, &v.Name, &v.Latitude, &v.Longitude, &v.IsActive, &v.IsLive)
		if err != nil {
			continue
		}
		vendors = append(vendors, v)
	}

	if len(vendors) == 0 {
		return nil, 0, fmt.Errorf("no active vendors available")
	}

	// Check which vendors have all required items in stock
	var bestVendor *user.Vendor
	var bestDistance = math.MaxFloat64
	var totalAmount float64

	for _, vendor := range vendors {
		hasAllItems, amount, err := checkVendorInventory(tx, vendor.ID, items)
		if err != nil {
			continue
		}

		if hasAllItems {
			distance := calculateDistance(customerLat, customerLon, vendor.Latitude, vendor.Longitude)
			if distance < bestDistance {
				bestDistance = distance
				bestVendor = &vendor
				totalAmount = amount
			}
		}
	}

	if bestVendor == nil {
		return nil, 0, fmt.Errorf("no vendor has all required items in stock")
	}

	return bestVendor, totalAmount, nil
}

func checkVendorInventory(tx pgx.Tx, vendorID uuid.UUID, items []user.OrderItem) (bool, float64, error) {
	// Create array of item IDs for the query
	itemIDs := make([]int64, len(items))
	itemQuantities := make(map[int64]int)

	for i, item := range items {
		itemIDs[i] = item.ItemID
		itemQuantities[item.ItemID] = item.Quantity
	}

	query := `
		SELECT
			vi.item_id,
			vi.qty,
			COALESCE(vi.wholesale_price_override, gi.price_wholesale) as price
		FROM quickkart.vendor.inventory vi
		JOIN quickkart.master.grocery_items gi ON vi.item_id = gi.item_id
		WHERE vi.vendor_id = $1 AND vi.item_id = ANY($2)`

	rows, err := tx.Query(context.Background(), query, vendorID, itemIDs)
	if err != nil {
		return false, 0, err
	}
	defer rows.Close()

	availableItems := make(map[int64]int)
	itemPrices := make(map[int64]float64)

	for rows.Next() {
		var itemID int64
		var qty int
		var price float64

		err := rows.Scan(&itemID, &qty, &price)
		if err != nil {
			continue
		}

		availableItems[itemID] = qty
		itemPrices[itemID] = price
	}

	// Check if all items are available in sufficient quantity
	var totalAmount float64
	for _, item := range items {
		availableQty, exists := availableItems[item.ItemID]
		if !exists || availableQty < item.Quantity {
			return false, 0, nil
		}
		totalAmount += itemPrices[item.ItemID] * float64(item.Quantity)
	}

	return true, totalAmount, nil
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
