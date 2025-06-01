package search_items_for_inventory_addition

//import (
//	"github.com/gin-gonic/gin"
//	"net/http"
//	"strconv"
//)

//func SearchItemsToAddHandler(c *gin.Context) {
//	vendorID := c.Param("vendor_id")
//	search := c.Query("search")
//	categoryID := c.Query("category_id")
//
//	if vendorID == "" {
//		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "vendor_id is required"})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
//}
