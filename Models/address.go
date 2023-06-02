package models

import (
	"gorm.io/gorm"
)

type Address struct {
	gorm.Model
	UserId       uint
	AddressLine1 string
	AddressLine2 string
	City         string
	State        string
	PostalCode   int
	Country      string
}
