package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/naqeeb9a/coffee-shop-apis/initializers"
	"github.com/naqeeb9a/coffee-shop-apis/models"
)

func AllPaymentCards(c *gin.Context) {
	data, usererr := c.Get("user")
	if !usererr {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid token",
			"details": "No user assigned to this token please login again and pass the new token",
		})
		return
	}
	var loggedInUser models.User = data.(models.User)
	var payemntCard []models.PaymentCard
	result := initializers.DB.Find(&payemntCard, "user_id = ?", loggedInUser.ID)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Error fetching cards",
			"reason": result.Error,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cardsCount": result.RowsAffected,
		"cards":      payemntCard,
	})
}

func AddPaymentCards(c *gin.Context) {
	var body struct {
		CardNumber  string
		Cvc         int
		CardExpDate string
		Name        string
	}

	if (c.Bind(&body)) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to read body",
			"details": c.Errors,
		})
		return
	}
	var validationError = make(map[string]string)
	if body.CardNumber == "" {
		validationError["CardNumber"] = "Please provide a Card Number"
	}
	if body.Cvc == 0 {
		validationError["Cvc"] = "Please provide a CVC"
	}
	if body.Name == "" {
		validationError["Name"] = "Please provide a Name"
	}
	if body.CardExpDate == "" {
		validationError["CardExpDate"] = "Please provide a card expiry date"
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
	payemntCard := models.PaymentCard{Name: body.Name, UserId: user.ID, CardExpDate: body.CardExpDate, Cvc: body.Cvc, CardNumber: body.CardNumber}

	result := initializers.DB.Create(&payemntCard)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error while adding payement card",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "payment card added successfully",
		"details": payemntCard,
	})
}
func RemovePaymentCard(c *gin.Context) {
	if c.Query("cardId") == "" {
		c.AbortWithStatusJSON(401, gin.H{
			"error":           "Query parameter (cardId) was not passed",
			"passed (cardId)": c.Query("cardId"),
		})
		return
	}

	result := initializers.DB.Unscoped().Delete(&models.PaymentCard{}, c.Query("cardId"))
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error deleting Card",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "Card deleted successfully",
	})
}
