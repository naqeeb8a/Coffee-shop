package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/coffee-shop5860322/api/initializers"
	"gitlab.com/coffee-shop5860322/api/models"
)

func AddAddress(c *gin.Context) {
	var body struct {
		AddressLine1 string
		AddressLine2 string
		City         string
		State        string
		PostalCode   int
		Country      string
	}

	if (c.Bind(&body)) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	var validationError = make(map[string]string)
	if body.AddressLine1 == "" {
		validationError["AddressLine1"] = "Please provide a AddressLine1"
	}
	if body.City == "" {
		validationError["City"] = "Please provide a City"
	}
	if body.State == "" {
		validationError["State"] = "Please provide a State"
	}
	if body.PostalCode == 0 {
		validationError["PostalCode"] = "Please provide a PostalCode"
	}
	if body.Country == "" {
		validationError["Country"] = "Please provide a Country"
	}
	if len(validationError) != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":    "Fill all fields",
			"required": validationError,
		})
		return
	}
	data, err := c.Get("user")
	if !err {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid token",
			"details": "No user assigned to this token please login again and pass the new token",
		})
		return
	}
	var user models.User = data.(models.User)
	address := models.Address{AddressLine1: body.AddressLine1, AddressLine2: body.AddressLine2, City: body.City, State: body.State, PostalCode: body.PostalCode, Country: body.Country, UserId: user.ID}

	result := initializers.DB.Create(&address)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error creating address",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "address added successfully",
		"details": address,
	})
}
func AllAddresses(c *gin.Context) {
	data, usererr := c.Get("user")
	if !usererr {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid token",
			"details": "No user assigned to this token please login again and pass the new token",
		})
		return
	}
	var loggedInUser models.User = data.(models.User)
	var address []models.Address
	result := initializers.DB.Find(&address, "user_id = ?", loggedInUser.ID)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Error fetch all items of that category",
			"reason": result.Error,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"itemCount": result.RowsAffected,
		"items":     address,
	})
}
func EditAddress(c *gin.Context) {
	var body struct {
		AddressLine1 string
		AddressLine2 string
		City         string
		State        string
		PostalCode   int
		Country      string
	}

	if (c.Bind(&body)) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	if c.Query("addressId") == "" {
		c.AbortWithStatusJSON(401, gin.H{
			"error":              "Query parameter (addressId) was not passed",
			"passed (addressId)": c.Query("addressId"),
		})
		return
	}
	var address models.Address
	initializers.DB.First(&address, c.Query("addressId"))
	if address.ID == 0 {
		c.AbortWithStatusJSON(401, gin.H{
			"error": "No record found of item with this id",
		})
		return
	}
	if body.AddressLine1 != "" {
		address.AddressLine1 = body.AddressLine1
	}
	if body.AddressLine2 != "" {
		address.AddressLine2 = body.AddressLine2
	}
	if body.City != "" {
		address.City = body.City
	}
	if body.State != "" {
		address.State = body.State
	}
	if body.PostalCode != 0 {
		address.PostalCode = body.PostalCode
	}
	if body.Country != "" {
		address.Country = body.Country
	}
	initializers.DB.Save(&address)
	c.JSON(http.StatusOK, gin.H{
		"success": "Address updated successfully",
		"details": address,
	})
}
func RemoveAddress(c *gin.Context) {
	if c.Query("addressId") == "" {
		c.AbortWithStatusJSON(401, gin.H{
			"error":              "Query parameter (addressId) was not passed",
			"passed (addressId)": c.Query("addressId"),
		})
		return
	}

	result := initializers.DB.Unscoped().Delete(&models.Address{}, c.Query("addressId"))
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error deleting address",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "address deleted from favourites",
	})
}
