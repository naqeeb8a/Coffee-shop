package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"coffee-shop.com/api/initializers"
	"coffee-shop.com/api/models"
)

func AddOffer(c *gin.Context) {
	var body struct {
		Discount  float32
		OfferCode string
		Status    string
		StartAt   string
		ExpiryAt  string
	}

	if (c.Bind(&body)) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	var validationError = make(map[string]string)
	if body.Discount == 0 {
		validationError["Discount"] = "Please provide a Discount"
	}
	if body.OfferCode == "" {
		validationError["OfferCode"] = "Please provide a offer Code"
	}
	if body.Status == "" {
		validationError["Status"] = "Please provide a status"
	}
	if body.StartAt == "" {
		validationError["StartAt"] = "Please provide a starting date"
	}
	if body.ExpiryAt == "" {
		validationError["ExpiryAt"] = "Please provide an expiry date"
	}
	const layout = "2006-Jan-02"
	startAt, err := time.Parse(layout, body.StartAt)
	if err != nil {
		validationError["startAt"] = "Error parsing start date please use this layout (2006-Jan-02)"
	}
	ExpiryAt, err2 := time.Parse(layout, body.ExpiryAt)
	if err2 != nil {
		validationError["ExpiryAt"] = "Error parsing expiry date please use this layout (2006-Jan-02)"
	}
	if len(validationError) != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":    "Fill all fields",
			"required": validationError,
		})
		return
	}

	offer := models.Offer{Discount: body.Discount, OfferCode: body.OfferCode, Status: body.Status, StartAt: startAt, ExpiryAt: ExpiryAt}

	result := initializers.DB.Create(&offer)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error creating offer",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "Offer created successfully",
		"details": offer,
	})
}

func Editoffer(c *gin.Context) {
	if c.Query("offerId") == "" {
		c.AbortWithStatusJSON(401, gin.H{
			"error":            "Query parameter (offerId) was not passed",
			"passed (offerId)": c.Query("offerId"),
		})
		return
	}
	var body struct {
		Discount  float32
		OfferCode string
		Status    string
		StartAt   string
		ExpiryAt  string
	}

	if (c.Bind(&body)) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	var validationError = make(map[string]string)
	const layout = "2006-Jan-02"
	var startAt time.Time
	var err error
	var ExpiryAt time.Time
	var err2 error

	if body.StartAt != "" {
		startAt, err = time.Parse(layout, body.StartAt)
		if err != nil {
			validationError["startAt"] = "Error parsing start date please use this layout (2006-Jan-02)"
		}
	}
	if body.ExpiryAt != "" {
		ExpiryAt, err2 = time.Parse(layout, body.ExpiryAt)
		if err2 != nil {
			validationError["ExpiryAt"] = "Error parsing expiry date please use this layout (2006-Jan-02)"
		}
	}

	if len(validationError) != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":    "Fill all fields",
			"required": validationError,
		})
		return
	}

	var offer models.Offer
	initializers.DB.First(&offer, c.Query("offerId"))
	if offer.ID == 0 {
		c.AbortWithStatusJSON(401, gin.H{
			"error": "No record found of user with this id",
		})
		return
	}
	if body.Discount != 0 {
		offer.Discount = body.Discount
	}
	if body.OfferCode != "" {
		offer.OfferCode = body.OfferCode
	}
	if body.Status != "" {
		offer.Status = body.Status
	}
	if body.StartAt != "" {
		offer.StartAt = startAt
	}
	if body.ExpiryAt != "" {
		offer.ExpiryAt = ExpiryAt
	}
	initializers.DB.Save(&offer)
	c.JSON(http.StatusOK, gin.H{
		"success": "offer updated successfully",
		"details": offer,
	})
}
func RemoveOffer(c *gin.Context) {
	if c.Query("offerId") == "" {
		c.AbortWithStatusJSON(401, gin.H{
			"error":            "Query parameter (offerId) was not passed",
			"passed (offerId)": c.Query("offerId"),
		})
		return
	}

	result := initializers.DB.Unscoped().Delete(&models.Offer{}, c.Query("offerId"))
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error deleting offer",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "offer deleted from favourites",
	})
}
func AllOffers(c *gin.Context) {
	var offers []models.Offer
	result := initializers.DB.Find(&offers)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error fetching offers",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"offersCount": result.RowsAffected,
		"offers":      offers,
	})
}
