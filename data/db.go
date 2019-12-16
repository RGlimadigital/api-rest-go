package data

import (
	"fmt"

	"github.com/RGlimadigital/Tareas-Go/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func ConnectDB() (*gorm.DB, error) {
	//dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", "localhost", "gorm", "gorm", "gorm123")
	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", "localhost", "rodrigo", "rodrigo", "rodrigo123")
	return gorm.Open("postgres", dbUri)
}

func InitDB() {
	dbCnx, err := ConnectDB()
	if err == nil {
		defer dbCnx.Close()

		dbCnx.AutoMigrate(&models.User{}, &models.Task{}, &models.Image{}, &models.Picture{})
	} else {
		panic(err)
	}
}
