package models

import (
	"time"

	"gorm.io/gorm"
)

type Coupon struct {
	gorm.Model
	Discount    float32
	CouponCode string
	Status     string
	StartAt    time.Time
	ExpiryAt   time.Time
}
