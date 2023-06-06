package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"coffee-shop.com/api/initializers"
	"coffee-shop.com/api/models"
)

func AddFavourite(c *gin.Context) {
	if c.Query("itemId") == "" {
		c.AbortWithStatusJSON(401, gin.H{
			"error":           "Query parameter (itemId) was not passed",
			"passed (itemId)": c.Query("itemId"),
		})
		return
	}

	itemId, err := strconv.Atoi(c.Query("itemId"))
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{
			"error":           "item id should be integer",
			"passed (itemId)": c.Query("itemId"),
		})
		return
	}
	data, usererr := c.Get("user")
	if !usererr {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid token",
			"details": "No user assigned to this token please login again and pass the new token",
		})
		return
	}
	var loggedInUser models.User = data.(models.User)
	item := models.FavouriteItem{ItemId: itemId, UserId: loggedInUser.ID}

	result := initializers.DB.Create(&item)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error creating item",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "Item added to favourites",
	})
}
func RemoveFavourite(c *gin.Context) {
	if c.Query("favouriteId") == "" {
		c.AbortWithStatusJSON(401, gin.H{
			"error":                "Query parameter (favouriteId) was not passed",
			"passed (favouriteId)": c.Query("favouriteId"),
		})
		return
	}

	result := initializers.DB.Unscoped().Delete(&models.FavouriteItem{}, c.Query("favouriteId"))
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error deleting item",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "Item deleted from favourites",
	})
}
func AllFavouriteItems(c *gin.Context) {
	data, usererr := c.Get("user")
	if !usererr {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid token",
			"details": "No user assigned to this token please login again and pass the new token",
		})
		return
	}
	var loggedInUser models.User = data.(models.User)
	var items []models.Item
	result := initializers.DB.Raw("SELECT * FROM `items` where id in(SELECT item_id FROM `favourite_items` where user_id=?);", loggedInUser.ID).Scan(&items)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error fetching favourites",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"itemsCount": result.RowsAffected,
		"items":      items,
	})
}
