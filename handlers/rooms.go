package handlers

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jonreesman/chat/room"
)

func GetRooms(c *fiber.Ctx, rooms map[string]*room.Room) error {
	type hubPayLoad struct {
		ID   string
		Name string
	}

	jsonHubs := make([]hubPayLoad, 0)

	for key, val := range rooms {
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
