package handler

import (
	"encoding/json"
	"fmt"

	"github.com/chunnior/users/internal/domain"
	"github.com/chunnior/users/internal/domain/callback"
	"github.com/chunnior/users/internal/domain/login"
	"github.com/gofiber/fiber/v2"
)

type LoginHandler struct {
	loginService login.LoginService
	logger       domain.Logger
}

func NewLoginHandler(loginService login.LoginService, logger domain.Logger) *LoginHandler {
	return &LoginHandler{
		loginService: loginService,
		logger:       logger,
	}
}

func (h *LoginHandler) Login(c *fiber.Ctx) error {
	var request login.LoginRequest

	if err := c.BodyParser(&request); err != nil {
		h.logger.Error(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	response, err := h.loginService.Login(c.Context(), request)
	if err != nil {
		h.logger.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Login failed",
		})
	}

	return c.JSON(response)
}

// ProviderProperties define las propiedades comunes para todos los proveedores
type ProviderProperties interface {
	// Agrega propiedades comunes aquí
	GetCode() string
	GetState() string
}

func (h *LoginHandler) Callback(c *fiber.Ctx) error {
	// Crear una instancia de la estructura que coincide con la estructura del cuerpo JSON
	var callbackRequest callback.CallbackRequestBody

	// Decodificar el cuerpo JSON de la solicitud manualmente
	if err := json.Unmarshal(c.Body(), &callbackRequest); err != nil {
		fmt.Println("Error al decodificar el cuerpo JSON:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Error al decodificar el cuerpo JSON"})
	}

	response, err := h.loginService.Callback(c.Context(), callbackRequest)
	if err != nil {
		h.logger.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Callback failed",
		})
	}

	return c.JSON(response)
}
