package initializers

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	dsn := os.Getenv("DB")
	d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to DB")
	}
	DB = d

}

func GetDB() *gorm.DB {
	return DB
}
