package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"qvickly/src/delivery_ec2/delivery_profile_details"
	order_details2 "qvickly/src/delivery_ec2/order_details"
	"qvickly/src/delivery_ec2/orders_summary"
	recentorders2 "qvickly/src/delivery_ec2/recent_orders"
	"qvickly/src/delivery_ec2/update_location"
	"qvickly/src/vendor_ec2/add_item_to_inventory"
	"qvickly/src/vendor_ec2/inventory_items"
	"qvickly/src/vendor_ec2/inventory_summary"
	"qvickly/src/vendor_ec2/order_details"
	"qvickly/src/vendor_ec2/profile_details"
	"qvickly/src/vendor_ec2/recent_orders"
	"qvickly/src/vendor_ec2/remove_item_from_inventory"
	vendororders "qvickly/src/vendor_ec2/summary"
	"qvickly/src/vendor_ec2/update_inventory"
	"qvickly/src/vendor_ec2/update_order_status"
)

func Router(app *gin.Engine) {
	// Swagger endpoint
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	{
		group := app.Group("/vendor")

		group.POST("/profile/details", vendor_profile_details.GetVendorProfileDetails)
		group.POST("/profile/details/create", vendor_profile_details.CreateVendorProfileDetails)

		group.GET("/orders/summary", vendororders.GetVendorOrderSummary)
		group.GET("/orders/order_details", order_details.GetVendorOrderDetail)
		group.GET("/orders/recent_orders", recent_orders.GetVendorOrdersHandler)

		group.POST("/orders/update_order_status", update_order_status.UpdateOrderStatusHandler)

		group.GET("/:vendor_id/inventory/summary", inventory_summary.GetVendorInventorySummaryHandler)
		group.GET("/:vendor_id/inventory", get_vendor_items.GetVendorInventoryHandler)
		//group.GET("/:vendor_id/inventory/search", SearchItemsToAddHandler)
		group.POST("/:vendor_id/inventory", add_item_to_inventory.AddItemToInventoryHandler)
		group.PUT("/:vendor_id/inventory/:item_id", update_inventory.UpdateInventoryItemHandler)
		group.DELETE("/:vendor_id/inventory/:item_id", remove_item_from_inventory.RemoveItemFromInventoryHandler)
		//group.POST("/:vendor_id/inventory/movement", BulkInventoryMovementHandler)
	}

	{
		group := app.Group("/delivery")

		group.GET("/profile/details", delivery_profile_details.GetDeliveryPartnerProfile)
		group.GET("/profile/orders/summary", orders_summary.GetDeliveryPartnerOrdersSummary)
		group.GET("/orders/recent", recentorders2.GetDeliveryPartnerRecentOrders)
		group.GET("/order/detail", order_details2.GetDeliveryOrderDetail)
		group.POST("/update_location", update_location.UpdateDeliveryPartnerLocation)
	}

}
