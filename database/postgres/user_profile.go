package postgres

import (
	"context"
	"database/sql"
	"errors"
	"qvickly/models/user"
	"time"

	"github.com/google/uuid"
)

func LoginC(req user.LoginRequest) (user.CustomerData, error) {
	ctx := context.Background()

	// Check if customer exists
	var customer user.CustomerData
	query := `SELECT id, full_name, phone, email FROM quickkart.profile.customer WHERE phone = $1`

	err := pgPool.QueryRow(ctx, query, req.Phone).Scan(
		&customer.ID, &customer.FullName, &customer.Phone, &customer.Email)

	return customer, err
}

func AddCAddress(req user.AddAddressRequest, customerID uuid.UUID) (err error) {
	ctx := context.Background()
	addressID := uuid.New()

	query := `
		INSERT INTO quickkart.profile.customer_addresses (
			address_id, customer_id, title, address_line1, address_line2, 
			city, state, postal_code, country, latitude, longitude, 
			is_default, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, false, $12, $12)`

	_, err = pgPool.Exec(ctx, query,
		addressID, customerID, req.Title, req.AddressLine1, req.AddressLine2,
		req.City, req.State, req.PostalCode, req.Country,
		req.Latitude, req.Longitude, time.Now())

	return err
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
