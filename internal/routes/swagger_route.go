package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func InitializeSwaggerRoute(a *fiber.App) {
	a.Get("/swagger/*", swagger.HandlerDefault)
}
