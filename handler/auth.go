package handler

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/jonreesman/chat/config"
	"github.com/jonreesman/chat/database"
	"github.com/jonreesman/chat/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func getClientByUsername(name string) (*model.Client, error) {
	fmt.Println(name)
	db := database.DB
	var client model.Client
	if err := db.Where(&model.Client{Username: name}).Find(&client).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &client, nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Identity string
		Password string
	}
	type ClientData struct {
		ID       uuid.UUID
		Username string
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

	return c.JSON(fiber.Map{"status": "success", "message": "Success login", "token": t, "user": user})
}
