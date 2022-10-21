package routes

import "github.com/gofiber/fiber/v2"

func InitializeRoutes(a *fiber.App) {
	route := a.Group("/api")

	route.Get("")
	route.Post("")
	route.Post("/reserve")
	route.Post("/purchase")
	route.Post("/services")
}
