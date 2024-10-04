package handler

import "github.com/gofiber/fiber/v2"

type HealthHandler struct {
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) HealthCheckHandler(c *fiber.Ctx) error {
	return c.SendString("OK")
}
