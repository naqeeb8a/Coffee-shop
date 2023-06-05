package initializers

import "gitlab.com/coffee-shop5860322/api/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Category{})
	DB.AutoMigrate(&models.Item{})
	DB.AutoMigrate(&models.FavouriteItem{})
	DB.AutoMigrate(&models.Address{})
	DB.AutoMigrate(&models.PaymentCard{})
}
