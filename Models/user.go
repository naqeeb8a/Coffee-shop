package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ProfilePicture  string
	FirstName       string
	LastName        string
	FullName        string
	Mobile          string
	Email           string `gorm:"unique"`
	Password        string
	AccessToken     string
	ExpiryAt        int64
	AppVersion      string
	DeviceOsVersion string
	DeviceModel     string
	DeviceUTCOffSet string
	DeviceToken     string
	LoyaltyPoints   int
}
