package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/onebluesky882/go_fiber_bun_template/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "go-fiber-bun",
			AppName:      "go-fiber-bun",
		}),

		db: database.New(),
	}

	return server
}
