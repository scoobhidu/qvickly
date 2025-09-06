package user_ec2

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qvickly/database/postgres"
	"strconv"
)

// 1. Get categories with max 5 items each
func GetCategoriesWithItems(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   postgres.GetDashboardItems(),
	})
}

// 2. Get all categories with their subcategories
func GetCategoriesWithSubcategories(c *gin.Context) {
	// Get all categories
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   postgres.GetCategoriesWithSubcategories(),
	})
}

// 3. Get all dashboard nudges
func GetDashboardNudges(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   postgres.GetNudges(),
	})
}

// Optional: Get items by category ID with pagination
func GetItemsByCategory(c *gin.Context) {
	categoryIDStr := c.DefaultQuery("subcategory_id", "1")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	// Get pagination parameters
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	items, _ := postgres.GetItemsBySubCategory(c, err, categoryID, limit, offset)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   items,
	})
}

// Optional: Get items by category ID with pagination
func GetItemsByFilter(c *gin.Context) {
	categoryIDStr := c.DefaultQuery("subcategory_id", "1")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	// Get pagination parameters
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	minP := c.DefaultQuery("min_price", "def")
	maxP := c.DefaultQuery("max_price", "def")

	searchQuery := c.DefaultQuery("search", "")

	minPrice := -1
	maxPrice := -1

	if minP != "def" {
		minPrice, _ = strconv.Atoi(minP)
	}
	if maxP != "def" {
		maxPrice, _ = strconv.Atoi(maxP)
	}

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	items, _ := postgres.GetItemsByFilter(categoryID, minPrice, maxPrice, searchQuery, limit, offset)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   items,
	})
}
