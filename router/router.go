package router

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"
	"github.com/jonreesman/chat/handler"
	"github.com/jonreesman/chat/middleware"
	"github.com/jonreesman/chat/room"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api", logger.New())
	api.Get("/", handler.Hello)

	auth := api.Group("auth")
	auth.Post("/login", handler.Login)

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
		fmt.Println(c)
		handler.ConnectToRoom(c)
	}))

	rooms.Post("/", middleware.Protected(), handler.CreateRoom)
	rooms.Delete("/:id", middleware.Protected(), handler.DeleteRoom)
	rooms.Patch("/:id", middleware.Protected(), handler.UpdateRoom)

	client := api.Group("/client")
	client.Get("/:id", middleware.Protected(), handler.GetClient)
	client.Post("/", handler.CreateClient)
	client.Patch("/:id", middleware.Protected(), handler.UpdateClient)
	client.Delete("/:id", middleware.Protected(), handler.DeleteClient)

}
