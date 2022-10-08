package router

import (
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
	auth.Get("/logout", middleware.Protected(), handler.Logout)

	rooms := api.Group("/rooms")
	rooms.Get("/", middleware.Protected(), func(c *fiber.Ctx) error {
		err := handler.GetRooms(c, room.Rooms)
		if err != nil {
			c.Status(500)
			log.Printf("error in getRooms(): %v", err)
			return err
		}
		return nil
	})

	rooms.Get("/room", middleware.Protected(), handler.GetRoomToken)

	rooms.Use("/room/:id", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	rooms.Get("/room/:id", websocket.New(func(c *websocket.Conn) {
		handler.ConnectToRoom(c)
	}))

	rooms.Post("/", middleware.Protected(), handler.CreateRoom)
	rooms.Delete("/:id", middleware.Protected(), handler.DeleteRoom)
	rooms.Patch("/:id", middleware.Protected(), handler.UpdateRoom)

	client := api.Group("/client")
	client.Get("/", middleware.Protected(), handler.GetClient)
	client.Get("/:id", middleware.Protected(), middleware.GetToken, handler.GetClient)
	client.Post("/", middleware.DemoAuth, handler.CreateClient)
	client.Patch("/:id", middleware.Protected(), middleware.GetToken, handler.UpdateClient)
	client.Delete("/:id", middleware.Protected(), handler.DeleteClient)

	uploads := app.Group("/uploads")
	uploads.Post("/avatar", middleware.Protected(), middleware.GetToken, handler.UploadAvatar)

}
