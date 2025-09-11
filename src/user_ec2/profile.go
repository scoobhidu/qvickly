package user_ec2

import (
	"database/sql"
	"errors"
	"net/http"
	"qvickly/database/postgres"
	"qvickly/models/user"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 1. Customer Login API
func LoginCustomer(c *gin.Context) {
	var req user.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, user.LoginResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}
	customer, err := postgres.LoginC(req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, user.LoginResponse{
				Success: false,
				Message: "Customer not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, user.LoginResponse{
			Success: false,
			Message: "Database error",
		})
		return
	}

	c.JSON(http.StatusOK, user.LoginResponse{
		Success:  true,
		Message:  "Login successful",
		Customer: &customer,
	})
}

// 2. Add Customer Address API
func AddCustomerAddress(c *gin.Context) {
	customerID := c.Param("customer_id")
	if customerID == "" {
		c.JSON(http.StatusBadRequest, user.AddressResponse{
			Success: false,
			Message: "Customer ID is required",
		})
		return
	}

	var req user.AddAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, user.AddressResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	cID, err := uuid.Parse(customerID)

	err = postgres.AddCAddress(req, cID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, user.AddressResponse{
			Success: false,
			Message: "Failed to add address",
		})
		return
	}
	c.JSON(http.StatusCreated, user.AddressResponse{
		Success: true,
		Message: "Address added successfully",
	})
}

// 3. Mark Address as Default API
func MarkAddressDefault(c *gin.Context) {
	customerID := c.Param("customer_id")
	if customerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Customer ID is required",
		})
		return
	}

	var req user.MarkDefaultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request format",
		})
		return
	}
	cID, err := uuid.Parse(customerID)

	err = postgres.SetDefaultAddress(req, cID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update addresses",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Address marked as default successfully",
	})
}

// 4. Get Coupons API
func GetCoupons(c *gin.Context) {
	coupons, _ := postgres.GetAllCoupons()

	c.JSON(http.StatusOK, user.CouponsResponse{
		Success: true,
		Message: "Coupons fetched successfully",
		Coupons: coupons,
	})
}

// 5. Apply Coupon API
//func ApplyCoupon(c *gin.Context) {
//	var req user.ApplyCouponRequest
//	if err := c.ShouldBindJSON(&req); err != nil {
//		c.JSON(http.StatusBadRequest, user.ApplyCouponResponse{
//			Success: false,
//			Message: "Invalid request format",
//		})
//		return
//	}
//
//	ctx := context.Background()
//
//	// Get order details
//	var orderAmount float64
//	orderQuery := `SELECT amount FROM quickkart.customer.orders WHERE order_id = $1`
//	err := pgPool.QueryRow(ctx, orderQuery, req.OrderID).Scan(&orderAmount)
//	if err != nil {
//		c.JSON(http.StatusNotFound, user.ApplyCouponResponse{
//			Success: false,
//			Message: "Order not found",
//		})
//		return
//	}
//
//	// Get coupon details
//	var coupon user.Coupon
//	couponQuery := `
//		SELECT coupon_id, discount_type, discount_value, max_discount,
//			   min_order_value, max_usages, used_count
//		FROM quickkart.master.coupons
//		WHERE code = $1 AND is_active = true
//		AND valid_from <= CURRENT_TIMESTAMP AND valid_to > CURRENT_TIMESTAMP`
//
//	err = pgPool.QueryRow(ctx, couponQuery, req.CouponCode).Scan(
//		&coupon.ID, &coupon.DiscountType, &coupon.DiscountValue,
//		&coupon.MaxDiscount, &coupon.MinOrderValue, &coupon.MaxUsages, &coupon.UsedCount,
//	)
//
//	if err != nil {
//		c.JSON(http.StatusNotFound, user.ApplyCouponResponse{
//			Success: false,
//			Message: "Invalid or expired coupon",
//		})
//		return
//	}
//
//	// Validate coupon conditions
//	if orderAmount < coupon.MinOrderValue {
//		c.JSON(http.StatusBadRequest, user.ApplyCouponResponse{
//			Success: false,
//			Message: fmt.Sprintf("Minimum order value should be %.2f", coupon.MinOrderValue),
//		})
//		return
//	}
//
//	if coupon.UsedCount >= coupon.MaxUsages {
//		c.JSON(http.StatusBadRequest, user.ApplyCouponResponse{
//			Success: false,
//			Message: "Coupon usage limit exceeded",
//		})
//		return
//	}
//
//	// Calculate discount
//	var discountAmount float64
//	if coupon.DiscountType == "percentage" {
//		discountAmount = orderAmount * (coupon.DiscountValue / 100)
//		if coupon.MaxDiscount != nil && discountAmount > *coupon.MaxDiscount {
//			discountAmount = *coupon.MaxDiscount
//		}
//	} else if coupon.DiscountType == "fixed" {
//		discountAmount = coupon.DiscountValue
//	}
//
//	finalAmount := orderAmount - discountAmount
//
//	// Start transaction to update order and coupon usage
//	tx, err := pgPool.Begin(ctx)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, user.ApplyCouponResponse{
//			Success: false,
//			Message: "Database error",
//		})
//		return
//	}
//	defer tx.Rollback(ctx)
//
//	// Update order with coupon details
//	_, err = tx.Exec(ctx,
//		`UPDATE quickkart.customer.orders
//		 SET amount = $1, coupon_code = $2, discount_amount = $3
//		 WHERE order_id = $4`,
//		finalAmount, req.CouponCode, discountAmount, req.OrderID)
//
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, user.ApplyCouponResponse{
//			Success: false,
//			Message: "Failed to apply coupon",
//		})
//		return
//	}
//
//	// Update coupon usage count
//	_, err = tx.Exec(ctx,
//		`UPDATE quickkart.master.coupons SET used_count = used_count + 1 WHERE coupon_id = $1`,
//		coupon.ID)
//
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, user.ApplyCouponResponse{
//			Success: false,
//			Message: "Failed to update coupon usage",
//		})
//		return
//	}
//
//	// Commit transaction
//	if err = tx.Commit(ctx); err != nil {
//		c.JSON(http.StatusInternalServerError, user.ApplyCouponResponse{
//			Success: false,
//			Message: "Failed to commit changes",
//		})
//		return
//	}
//
//	c.JSON(http.StatusOK, user.ApplyCouponResponse{
//		Success:        true,
//		Message:        "Coupon applied successfully",
//		DiscountAmount: &discountAmount,
//		FinalAmount:    &finalAmount,
//	})
//}

// 6. Get Order Status API
func GetOrderStatus(c *gin.Context) {
	orderID := c.Param("order_id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, user.OrderStatusResponse{
			Success: false,
			Message: "Order ID is required",
		})
		return
	}

	oID, _ := uuid.Parse(orderID)
	order, err := postgres.GetOrderStatus(oID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, user.OrderStatusResponse{
			Success: false,
			Message: "Failed to fetch order status" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user.OrderStatusResponse{
		Success: true,
		Message: "Order status fetched successfully",
		Order:   &order,
	})
}

// Get All Customer Addresses API
func GetCustomerAddresses(c *gin.Context) {
	customerID := c.Query("customer_id")
	if customerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Customer ID is required",
		})
		return
	}
	cID, err := uuid.Parse(customerID)
	addresses, err := postgres.GetAddresses(cID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, user.OrderStatusResponse{})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"message":   "Addresses fetched successfully",
		"addresses": addresses,
	})
	return
}
