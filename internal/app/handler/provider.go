package handler

import (
	"encoding/json"

	"github.com/chunnior/users/internal/domain/provider"
	"github.com/gofiber/fiber/v2"
)

type ProviderHandler struct {
	providerService provider.ProviderService
}

func NewProviderHandler(providerService provider.ProviderService) *ProviderHandler {
	return &ProviderHandler{
		providerService: providerService,
	}
}

func (ph *ProviderHandler) ProviderInfo(c *fiber.Ctx) error {
	providerName, dataType, userId := c.Params("provider"), c.Params("dataType"), c.Params("userId")
	if providerName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Provider is required",
		})
	}
	if dataType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Data type is required",
		})
	}
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}
	request := provider.ProviderInfoRequest{
		Provider: providerName,
		DataType: dataType,
		UserID:   userId,
	}
	response, err := ph.providerService.GetProviderInfo(c.Context(), request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Provider info failed",
		})
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error converting response to JSON",
		})
	}

	return c.Send(responseJSON)
}
