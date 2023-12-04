package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// AuthRequired is a middleware to check if the user is authenticated
func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Here you can add your authentication logic
		// For example, you can check if a session exists

		// If the user is not authenticated, return an error
		// return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})

		// If the user is authenticated, call the next handler
		return c.Next()
	}
}