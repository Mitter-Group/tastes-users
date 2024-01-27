// user-service/infrastructure/middleware/auth_middleware.go
package middleware

import (
	"github.com/chunnior/users/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware es un middleware global para validar la clave pública.
func AuthMiddleware(encryptedAPIKey string, secretKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		requestEncryptedApiKey := string(c.Request().Header.Peek("X-API-Key"))

		if requestEncryptedApiKey == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Must send an API KEY",
			})
		}
		decryptedRequestAPIKey, err := utils.DecryptAndValidateAPIKey(requestEncryptedApiKey, secretKey)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "API KEY cant be decrypted",
			})
		}

		decryptedAPIKey, err := utils.DecryptAndValidateAPIKey(encryptedAPIKey, secretKey)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "API KEY cant be decrypted",
			})
		}

		if decryptedRequestAPIKey != decryptedAPIKey {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid API KEY",
			})
		}
		return c.Next()
	}
}
