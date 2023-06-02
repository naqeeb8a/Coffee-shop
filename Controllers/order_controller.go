package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/naqeeb9a/coffee-shop-apis/initializers"
	"github.com/naqeeb9a/coffee-shop-apis/models"
)

func CreateOrder(c *gin.Context) {
	var body struct {
		TotalItemCount int
		TotalPrice     int
		PaymentMethod  string
		AddressId      int
		OrderStatus    string
		Items          []models.OrderItem
	}

	if (c.Bind(&body)) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	var validationError = make(map[string]string)
	if body.TotalItemCount == 0 {
		validationError["TotalItemCount"] = "Please provide a total item count"
	}
	if body.TotalPrice == 0 {
		validationError["TotalPrice"] = "Please provide an total price"
	}
	if body.AddressId == 0 {
		validationError["AddressId"] = "Please provide an address Id"
	}
	if body.OrderStatus == "" {
		validationError["OrderStatus"] = "Please provide an order status"
	}
	if body.PaymentMethod == "" {
		validationError["PaymentMethod"] = "Please provide an payment method"
	}
	if len(validationError) != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":    "Fill all fields",
			"required": validationError,
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
	order := models.Order{UserId: loggedInUser.ID, TotalItemCount: body.TotalItemCount, TotalPrice: body.TotalPrice, AddressId: body.AddressId, OrderStatus: body.OrderStatus, PaymentMethod: body.PaymentMethod}

	result := initializers.DB.Create(&order)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error placing Order",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "Order placed successfully",
	})
}
