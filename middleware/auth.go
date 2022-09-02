package middleware

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jonreesman/chat/config"
)

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}

func Protected() fiber.Handler {
	log.Printf("User attempted protected route.")
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(config.GetConfig("SECRET")),
		ErrorHandler: jwtError,
		TokenLookup:  "cookie:token",
	})
}

func GetToken(c *fiber.Ctx) error {
	token := c.Cookies("token")
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetConfig("SECRET")), nil
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Invalid token id", "data": nil})
	}
	fmt.Println("token: " + token)
	fmt.Println("user")
	fmt.Println(c.Locals("user"))
	c.Locals("user", parsedToken)
	return c.Next()
}

func ParseToken(token string) (*jwt.Token, error) {
	token = strings.ReplaceAll(token, "Bearer ", "")
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.GetConfig("SECRET")), nil
	})
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		fmt.Println(claims["user_id"])
		return parsedToken, nil
	} else {
		return nil, err
	}
}
