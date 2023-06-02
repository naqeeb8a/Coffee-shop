package models

import (
	"gorm.io/gorm"
)

type FavouriteItem struct {
	gorm.Model
	ItemId int
	UserId uint
}
