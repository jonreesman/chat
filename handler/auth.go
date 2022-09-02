package handler

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jonreesman/chat/config"
	"golang.org/x/crypto/bcrypt"
)

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Recieves login info from the body of the HTTP request.
// Creates a new JSON Web Token, with a 72 hour expiry and
// returns it to the client.
/*
	POST Request Form: http://[ip]:[port]/api/auth/login
	Request Body (JSON): {
		"Identity": [username],
		"Password": [user password]
	}
	Response Form:
		{
			"status": [response status],
			"message": [success/error],
			"token": [JWT],
			"user": [model.Client]
		}
*/

func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Identity string
		Password string
	}
	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on login request", "data": err})
	}

	identity := input.Identity
	password := input.Password

	user, err := getClientByUsername(identity)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Error on username", "data": err})
	}
	if !CheckPasswordHash(password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid password", "data": nil})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["user_id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(config.GetConfig("SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    t,
		Expires:  time.Now().Add(time.Hour * 72),
		HTTPOnly: true,
		SameSite: "Lax",
	})
	return c.JSON(fiber.Map{"status": "success", "message": "Success login", "token": t, "user": user})
}

// Returns a JWT for the specific session of a
// user connecting to a given rooms websocket
/*
	POST Request Form: http://[ip]:[port]/api/rooms/room
	Query Params (JSON): {
		"user_id": [user ID],
		"room_id": [room ID]
	}
	Response Form:
		{
			"status": [response status],
			"message": [success/error],
			"token": [JWT],
		}
*/
func GetRoomToken(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	fmt.Println(token)
	userID := c.Query("user_id")
	roomID := c.Query("room_id")
	if !validToken(token, userID) {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Invalid token id", "data": nil})
	}

	userToken, err := token.SignedString([]byte(config.GetConfig("SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	roomToken := jwt.New(jwt.SigningMethodHS256)
	claims := roomToken.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["room_id"] = roomID
	claims["user_token"] = userToken
	claims["exp"] = time.Now().Add(time.Second * 5).Unix()

	t, err := roomToken.SignedString([]byte(config.GetConfig("SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Room auth established", "room_token": t})
}
