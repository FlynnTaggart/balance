package routes

import (
	"balance/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func InitializeRoutes(a *fiber.App, handler *handlers.Handler) {
	route := a.Group("/api")

	route.Get("", handler.GetBalance)
	route.Post("", handler.AddBalance)
	//route.Post("/reserve")
	//route.Post("/purchase")
	//route.Post("/services")
}
