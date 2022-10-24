package routes

import (
	"balance/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func InitializeRoutes(a *fiber.App, handler *handlers.Handler) {
	route := a.Group("/api")

	route.Get("", handler.GetBalance)
	route.Post("", handler.AddBalance)
	route.Post("/reserve", handler.Reserve)
	route.Get("/reserve", handler.GetReserve)
	route.Delete("/reserve", handler.DeleteReserve)
	route.Post("/purchase", handler.Purchase)
	route.Post("/services", handler.AddServices)
	route.Get("/services", handler.GetService)
}
