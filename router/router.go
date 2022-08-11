package router

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"
	"github.com/jonreesman/chat/handler"
	"github.com/jonreesman/chat/room"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api", logger.New())
	api.Get("/", handler.Hello)
	rooms := api.Group("/rooms")
	rooms.Get("/", func(c *fiber.Ctx) error {
		err := handler.GetRooms(c, room.Rooms)
		if err != nil {
			c.Status(500)
			log.Printf("error in getRooms(): %v", err)
			return err
		}
		return nil
	})

	rooms.Get("/:id", websocket.New(func(c *websocket.Conn) {
		handler.ConnectToRoom(c)
	}))
}
