package initializers

import "coffee-shop.com/api/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Category{})
	DB.AutoMigrate(&models.Item{})
	DB.AutoMigrate(&models.FavouriteItem{})
	DB.AutoMigrate(&models.Address{})
	DB.AutoMigrate(&models.PaymentCard{})
	DB.AutoMigrate(&models.Offer{})
	DB.AutoMigrate(&models.AvailedOffer{})

}
