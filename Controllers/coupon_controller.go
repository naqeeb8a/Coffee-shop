package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/coffee-shop5860322/api/initializers"
	"gitlab.com/coffee-shop5860322/api/models"
)

func CreateCoupon(c *gin.Context) {
	var body struct {
		Discount   float32
		CouponCode string
		Status     string
		StartAt    string
		ExpiryAt   string
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
	if body.CouponCode == "" {
		validationError["CouponCode"] = "Please provide a Coupon Code"
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

	coupon := models.Coupon{Discount: body.Discount, CouponCode: body.CouponCode, Status: body.Status, StartAt: startAt, ExpiryAt: ExpiryAt}

	result := initializers.DB.Create(&coupon)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error creating coupon",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "Coupon created successfully",
		"details": coupon,
	})
}
func AllCoupons(c *gin.Context) {
	var coupons []models.Coupon
	result := initializers.DB.Find(&coupons)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error fetching coupons",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"itemsCount": result.RowsAffected,
		"items":      coupons,
	})
}
