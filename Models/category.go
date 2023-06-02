package models

import (
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	CategoryImage string
	Name          string `gorm:"unique"`
}
