package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/naqeeb8a/Coffee-shop/initializers"
	"github.com/naqeeb8a/Coffee-shop/models"
)

func AddCategory(c *gin.Context) {
	var body struct {
		CategoryImage string
		Name          string
	}

	if (c.Bind(&body)) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	var validationError = make(map[string]string)
	if body.Name == "" {
		validationError["name"] = "Please provide a name"
	}
	if body.CategoryImage == "" {
		validationError["categoryImage"] = "Please provide an image"
	}
	if len(validationError) != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":    "Fill all fields",
			"required": validationError,
		})
		return
	}
	category := models.Category{CategoryImage: body.CategoryImage, Name: body.Name}

	result := initializers.DB.Create(&category)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error creating category",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"categoryImage": body.CategoryImage,
		"name":          body.Name,
	})
}
func AllCategories(c *gin.Context) {
	var categories []models.Category
	result := initializers.DB.Find(&categories)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error fetching all categories",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"categoriesCount": result.RowsAffected,
		"categories":      categories,
	})
}
func CategoryItem(c *gin.Context) {
	if c.Query("categoryId") == "" {
		c.AbortWithStatusJSON(401, gin.H{
			"error":               "Query parameter (categoryId) was not passed",
			"passed (categoryId)": c.Query("categoryId"),
		})
		return
	}
	var items []models.Item
	result := initializers.DB.Find(&items, "category_id = ?", c.Query("categoryId"))
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Error fetch all items of that category",
			"reason": result.Error,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"itemCount": result.RowsAffected,
		"items":     items,
	})
}
func EditCategory(c *gin.Context) {
	var body struct {
		CategoryImage string
		Name          string
	}

	if (c.Bind(&body)) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	if c.Query("categoryId") == "" {
		c.AbortWithStatusJSON(401, gin.H{
			"error":               "Query parameter (categoryId) was not passed",
			"passed (categoryId)": c.Query("categoryId"),
		})
		return
	}
	var category models.Category
	initializers.DB.First(&category, c.Query("categoryId"))
	if category.ID == 0 {
		c.AbortWithStatusJSON(401, gin.H{
			"error": "No record found of category with this id",
		})
		return
	}
	if body.Name != "" {
		category.Name = body.Name
	}
	if body.CategoryImage != "" {
		category.CategoryImage = body.CategoryImage
	}
	initializers.DB.Save(&category)
	c.JSON(http.StatusOK, gin.H{
		"categoryImage": category.CategoryImage,
		"name":          category.Name,
	})
}
