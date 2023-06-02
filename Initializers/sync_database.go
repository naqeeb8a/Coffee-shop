package initializers

import "github.com/naqeeb9a/coffee-shop-apis/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Category{})
	DB.AutoMigrate(&models.Item{})
	DB.AutoMigrate(&models.FavouriteItem{})
	DB.AutoMigrate(&models.Address{})
	DB.AutoMigrate(&models.PaymentCard{})
}
