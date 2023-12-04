package handler

import (
	"encoding/json"
	"fmt"

	"github.com/chunnior/users/internal/domain/callback"
	"github.com/chunnior/users/internal/domain/login"
	"github.com/gofiber/fiber/v2"
)

type LoginHandler struct {
	loginService login.LoginService
}

func NewLoginHandler(loginService login.LoginService) *LoginHandler {
	return &LoginHandler{
		loginService: loginService,
	}
}

func (h *LoginHandler) Login(c *fiber.Ctx) error {
	var request login.LoginRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	response, err := h.loginService.Login(c.Context(), request)
	if err != nil {
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Callback failed",
		})
	}

	return c.JSON(response)
	/*
		// Acceder a los valores del cuerpo de la solicitud
		provider := callbackRequest.Provider
		code := callbackRequest.Code
		state := callbackRequest.State
		token := callbackRequest.Token

		// Realizar acciones con los valores recibidos
		fmt.Printf("Provider: %s\n", provider)

		// Acceder a propiedades específicas según el proveedor
		switch provider {
		case "spotify":
			// Realizar acciones específicas para Spotify
			fmt.Printf("Code: %s\n", code)
			fmt.Printf("State: %s\n", state)

		case "google":
			// Acciones específicas para Google
			fmt.Printf("Code: %s\n", code)
			fmt.Printf("State: %s\n", state)
			fmt.Printf("Token: %s\n", token)

		case "instagram":
			// Acciones específicas para Instagram
			fmt.Printf("Code: %s\n", code)

		default:
			// Acciones para otros proveedores
		}

		return c.SendString("Procesamiento del callback completado")
	*/
	/*
		var callbackRequest callback.CallbackRequest

		// Decodificar el cuerpo JSON de la solicitud manualmente
		if err := json.Unmarshal(c.Body(), &callbackRequest); err != nil {
			fmt.Println("Error al decodificar el cuerpo JSON:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Error al decodificar el cuerpo JSON"})
		}

		// Acceder a los valores del cuerpo de la solicitud
		provider := callbackRequest.Provider
		code := callbackRequest.Code
		state := callbackRequest.State

		// Realizar acciones con los valores recibidos
		fmt.Printf("Provider: %s\n", provider)
		fmt.Printf("Code: %s\n", code)
		fmt.Printf("State: %s\n", state)

		return c.SendString("Procesamiento del callback completado")
	*/
	/*
		var callbackRequest callback.CallbackRequest

		// Analizar el cuerpo JSON de la solicitud y almacenar los valores en la instancia
		if err := c.BodyParser(&callbackRequest); err != nil {
			fmt.Println("Error al analizar el cuerpo JSON:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Error al analizar el cuerpo JSON"})
		}

		// Acceder a los valores a través de la interfaz
		provider := callbackRequest.GetProvider()
		code := callbackRequest.GetCode()
		state := callbackRequest.GetState()

		// Realizar acciones con los valores recibidos
		fmt.Printf("Provider: %s\n", provider)
		fmt.Printf("Code: %s\n", code)
		fmt.Printf("State: %s\n", state)

		return c.SendString("Procesamiento del callback completado")
	*/
	/*
		var requestBody map[string]interface{}

		// Analizar el cuerpo JSON de la solicitud y almacenar los valores en el mapa
		if err := c.BodyParser(&requestBody); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Error al analizar el cuerpo JSON"})
		}

		// Imprimir todos los valores del cuerpo de la solicitud
		fmt.Println("Request Body:")
		for key, value := range requestBody {
			fmt.Printf("%s: %v\n", key, value)
		}

		return c.SendString("Procesamiento de parámetros POST completado")
	*/
	/*
		var request callback.CallbackRequest

		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse JSON",
			})
		}

		response, err := h.loginService.Callback(c.Context(), request)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Callback failed",
			})
		}

		return c.JSON(response)
	*/
}
