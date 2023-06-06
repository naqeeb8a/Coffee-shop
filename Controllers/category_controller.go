package controllers

import (
	"net/http"

	"coffee-shop.com/api/initializers"
	"coffee-shop.com/api/models"
	"github.com/gin-gonic/gin"
)

func AddCategory(c *gin.Context) {
	CategoryImage, err := c.FormFile("image")
	CategoryName := c.PostForm("name")
	if CategoryName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":    "Empty name",
			"required": "provide a string name in multipart form",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":    "failed to upload image",
			"required": "provide a file in multipart form",
		})
		return
	}
	err = c.SaveUploadedFile(CategoryImage, "assets/"+CategoryImage.Filename)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":    "failed to upload image",
			"required": "provide a file in multipart form",
		})
		return
	}
	category := models.Category{CategoryImage: "assets/" + CategoryImage.Filename, Name: CategoryName}

	result := initializers.DB.Create(&category)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error creating category",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"details": category,
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
func RemoveCategory(c *gin.Context) {
	if c.Query("categoryId") == "" {
		c.AbortWithStatusJSON(401, gin.H{
			"error":               "Query parameter (categoryId) was not passed",
			"passed (categoryId)": c.Query("categoryId"),
		})
		return
	}
	result := initializers.DB.Unscoped().Delete(&models.Category{}, c.Query("categoryId"))
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error deleting category",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "category deleted successfully",
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

	if c.Query("categoryId") == "" {
		c.AbortWithStatusJSON(401, gin.H{
			"error":               "Query parameter (categoryId) was not passed",
			"passed (categoryId)": c.Query("categoryId"),
		})
		return
	}
	CategoryImage, err := c.FormFile("image")
	CategoryName := c.PostForm("name")
	if err != nil && CategoryName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":    "fields cannot be empty",
			"required": "provide atleast one field",
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
	if err == nil && CategoryImage.Filename != "" {
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":    "failed to upload image",
				"required": "provide a file in multipart form",
			})
			return
		}
		err = c.SaveUploadedFile(CategoryImage, "assets/"+CategoryImage.Filename)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":    "failed to upload image",
				"required": "provide a file in multipart form",
			})
			return
		}
		category.CategoryImage = "assets/" + CategoryImage.Filename
	}

	if CategoryName != "" {
		category.Name = CategoryName
	}

	initializers.DB.Save(&category)
	c.JSON(http.StatusOK, gin.H{
		"categoryImage": category.CategoryImage,
		"name":          category.Name,
	})
}
