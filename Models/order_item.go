package models

import (
	"gorm.io/gorm"
)

type OrderItem struct {
	gorm.Model
	ItemId   uint
	OrderId  uint
	Quantity int
}
