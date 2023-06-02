package models

import (
	"gorm.io/gorm"
)

type PaymentCard struct {
	gorm.Model
	UserId      uint
	CardNumber  string
	Cvc         int
	CardExpDate string
	Name        string
}
