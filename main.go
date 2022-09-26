package main

import (
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jonreesman/chat/config"
	_ "github.com/jonreesman/chat/config"
	"github.com/jonreesman/chat/database"
	"github.com/jonreesman/chat/room"
	"github.com/jonreesman/chat/router"
)

func main() {
	database.Connect()
	app := fiber.New()

	room.RoomSetup()

	app.Static("/", "./home.html")
	app.Static("/avatars", "./uploads/avatars")

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	router.SetupRoutes(app)
	addr := config.GetConfig("HOST")
	flag.Parse()
	log.Fatal(app.Listen(addr))
}
