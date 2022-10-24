package handlers

import (
	"balance/internal/databases"
	"balance/internal/models"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"math"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	DB databases.DBInt
}

func NewHandler(DB databases.DBInt) *Handler {
	return &Handler{DB: DB}
}

func returnBadRequest(err error, c *fiber.Ctx) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"message": err.Error(),
	})
}

func (h *Handler) GetBalance(c *fiber.Ctx) error {
	payload := struct {
		ID uint64 `json:"id"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return returnBadRequest(err, c)
	}

	balance, err := h.DB.GetBalance(payload.ID)
	if err != nil {
		return returnBadRequest(err, c)
	}

	return c.JSON(fiber.Map{
		"balance": float32(balance) * 0.01,
	})
}

func (h *Handler) AddBalance(c *fiber.Ctx) error {
	payload := struct {
		ID     uint64  `json:"id"`
		Amount float32 `json:"amount"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return returnBadRequest(err, c)
	}
	err := h.DB.AddBalance(payload.ID, int64(math.Ceil(float64(payload.Amount*100))))
	if err != nil {
		return returnBadRequest(err, c)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) DeleteUser(c *fiber.Ctx) error {
	payload := struct {
		ID uint64 `json:"id"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return returnBadRequest(err, c)
	}

	err := h.DB.DeleteUser(payload.ID)
	if err != nil {
		return returnBadRequest(err, c)
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
		return returnBadRequest(err, c)
	}

	err := h.DB.Reserve(payload.UserID, payload.ServiceID, payload.OrderID, int64(math.Ceil(float64(payload.Amount*100))))
	if err != nil {
		return returnBadRequest(err, c)
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
		return returnBadRequest(err, c)
	}

	reserve, err := h.DB.GetReserve(payload.UserID, payload.ServiceID, payload.OrderID)
	if err != nil {
		return returnBadRequest(err, c)
	}
	outPayload := struct {
		OrderID     uint64         `json:"order_id" gorm:"primaryKey;autoIncrement:false"`
		UserID      uint64         `json:"-"`
		User        models.User    `json:"user"`
		ServiceID   uint64         `json:"-"`
		Service     models.Service `json:"service"`
		Amount      float32        `json:"amount"`
		Purchased   bool           `json:"purchased"`
		ReservedAt  time.Time      `json:"reserved_at"`
		PurchasedAt *time.Time     `json:"purchased_at,omitempty"` // nullable
	}{
		OrderID:     reserve.OrderID,
		UserID:      reserve.UserID,
		User:        reserve.User,
		ServiceID:   reserve.ServiceID,
		Service:     reserve.Service,
		Amount:      float32(reserve.Amount) * 0.01,
		Purchased:   reserve.Purchased,
		ReservedAt:  reserve.ReservedAt,
		PurchasedAt: reserve.PurchasedAt,
	}
	return c.JSON(outPayload)
}

func (h *Handler) DeleteReserve(c *fiber.Ctx) error {
	payload := struct {
		UserID    uint64  `json:"user_id"`
		ServiceID uint64  `json:"service_id"`
		OrderID   uint64  `json:"order_id"`
		Amount    float32 `json:"amount"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return returnBadRequest(err, c)
	}

	err := h.DB.DeleteReserve(payload.UserID, payload.ServiceID, payload.OrderID, int64(math.Ceil(float64(payload.Amount*100))))
	if err != nil {
		return returnBadRequest(err, c)
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
		return returnBadRequest(err, c)
	}

	err := h.DB.Purchase(payload.UserID, payload.ServiceID, payload.OrderID, int64(math.Ceil(float64(payload.Amount*100))))
	if err != nil {
		return returnBadRequest(err, c)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) AddServices(c *fiber.Ctx) error {
	payload := struct {
		Services []models.Service `json:"services"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return returnBadRequest(err, c)
	}

	err := h.DB.AddServices(payload.Services[:])
	if err != nil {
		return returnBadRequest(err, c)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) GetService(c *fiber.Ctx) error {
	payload := struct {
		ServiceID uint64 `json:"id"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return returnBadRequest(err, c)
	}

	service, err := h.DB.GetService(payload.ServiceID)
	if err != nil {
		return returnBadRequest(err, c)
	}

	return c.JSON(service)
}

func (h *Handler) DeleteService(c *fiber.Ctx) error {
	payload := struct {
		ID uint64 `json:"id"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return returnBadRequest(err, c)
	}

	err := h.DB.DeleteService(payload.ID)
	if err != nil {
		return returnBadRequest(err, c)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) GetReport(c *fiber.Ctx) error {
	payload := struct {
		Year  int `params:"year"`
		Month int `params:"month"`
	}{}
	if err := c.ParamsParser(&payload); err != nil {
		return returnBadRequest(err, c)
	}
	if payload.Month < 1 || payload.Month > 12 {
		return returnBadRequest(errors.New("handler: get report: wrong month input"), c)
	}
	filePath := "./report/" + strconv.Itoa(payload.Year) + "/" + strconv.Itoa(payload.Month) + "/report.csv"
	if _, err := os.Stat(filePath); err == nil {
		return c.SendFile(filePath, false)
	} else if errors.Is(err, os.ErrNotExist) {
		return returnBadRequest(fmt.Errorf("handler: get report: report from %d.%d doesn't exist", payload.Month, payload.Year), c)
	} else {
		return returnBadRequest(err, c)
	}
}

func (h *Handler) CreateReport(c *fiber.Ctx) error {
	payload := struct {
		Year  int `params:"year"`
		Month int `params:"month"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return returnBadRequest(err, c)
	}
	if payload.Month < 1 || payload.Month > 12 {
		return returnBadRequest(errors.New("handler: get report: wrong month input"), c)
	}

	relativeLink, err := h.DB.CreateReport(payload.Year, payload.Month)
	if err != nil {
		return returnBadRequest(err, c)
	}
	link := c.BaseURL() + "/api" + relativeLink

	return c.JSON(fiber.Map{
		"report_link": link,
	})
}
