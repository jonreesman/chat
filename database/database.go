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
	DB.AutoMigrate(&model.Client{}, &model.Message{}, &model.Room{})
}

func CreateClient(client *model.Client) error {
	if err := DB.Create(&client).Error; err != nil {
		return err
	}
	return nil
}

func FindClient(id string) model.Client {
	var client model.Client
	DB.First(&client, id)
	return client
}

func SaveClient(client *model.Client) {
	DB.Save(&client)
}

func DeleteClient(client *model.Client) {
	DB.Delete(&client)
}

func UpdateAvatar(client *model.Client) {
	fmt.Println("Updating avatar: " + client.AvatarURL)
	DB.Model(&model.Client{}).Where("id = ?", client.ID).Update("avatar_url", client.AvatarURL)
}

func FindMessage(id string) model.Message {
	var message model.Message
	DB.First(&message, id)
	return message
}

func SaveMessage(message *model.Message) {
	DB.Save(&message)
}

func GetRoomMessages(id string) []model.Message {
	var messages []model.Message
	DB.Find(&messages, "room_id = ?", id)
	return messages
}

func DeleteMessage(message *model.Message) {
	DB.Delete(&message)
}

func CreateRoom(room *model.Room) error {
	if err := DB.Create(&room).Error; err != nil {
		return err
	}
	return nil
}

func GetRooms() []model.Room {
	var rooms []model.Room
	DB.Find(&rooms)
	return rooms
}

func FindRoom(id string) model.Room {
	var room model.Room
	DB.First(&room, "id = ?", id)
	return room
}

func SaveRoom(room *model.Room) {
	DB.Save(&room)
}

func DeleteRoom(room *model.Room) {
	DB.Unscoped().Where(&model.Message{RoomID: room.ID}).Delete(&model.Message{})
	DB.Delete(&room)
}
