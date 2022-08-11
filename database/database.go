package database

import (
	"fmt"

	"github.com/jonreesman/chat/config"
	"github.com/jonreesman/chat/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error
	port := config.GetConfig("DB_PORT")
	host := config.GetConfig("DB_HOST")
	user := config.GetConfig("DB_USER")
	password := config.GetConfig("DB_PASSWORD")
	name := config.GetConfig("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, name)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(fmt.Sprintf("Failed to connect to databse: %v", err))
	}
	DB.AutoMigrate(&model.Client{})
}
