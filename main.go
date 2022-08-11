package main

import (
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	app.Use(cors.New())

	router.SetupRoutes(app)

	addr := flag.String("addr", ":8080", "http service address")
	flag.Parse()
	log.Fatal(app.Listen(*addr))
}
