package handler

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/jonreesman/chat/database"
	"github.com/jonreesman/chat/model"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func validToken(t *jwt.Token, id string) bool {
	claims := t.Claims.(jwt.MapClaims)
	uid := claims["user_id"]

	if uid != id {
		return false
	}

	return true
}

func validClient(id string, p string) bool {
	db := database.DB
	var user model.Client
	db.First(&user, id)
	if user.Username == "" {
		return false
	}
	if !CheckPasswordHash(p, user.Password) {
		return false
	}
	return true
}

func GetClientByID(id string) (*model.Client, error) {
	db := database.DB
	var user model.Client
	db.Find(&user, "id = ?", id)
	if user.Username == "" {
		return nil, errors.New("failed to find user with that id")
	}
	return &user, nil
}

func GetClient(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	var user model.Client
	db.Find(&user, "id = ?", id)
	if user.Username == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No user found with ID", "data": nil})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Product found", "data": user})
}

func CreateClient(c *fiber.Ctx) error {
	type NewUser struct {
		Username string
		Password string
	}

	client := new(model.Client)
	if err := c.BodyParser(client); err != nil {
		fmt.Println(client)
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})

	}
	client.ID = uuid.New()
	client.DisplayName = client.Username
	hash, err := hashPassword(client.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't hash password", "data": err})

	}

	client.Password = hash
	if err := database.CreateClient(client); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't create user", "data": err})
	}

	newUser := NewUser{
		Username: client.Username,
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Created user", "data": newUser})
}

// UpdateUser update user
func UpdateClient(c *fiber.Ctx) error {
	type UpdateUserInput struct {
		Names string `json:"names"`
	}
	var uui UpdateUserInput
	if err := c.BodyParser(&uui); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}
	id := c.Params("id")
	token := c.Locals("user").(*jwt.Token)

	if !validToken(token, id) {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Invalid token id", "data": nil})
	}

	user := database.FindClient(id)
	user.DisplayName = uui.Names
	database.SaveClient(&user)
	return c.JSON(fiber.Map{"status": "success", "message": "User successfully updated", "data": user})
}
func DeleteClient(c *fiber.Ctx) error {
	type PasswordInput struct {
		Password string `json:"password"`
	}
	var pi PasswordInput
	if err := c.BodyParser(&pi); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}
	id := c.Params("id")
	token := c.Locals("user").(*jwt.Token)

	if !validToken(token, id) {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Invalid token id", "data": nil})

	}

	if !validClient(id, pi.Password) {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Not valid user", "data": nil})

	}

	client := database.FindClient(id)
	database.DeleteClient(&client)
	return c.JSON(fiber.Map{"status": "success", "message": "User successfully deleted", "data": nil})
}
