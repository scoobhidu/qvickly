package postgres

import (
	"context"
	"database/sql"
	"errors"
	"qvickly/models/user"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func LoginC(req user.LoginRequest) (user.CustomerData, error) {
	ctx := context.Background()

	// Check if customer exists
	var customer user.CustomerData
	query := `SELECT id, full_name, phone, email, ca.latitude, ca.longitude, ca.title, COALESCE(ca.address_line1, ca.address_line2, ca.city, ca.state, ca.country, ca.postal_code)  FROM quickkart.profile.customer left join quickkart.profile.customer_addresses ca on customer.id = ca.customer_id WHERE phone = $1 and ca.is_default=true`

	err := pgPool.QueryRow(ctx, query, req.Phone).Scan(
		&customer.ID, &customer.FullName, &customer.Phone, &customer.Email, &customer.Latitude, &customer.Longitude, &customer.Title, &customer.Address)

	return customer, err
}

func AddCAddress(req user.AddAddressRequest, customerID uuid.UUID) (err error) {
	ctx := context.Background()

	query := `
		INSERT INTO quickkart.profile.customer_addresses (
			customer_id, title, address_line1, address_line2, 
			city, state, postal_code, country, latitude, longitude, 
			created_at, updated_at, is_default
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $11, CASE 
				WHEN NOT EXISTS (
					SELECT 1 
					FROM quickkart.profile.customer_addresses 
					WHERE customer_id = $1
				) THEN true 
				ELSE false
				END)`

	_, err = pgPool.Exec(ctx, query,
		customerID, req.Title, req.AddressLine1, req.AddressLine2,
		req.City, req.State, req.PostalCode, req.Country,
		req.Latitude, req.Longitude, time.Now())

	return err
}

func AddUser(req user.SignUpRequest) (cID string, err error) {
	ctx := context.Background()

	queryRow := `
		INSERT INTO quickkart.profile.customer (
			phone, full_name, email
		) VALUES ($1, coalesce($2, $3), $4) returning id`
	err = pgPool.QueryRow(ctx, queryRow, req.Phone, req.FirstName, req.LastName, req.Email).Scan(&cID)

	return
}

func SetDefaultAddress(req user.MarkDefaultRequest, customerID uuid.UUID) (err error) {
	ctx := context.Background()

	// Start transaction
	tx, err := pgPool.Begin(ctx)
	if err != nil {
		return
	}
	defer tx.Rollback(ctx)

	// First, set all addresses as non-default for this customer
	_, err = tx.Exec(ctx,
		`UPDATE quickkart.profile.customer_addresses SET is_default = false WHERE customer_id = $1`,
		customerID)

	// Then set the specified address as default
	result, err := tx.Exec(ctx,
		`UPDATE quickkart.profile.customer_addresses SET is_default = true 
		 WHERE address_id = $1 AND customer_id = $2`,
		req.AddressID, customerID)

	if result.RowsAffected() == 0 {
		return errors.New("address not found")
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		return errors.New("failed to commit changes")
	}

	return err
}

func GetAllCoupons() ([]user.Coupon, error) {
	ctx := context.Background()

	query := `
		SELECT coupon_id, title, code, discount_type, discount_value, 
			   max_discount, min_order_value, max_usages, used_count, 
			   valid_from, valid_to, is_active, created_at
		FROM quickkart.master.coupons 
		WHERE is_active = true AND valid_to > CURRENT_TIMESTAMP
		ORDER BY created_at DESC`

	rows, err := pgPool.Query(ctx, query)
	if err != nil {
		return []user.Coupon{}, err
	}
	defer rows.Close()

	var coupons []user.Coupon
	for rows.Next() {
		var coupon user.Coupon
		err := rows.Scan(
			&coupon.ID, &coupon.Title, &coupon.Code, &coupon.DiscountType,
			&coupon.DiscountValue, &coupon.MaxDiscount, &coupon.MinOrderValue,
			&coupon.MaxUsages, &coupon.UsedCount, &coupon.ValidFrom,
			&coupon.ValidTo, &coupon.IsActive, &coupon.CreatedAt,
		)
		if err != nil {
			continue
		}
		coupons = append(coupons, coupon)
	}

	return coupons, err
}

func GetOrderStatus(orderID uuid.UUID) (user.OrderStatus, error) {
	ctx := context.Background()

	var order user.OrderStatus
	query := `
		SELECT order_id, status, order_time, amount, pack_by_time, paid_time, delivery_time
		FROM quickkart.customer.orders 
		WHERE order_id = $1`

	err := pgPool.QueryRow(ctx, query, orderID).Scan(
		&order.OrderID, &order.Status, &order.OrderTime, &order.Amount,
		&order.PackByTime, &order.PaidTime, &order.DeliveryTime,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return order, err
		}
		return order, err
	}

	return order, err
}

func GetAddresses(customerID uuid.UUID) (addresses []gin.H, err error) {
	ctx := context.Background()

	query := `
		SELECT address_id, title, address_line1, address_line2, city, state, 
			   postal_code, country, latitude, longitude, is_default, created_at
		FROM quickkart.profile.customer_addresses 
		WHERE customer_id = $1
		ORDER BY is_default DESC, created_at DESC`

	rows, err := pgPool.Query(ctx, query, customerID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var addressID, title, addressLine1, city, state, postalCode, country string
		var addressLine2 *string
		var latitude, longitude *float64
		var isDefault bool
		var createdAt time.Time

		err := rows.Scan(
			&addressID, &title, &addressLine1, &addressLine2, &city, &state,
			&postalCode, &country, &latitude, &longitude, &isDefault, &createdAt,
		)
		if err != nil {
			continue
		}

		address := gin.H{
			"address_id":    addressID,
			"title":         title,
			"address_line1": addressLine1,
			"address_line2": addressLine2,
			"city":          city,
			"state":         state,
			"postal_code":   postalCode,
			"country":       country,
			"latitude":      latitude,
			"longitude":     longitude,
			"is_default":    isDefault,
			"created_at":    createdAt,
		}
		addresses = append(addresses, address)
	}

	if err = rows.Err(); err != nil {
		return
	}

	return
}
