package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/coffee-shop5860322/api/initializers"
	"gitlab.com/coffee-shop5860322/api/models"
)

func AddItem(c *gin.Context) {
	var body struct {
		Name        string
		Image       string
		Description string
		Price       int
		IsEnabled   bool
		CategoryId  int
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
	if body.Image == "" {
		validationError["image"] = "Please provide an image"
	}
	if body.Description == "" {
		validationError["description"] = "Please provide a description"
	}
	if body.Price == 0 {
		validationError["price"] = "Please provide an price"
	}
	if body.CategoryId == 0 {
		validationError["categoryId"] = "Please provide a category Id"
	}
	if len(validationError) != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":    "Fill all fields",
			"required": validationError,
		})
		return
	}
	item := models.Item{Name: body.Name, Image: body.Image, Description: body.Description, Price: body.Price, IsEnabled: body.IsEnabled, CategoryId: body.CategoryId}

	result := initializers.DB.Create(&item)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error creating item",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name":        body.Name,
		"image":       body.Image,
		"description": body.Description,
		"price":       body.Price,
		"isEnabled":   body.IsEnabled,
		"categoryId":  body.CategoryId,
	})
}
func EditItem(c *gin.Context) {
	var body struct {
		Name        string
		Image       string
		Description string
		Price       int
		IsEnabled   bool
		CategoryId  int
	}

	if (c.Bind(&body)) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	if c.Query("itemId") == "" {
		c.AbortWithStatusJSON(401, gin.H{
			"error":           "Query parameter (itemId) was not passed",
			"passed (itemId)": c.Query("itemId"),
		})
		return
	}
	var item models.Item
	initializers.DB.First(&item, c.Query("itemId"))
	if item.ID == 0 {
		c.AbortWithStatusJSON(401, gin.H{
			"error": "No record found of item with this id",
		})
		return
	}
	if body.Name != "" {
		item.Name = body.Name
	}
	if body.Image != "" {
		item.Image = body.Image
	}
	if body.Description != "" {
		item.Description = body.Description
	}
	if body.Price != 0 {
		item.Price = body.Price
	}
	if body.CategoryId != 0 {
		item.CategoryId = body.CategoryId
	}
	initializers.DB.Save(&item)
	c.JSON(http.StatusOK, gin.H{
		"id":          item.ID,
		"name":        item.Name,
		"image":       item.Image,
		"description": item.Description,
		"price":       item.Price,
		"isEnabled":   item.IsEnabled,
		"categoryId":  item.CategoryId,
	})
}
func AllItems(c *gin.Context) {
	var items []models.Item
	result := initializers.DB.Find(&items)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error fetching items",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"itemsCount": result.RowsAffected,
		"items":      items,
	})
}
func ItemDetails(c *gin.Context) {
	if c.Query("itemId") == "" {
		c.AbortWithStatusJSON(401, gin.H{
			"error":           "Query parameter (itemId) was not passed",
			"passed (itemId)": c.Query("itemId"),
		})
		return
	}
	var item models.Item
	initializers.DB.First(&item, c.Query("itemId"))
	if item.ID == 0 {
		c.AbortWithStatusJSON(401, gin.H{
			"error": "No record found of item with this id",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"itemDetails": item})
}
