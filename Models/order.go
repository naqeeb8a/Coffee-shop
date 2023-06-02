package models

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserId         uint
	TotalItemCount int
	TotalPrice     int
	PaymentMethod  string
	AddressId      int
	OrderStatus    string
}
