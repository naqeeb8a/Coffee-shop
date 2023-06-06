package controllers

import (
	"net/http"
	"strconv"

	"coffee-shop.com/api/initializers"
	"coffee-shop.com/api/models"
	"github.com/gin-gonic/gin"
)

func AddOfferAvailedUser(c *gin.Context) {
	if c.Query("offerId") == "" {
		c.AbortWithStatusJSON(401, gin.H{
			"error":            "Query parameter (offerId) was not passed",
			"passed (offerId)": c.Query("offerId"),
		})
		return
	}

	OfferId, err := strconv.Atoi(c.Query("offerId"))
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{
			"error":            "points should be integer",
			"passed (offerId)": c.Query("offerId"),
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
	offerAvailed := models.AvailedOffer{OfferId: OfferId, UserId: loggedInUser.ID}

	result := initializers.DB.Create(&offerAvailed)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error availing offer",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "offer availed successfully",
		"details": offerAvailed,
	})

}
