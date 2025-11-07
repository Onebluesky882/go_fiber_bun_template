package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
)

func RegisterRoutes(router fiber.Router, db *bun.DB) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)
	router.Get("/", handler.GetAllUsers)
}
