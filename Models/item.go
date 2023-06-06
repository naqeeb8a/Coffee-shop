package models

import (
	"gorm.io/gorm"
)

type Item struct {
	gorm.Model
	Name        string `gorm:"unique"`
	Image       string
	Description string
	Price       int
	IsEnabled   bool
	CategoryId  int
	Rating      int16
}
