package add_item_to_inventory

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qvickly/database/postgres"
)

// GetCategoriesHandler godoc
// @Summary Get All Categories
// @Description Retrieve all available product categories for filtering and organization purposes
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {object} []vendors.Category "Categories retrieved successfully"
// @Router /api/categories [get]
func GetCategoriesHandler(c *gin.Context) {
	categories, err := postgres.GetItemCategories()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to fetch categories"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": categories})
}
