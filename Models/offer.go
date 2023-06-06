package models

import (
	"time"

	"gorm.io/gorm"
)

type Offer struct {
	gorm.Model
	Discount    float32
	OfferCode string
	Status     string
	StartAt    time.Time
	ExpiryAt   time.Time
}
