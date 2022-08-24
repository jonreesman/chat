package handler

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jonreesman/chat/database"
)

func UploadAvatar(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	headers := c.GetReqHeaders()
	id := headers["Id"]
	fmt.Println(headers)
	if err != nil {
		log.Printf("Avatar save error: %v", err)
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "failed to upload avatar"})
	}
	splitFileName := strings.Split(file.Filename, ".")
	fileEnding := splitFileName[len(splitFileName)-1]
	if err := c.SaveFile(file, fmt.Sprintf("./uploads/avatars/%s.%s", id, fileEnding)); err != nil {
		log.Printf("Avatar save error: %v", err)
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "failed to upload avatar"})
	}
	client, err := GetClientByID(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "invalid client"})
	}
	client.AvatarURL = client.ID.String() + "." + fileEnding
	database.UpdateAvatar(client)
	return c.JSON(fiber.Map{"status": "success", "message": "room successfully updated", "data": client})

}
