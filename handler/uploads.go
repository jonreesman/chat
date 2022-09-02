package handler

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jonreesman/chat/database"
)

// Recieves a user id from the request parameters, the user token,
// and a new avatar image from the request body. Validates the user
// has permission to update the user avatar, saves the image, and
// sets the new user avatar URL to the user object in the database.
/*
	POST Request Form: http://[ip]:[port]/uploads/avatar/:id
	Request Body (JSON): None
	Response Form:
		{
			"status": [response status],
			"message": [success/error],
			"data": {
				"Name": [new room name],
			}
		}
*/
func UploadAvatar(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	id := c.Query("id")
	token := c.Locals("user")
	convertedToken := token.(*jwt.Token)

	if !validToken(convertedToken, id) {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Invalid token id", "data": nil})
	}
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
	user, err := GetClientByID(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "invalid client"})
	}
	user.AvatarURL = user.ID.String() + "." + fileEnding
	database.UpdateAvatar(user)

	type updatedUser struct {
		Username    string
		DisplayName string
		AvatarURL   string
		ID          uuid.UUID
	}

	uu := updatedUser{
		Username:    user.Username,
		DisplayName: user.DisplayName,
		ID:          user.ID,
		AvatarURL:   user.AvatarURL,
	}

	return c.JSON(fiber.Map{"status": "success", "message": "room successfully updated", "data": uu})

}
