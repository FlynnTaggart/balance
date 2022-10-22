package handlers

import (
	"balance/internal/databases"
	"balance/internal/models"

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

func (h *Handler) Reserve(c *fiber.Ctx) error {
	payload := struct {
		UserID    uint64  `json:"user_id"`
		ServiceID uint64  `json:"service_id"`
		OrderID   uint64  `json:"order_id"`
		Amount    float32 `json:"amount"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	err := h.DB.Reserve(payload.UserID, payload.ServiceID, payload.OrderID, payload.Amount)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) Purchase(c *fiber.Ctx) error {
	payload := struct {
		UserID    uint64  `json:"user_id"`
		ServiceID uint64  `json:"service_id"`
		OrderID   uint64  `json:"order_id"`
		Amount    float32 `json:"amount"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	err := h.DB.Purchase(payload.UserID, payload.ServiceID, payload.OrderID, payload.Amount)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) AddServices(c *fiber.Ctx) error {
	payload := struct {
		Services []models.Service `json:"services"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	err := h.DB.AddServices(payload.Services[:])
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) GetReserve(c *fiber.Ctx) error {
	payload := struct {
		UserID    uint64 `json:"user_id"`
		ServiceID uint64 `json:"service_id"`
		OrderID   uint64 `json:"order_id"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	reserve, err := h.DB.GetReserve(payload.UserID, payload.ServiceID, payload.OrderID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(reserve)
}

func (h *Handler) GetService(c *fiber.Ctx) error {
	payload := struct {
		ServiceID uint64 `json:"id"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	service, err := h.DB.GetService(payload.ServiceID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(service)
}
