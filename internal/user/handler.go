package user

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.service.GetAllUsers(context.Background())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"status": "success",
		"data":   users,
	})
}
