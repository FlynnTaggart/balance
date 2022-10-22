package handlers

import (
	"balance/internal/databases"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	DB databases.DBInt
}

func NewHandler(DB databases.DBInt) *Handler {
	return &Handler{DB: DB}
}

func (h *Handler) GetBalance(c *fiber.Ctx) error {
	payload := struct {
		ID uint64 `json:"id"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	balance, err := h.DB.GetBalance(payload.ID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"balance": balance,
	})
}

func (h *Handler) AddBalance(c *fiber.Ctx) error {
	payload := struct {
		ID     uint64  `json:"id"`
		Amount float32 `json:"amount"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	err := h.DB.AddBalance(payload.ID, payload.Amount)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}
