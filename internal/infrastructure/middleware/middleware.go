// user-service/infrastructure/middleware/auth_middleware.go
package middleware

import (
	"fmt"

	"github.com/chunnior/users/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware es un middleware global para validar la clave pública.
func AuthMiddleware(encryptedAPIKey string, secretKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		allHeaders := c.Request().Header.String()
		fmt.Println("All Headers:", allHeaders)

		requestEncryptedApiKey := string(c.Request().Header.Peek("X-API-Key"))
		fmt.Println("API Key:", requestEncryptedApiKey)
		decryptedRequestAPIKey, err := utils.DecryptAndValidateAPIKey(requestEncryptedApiKey, secretKey)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid api key 1",
			})
		}

		decryptedAPIKey, err := utils.DecryptAndValidateAPIKey(encryptedAPIKey, secretKey)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid api key 2",
			})
		}

		if decryptedRequestAPIKey != decryptedAPIKey {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid api key 3",
			})
		}
		return c.Next()
	}
}
