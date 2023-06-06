package controllers

import (
	"net/http"
	"strconv"

	"coffee-shop.com/api/initializers"
	"coffee-shop.com/api/models"
	"github.com/gin-gonic/gin"
)

func AddItem(c *gin.Context) {
	ItemImage, err := c.FormFile("image")
	ItemName := c.PostForm("name")
	ItemDescription := c.PostForm("description")
	ItemPrice := c.PostForm("price")
	ItemIsEnabled := c.PostForm("isEnabled")
	ItemCategoryId := c.PostForm("categoryId")
	var validationError = make(map[string]string)
	if ItemName == "" {
		validationError["name"] = "Please provide a name"
	}
	if err != nil {
		validationError["image"] = "Please provide an image"
	}
	if ItemDescription == "" {
		validationError["description"] = "Please provide a description"
	}
	if ItemPrice == "" {
		validationError["price"] = "Please provide an price"
	}
	if ItemIsEnabled == "" {
		ItemIsEnabled = "true"
	}
	if ItemCategoryId == "" {
		validationError["categoryId"] = "Please provide a category Id"
	}
	if len(validationError) != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":    "Fill all fields",
			"required": validationError,
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
	err = c.SaveUploadedFile(ItemImage, "assets/"+ItemImage.Filename)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":    "failed to upload image",
			"required": "provide a file in multipart form",
		})
		return
	}
	price, intErr := strconv.Atoi(ItemPrice)
	if intErr != nil {
		c.AbortWithStatusJSON(401, gin.H{
			"error": "Enter a valid price that can be parsed into integer",
		})
		return
	}
	categoyId, intErr := strconv.Atoi(ItemCategoryId)
	if intErr != nil {
		c.AbortWithStatusJSON(401, gin.H{
			"error": "Enter a valid categoryId that can be parsed into integer",
		})
		return
	}
	isEnabled, boolErr := strconv.ParseBool(ItemIsEnabled)
	if boolErr != nil {
		c.AbortWithStatusJSON(401, gin.H{
			"error": "Enter a valid boolean that can be parsed into bool",
		})
		return
	}
	item := models.Item{Name: ItemName, Image: "assets/" + ItemImage.Filename, Description: ItemDescription, Price: price, IsEnabled: isEnabled, CategoryId: categoyId}

	result := initializers.DB.Create(&item)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error creating item",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"details": item,
	})
}
func EditItem(c *gin.Context) {
	if c.Query("itemId") == "" {
		c.AbortWithStatusJSON(401, gin.H{
			"error":           "Query parameter (itemId) was not passed",
			"passed (itemId)": c.Query("itemId"),
		})
		return
	}
	ItemImage, err := c.FormFile("image")
	ItemName := c.PostForm("name")
	ItemDescription := c.PostForm("description")
	ItemPrice := c.PostForm("price")
	ItemIsEnabled := c.PostForm("isEnabled")
	ItemCategoryId := c.PostForm("categoryId")
	if err != nil && ItemName == "" && ItemPrice == "" && ItemDescription == "" && ItemIsEnabled == "" && ItemCategoryId == "" {
		c.AbortWithStatusJSON(401, gin.H{
			"error":    "fields cannot be empty",
			"required": "provide atleast one field",
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
	if ItemName != "" {
		item.Name = ItemName
	}
	if err == nil && ItemImage.Filename != "" {
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":    "failed to upload image",
				"required": "provide a file in multipart form",
			})
			return
		}
		err = c.SaveUploadedFile(ItemImage, "assets/"+ItemImage.Filename)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":    "failed to upload image",
				"required": "provide a file in multipart form",
			})
			return
		}
		item.Image = "assets/" + ItemImage.Filename
	}
	if ItemDescription != "" {
		item.Description = ItemDescription
	}
	if ItemIsEnabled != "" {
		isEnabled, boolErr := strconv.ParseBool(ItemIsEnabled)
		if boolErr != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "Enter a valid boolean that can be parsed into bool",
			})
			return
		}
		item.IsEnabled = isEnabled
	}
	if ItemPrice != "" {
		price, intErr := strconv.Atoi(ItemPrice)
		if intErr != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "Enter a valid price that can be parsed into integer",
			})
			return
		}
		item.Price = price
	}
	if ItemCategoryId != "" {
		categoyId, intErr := strconv.Atoi(ItemCategoryId)
		if intErr != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "Enter a valid categoryId that can be parsed into integer",
			})
			return
		}
		item.CategoryId = categoyId
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
func RemoveItem(c *gin.Context) {
	if c.Query("itemId") == "" {
		c.AbortWithStatusJSON(401, gin.H{
			"error":           "Query parameter (itemId) was not passed",
			"passed (itemId)": c.Query("itemId"),
		})
		return
	}
	result := initializers.DB.Unscoped().Delete(&models.Item{}, c.Query("itemId"))
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error deleting item",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "item deleted successfully",
	})
}
