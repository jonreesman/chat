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
	"gorm.io/gorm"
)

// Recieves a user ID from the request parameters
// and returns a copy of the user object from the database
/*
	GET Request Form: http://[ip]:[port]/api/client/:id
	Request Body (JSON): None
	Response Form:
		{
			"status": [response status],
			"message": [success/error],
			"data": [model.Client]
		}
*/
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

// Recieves a new username and password, creates the
// new user in the database, and returns the new user
// object to the creator.
/*
	POST Request Form: http://[ip]:[port]/api/client/
	Request Body (JSON): {
		"Username": [username],
		"Password": [password],
	}
	Response Form:
		{
			"status": [response status],
			"message": [success/error],
			"data": {
				"Username:" [username],
			}
		}
*/
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

// Recieves a new display name from the request body
// and updates the user object in the database.
/*
	PATCH Request Form: http://[ip]:[port]/api/client/:id
	Request Body (JSON): {
		"DisplayName": [new display name]
	}
	Response Form:
		{
			"status": [response status],
			"message": [success/error],
			"data": {
				"Username": [username],
				"DisplayName": [display name],
				"ID": [user id],
			}
		}
*/
func UpdateClient(c *fiber.Ctx) error {
	type UpdateUserInput struct {
		DisplayName string
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

	user, err := GetClientByID(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "internal error", "data": nil})
	}
	user.DisplayName = uui.DisplayName
	database.SaveClient(user)

	//Validate changes
	userCheck, err := GetClientByID(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "internal error", "data": nil})
	}
	if userCheck.DisplayName != uui.DisplayName {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "failed to update display name", "data": nil})
	}

	type updatedUser struct {
		Username    string
		DisplayName string
		ID          uuid.UUID
	}

	uu := updatedUser{
		Username:    user.Username,
		DisplayName: user.DisplayName,
		ID:          user.ID,
	}

	return c.JSON(fiber.Map{"status": "success", "message": "User successfully updated", "data": uu})
}

// Recieves a user ID from the request parameters to delete.
// Validates the delete as a valid operation by checking the
// provided password from the request body as being the password
// of the user to be deleted.
/*
	DELETE Request Form: http://[ip]:[port]/api/client/:id
	Request Body (JSON): {
		"Password": [password]
	}
	Response Form:
		{
			"status": [response status],
			"message": [success/error],
			"data": {
				"Username": [username],
				"DisplayName": [display name],
				"ID": [user id],
			}
		}
*/
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

func GetClientByID(id string) (*model.Client, error) {
	db := database.DB
	var user model.Client
	db.Find(&user, "id = ?", id)
	if user.Username == "" {
		return nil, errors.New("failed to find user with that id")
	}
	return &user, nil
}

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
