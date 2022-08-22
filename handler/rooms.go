package handler

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jonreesman/chat/database"
	"github.com/jonreesman/chat/model"
	"github.com/jonreesman/chat/room"
)

func CreateRoom(c *fiber.Ctx) error {
	type NewRoom struct {
		Name string
	}

	roomModel := new(model.Room)
	if err := c.BodyParser(roomModel); err != nil {
		fmt.Println(roomModel)
		return c.Status(50).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}
	roomModel.ID = uuid.New()
	if err := database.CreateRoom(roomModel); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't create user", "data": err})
	}
	newRoom := NewRoom{
		Name: roomModel.Name,
	}
	room.CreateRoom(roomModel)
	return c.JSON(fiber.Map{"status": "success", "message": "Created room", "data": newRoom})
}

func DeleteRoom(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "no id selected"})
	}
	roomModel := database.FindRoom(id)
	if roomModel.Name == "" {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "no room by that ID exists"})
	}
	database.DeleteRoom(&roomModel)
	if ok := database.FindRoom(id); ok.Name == "" {
		room.DeleteRoom(id)
		return c.JSON(fiber.Map{"status": "success", "message": "room successfully deleted", "data": nil})
	}
	return c.Status(500).JSON(fiber.Map{"status": "error", "message": "unable to delete room"})
}

func UpdateRoom(c *fiber.Ctx) error {
	id := c.Params("id")
	name := c.Params("name")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "no id selected"})
	}
	room := database.FindRoom(id)
	if room.Name == "" {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "no room by that ID exists"})
	}
	room.Name = name
	database.SaveRoom(&room)
	room = database.FindRoom(id)
	if room.Name == name {
		return c.JSON(fiber.Map{"status": "success", "message": "room successfully updated", "data": nil})
	}
	return c.Status(500).JSON(fiber.Map{"status": "error", "message": "unable to update room"})

}

func GetRooms(c *fiber.Ctx, rooms map[string]*room.Room) error {
	type hubPayLoad struct {
		ID   string
		Name string
	}

	jsonHubs := make([]hubPayLoad, 0)

	for key, val := range rooms {
		fmt.Printf("Key: %s, val: %s\n", key, val.Name)
		jsonHubs = append(jsonHubs, hubPayLoad{
			ID:   key,
			Name: val.Name,
		})
	}

	res, err := json.Marshal(jsonHubs)
	if err != nil {
		log.Printf("Failed to marshall room list: %v", err)
		return err
	}
	return c.Send(res)
}
