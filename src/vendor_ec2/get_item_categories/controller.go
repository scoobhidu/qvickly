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
// @Router /vendor/categories [get]
func GetCategoriesHandler(c *gin.Context) {
	categories, err := postgres.GetItemCategories()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to fetch categories"})
		return
	}

	m := make(map[string][]interface{})

	for _, v := range categories {
		l := make([]interface{}, 0)

		l = append(l, v.Name)
		l = append(l, v.ID)

		m[v.SuperCategory] = append(m[v.SuperCategory], l)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": m})
}
