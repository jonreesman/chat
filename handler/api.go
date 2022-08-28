package handler

import "github.com/gofiber/fiber/v2"

// Recieves login info from the body of the HTTP request.
// Creates a new JSON Web Token, with a 72 hour expiry and
// returns it to the client.
/*
	GET Request Form: http://[ip]:[port]/
	Request Body (JSON): None
	Response Form:
		{
			"status": [response status],
			"message": [success/error],
			"data": nil
		}
*/
func Hello(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "success", "message": "Hello i'm ok!", "data": nil})
}
