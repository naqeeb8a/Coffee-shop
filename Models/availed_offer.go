package models

import (
	"gorm.io/gorm"
)

type AvailedOffer struct {
	gorm.Model
	OfferId int
	UserId   uint
}
