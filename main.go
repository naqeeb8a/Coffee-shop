package main

import (
	"coffee-shop.com/api/controllers"
	"coffee-shop.com/api/initializers"
	"coffee-shop.com/api/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
}

func main() {

	r := gin.Default()
	r.Static("/assets", "./assets")
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.POST("/user/signup", controllers.SignUp)
	r.POST("/user/login", controllers.Login)
	r.PUT("/user/edit", middleware.RequireAuth, controllers.EditUser)
	r.GET("/user/profile", middleware.RequireAuth, controllers.GetUser)
	r.GET("/user/points", middleware.RequireAuth, controllers.AddUserLoyaltyPoints)
	r.GET("/category", middleware.RequireAuth, controllers.AllCategories)
	r.POST("/category/add", middleware.RequireAuth, controllers.AddCategory)
	r.PUT("/category/edit", middleware.RequireAuth, controllers.EditCategory)
	r.DELETE("/category/remove", middleware.RequireAuth, controllers.RemoveCategory)
	r.GET("/category/item", middleware.RequireAuth, controllers.CategoryItem)
	r.GET("/item", middleware.RequireAuth, controllers.AllItems)
	r.POST("/item/add", middleware.RequireAuth, controllers.AddItem)
	r.PUT("/item/edit", middleware.RequireAuth, controllers.EditItem)
	r.DELETE("/item/remove", middleware.RequireAuth, controllers.RemoveItem)
	r.GET("/item/details", middleware.RequireAuth, controllers.ItemDetails)
	r.GET("/favourite", middleware.RequireAuth, controllers.AllFavouriteItems)
	r.POST("/favourite/add", middleware.RequireAuth, controllers.AddFavourite)
	r.DELETE("/favourite/remove", middleware.RequireAuth, controllers.RemoveFavourite)
	r.GET("/address", middleware.RequireAuth, controllers.AllAddresses)
	r.POST("/address/add", middleware.RequireAuth, controllers.AddAddress)
	r.PUT("/address/edit", middleware.RequireAuth, controllers.EditAddress)
	r.DELETE("/address/remove", middleware.RequireAuth, controllers.RemoveAddress)
	r.GET("/paymentCard", middleware.RequireAuth, controllers.AllPaymentCards)
	r.POST("/paymentCard/add", middleware.RequireAuth, controllers.AddPaymentCards)
	r.DELETE("/paymentCard/remove", middleware.RequireAuth, controllers.RemovePaymentCard)
	r.GET("/offer", middleware.RequireAuth, controllers.AllOffers)
	r.POST("/offer/add", middleware.RequireAuth, controllers.AddOffer)
	r.PUT("/offer/edit", middleware.RequireAuth, controllers.Editoffer)
	r.DELETE("/offer/remove", middleware.RequireAuth, controllers.RemoveOffer)
	r.GET("/availedOffer/add", middleware.RequireAuth, controllers.AddOfferAvailedUser)
	r.Run()
}
