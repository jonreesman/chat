package main

import (
	"flag"
	"log"
	"os"

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

	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.GetConfig("ORIGIN"),
		AllowCredentials: true,
	}))

	app.Static("/", "./home.html")
	if _, err := os.Stat("./uploads/avatars"); os.IsNotExist(err) {
		err := os.MkdirAll("./uploads/avatars", os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
	app.Static("/avatars", "./uploads/avatars")

	router.SetupRoutes(app)
	addr := config.GetConfig("HOST")
	flag.Parse()
	log.Fatal(app.Listen(addr))
}
