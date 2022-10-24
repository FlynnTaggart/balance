package handlers

import (
	"balance/internal/databases"
	"balance/internal/models"
	"balance/internal/utils"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
)

type Handler struct {
	DB databases.DBInt
}

// NewHandler creates new Handler instance
func NewHandler(DB databases.DBInt) *Handler {
	return &Handler{DB: DB}
}

// returnBadRequest wraps bad request
func returnBadRequest(err error, c *fiber.Ctx) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"message": err.Error(),
	})
}

// GetBalance gets user balance by id
// @Description Get user balance by given id
// @Summary     Get user balance
// @Tags        Balance
// @Accept      json
// @Produce     json
// @Param       inJSON body     models.PayloadId      true "In JSON with User ID"
// @Success     200    {object} models.PayloadBalance "User's balance"
// @Failure     400    {object} models.PayloadErr     "Error"
// @Router      / [get]
func (h *Handler) GetBalance(c *fiber.Ctx) error {
	payload := models.PayloadId{}
	if err := c.BodyParser(&payload); err != nil {
		return returnBadRequest(err, c)
	}

	balance, err := h.DB.GetBalance(payload.ID)
	if err != nil {
		return returnBadRequest(err, c)
	}

	return c.JSON(fiber.Map{
		"balance": utils.MoneyToFloat(balance),
	})
}

// AddBalance adds amount money to user by id
// @Description Add user balance by given id
// @Summary     Add user balance
// @Tags        Balance
// @Accept      json
// @Produce     json
// @Param       inJSON body     models.PayloadAddBalance true "In JSON with User ID and Amount"
// @Success     200    {string} status                   "OK"
// @Failure     400    {object} models.PayloadErr        "Error"
// @Router      / [post]
func (h *Handler) AddBalance(c *fiber.Ctx) error {
	payload := models.PayloadAddBalance{}
	if err := c.BodyParser(&payload); err != nil {
		return returnBadRequest(err, c)
	}
	err := h.DB.AddBalance(payload.ID, utils.MoneyToInt(payload.Amount))
	if err != nil {
		return returnBadRequest(err, c)
	}

	return c.SendStatus(fiber.StatusOK)
}

// DeleteUser deletes user id
// @Description Delete user by given id
// @Summary     Delete user
// @Tags        Users
// @Accept      json
// @Produce     json
// @Param       inJSON body     models.PayloadId  true "In JSON with User ID"
// @Success     200    {string} status            "OK"
// @Failure     400    {object} models.PayloadErr "Error"
// @Router      /users/ [delete]
func (h *Handler) DeleteUser(c *fiber.Ctx) error {
	payload := models.PayloadId{}
	if err := c.BodyParser(&payload); err != nil {
		return returnBadRequest(err, c)
	}

	err := h.DB.DeleteUser(payload.ID)
	if err != nil {
		return returnBadRequest(err, c)
	}

	return c.SendStatus(fiber.StatusOK)
}

// Reserve performs money reserve transaction for given orderId, userId, serviceId and amount.
// @Description Reserve money for given orderId, userId, serviceId and amount.
// @Summary     Reserve money
// @Tags        Reserves
// @Accept      json
// @Produce     json
// @Param       inJSON body     models.PayloadReserve true "In JSON with user_id, service_id, order_id and amount"
// @Success     200    {string} status                "OK"
// @Failure     400    {object} models.PayloadErr     "Error"
// @Router      /reserve/ [post]
func (h *Handler) Reserve(c *fiber.Ctx) error {
	payload := models.PayloadReserve{}
	if err := c.BodyParser(&payload); err != nil {
		return returnBadRequest(err, c)
	}

	err := h.DB.Reserve(payload.UserID, payload.ServiceID, payload.OrderID, utils.MoneyToInt(payload.Amount))
	if err != nil {
		return returnBadRequest(err, c)
	}

	return c.SendStatus(fiber.StatusOK)
}

// GetReserve gets reserve by user_id, service_id, order_id
// @Description Get reserve by user_id, service_id, order_id
// @Summary     Get reserve
// @Tags        Reserves
// @Accept      json
// @Produce     json
// @Param       inJSON body     models.PayloadReserve true "In JSON with user_id, service_id, order_id"
// @Success     200    {object} models.Reserve        "Service"
// @Failure     400    {object} models.PayloadErr     "Error"
// @Router      /reserve/ [get]
func (h *Handler) GetReserve(c *fiber.Ctx) error {
	payload := models.PayloadReserve{}
	if err := c.BodyParser(&payload); err != nil {
		return returnBadRequest(err, c)
	}

	reserve, err := h.DB.GetReserve(payload.UserID, payload.ServiceID, payload.OrderID)
	if err != nil {
		return returnBadRequest(err, c)
	}
	outPayload := models.ReserveFloatAmount{
		OrderID:     reserve.OrderID,
		UserID:      reserve.UserID,
		User:        reserve.User,
		ServiceID:   reserve.ServiceID,
		Service:     reserve.Service,
		Amount:      utils.MoneyToFloat(reserve.Amount),
		Purchased:   reserve.Purchased,
		ReservedAt:  reserve.ReservedAt,
		PurchasedAt: reserve.PurchasedAt,
	}
	return c.JSON(outPayload)
}

// DeleteReserve remove reserve for given orderId, userId, serviceId and amount.
// @Description Remove reserve for given orderId, userId, serviceId and amount.
// @Summary     Remove reserve
// @Tags        Reserves
// @Accept      json
// @Produce     json
// @Param       inJSON body     models.PayloadReserve true "In JSON with user_id, service_id, order_id and amount"
// @Success     200    {string} status                "OK"
// @Failure     400    {object} models.PayloadErr     "Error"
// @Router      /reserve/ [delete]
func (h *Handler) DeleteReserve(c *fiber.Ctx) error {
	payload := models.PayloadReserve{}
	if err := c.BodyParser(&payload); err != nil {
		return returnBadRequest(err, c)
	}

	err := h.DB.DeleteReserve(payload.UserID, payload.ServiceID, payload.OrderID, utils.MoneyToInt(payload.Amount))
	if err != nil {
		return returnBadRequest(err, c)
	}

	return c.SendStatus(fiber.StatusOK)
}

// Purchase performs purchase for given orderId, userId, serviceId and amount.
// @Description Perform purchase for given orderId, userId, serviceId and amount.
// @Summary     Perform purchase
// @Tags        Purchases
// @Accept      json
// @Produce     json
// @Param       inJSON body     models.PayloadReserve true "In JSON with user_id, service_id, order_id and amount"
// @Success     200    {string} status                "OK"
// @Failure     400    {object} models.PayloadErr     "Error"
// @Router      /purchase/ [post]
func (h *Handler) Purchase(c *fiber.Ctx) error {
	payload := models.PayloadReserve{}
	if err := c.BodyParser(&payload); err != nil {
		return returnBadRequest(err, c)
	}

	err := h.DB.Purchase(payload.UserID, payload.ServiceID, payload.OrderID, utils.MoneyToInt(payload.Amount))
	if err != nil {
		return returnBadRequest(err, c)
	}

	return c.SendStatus(fiber.StatusOK)
}

// AddServices adds multiple services
// @Description Add multiple services
// @Summary     Add multiple services
// @Tags        Services
// @Accept      json
// @Produce     json
// @Param       inJSON body     []models.Service  true "Array of services"
// @Success     200    {string} status            "OK"
// @Failure     400    {object} models.PayloadErr "Error"
// @Router      /services/ [post]
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

// GetService gets service by id
// @Description Get service by given id
// @Summary     Get service
// @Tags        Services
// @Accept      json
// @Produce     json
// @Param       inJSON body     models.PayloadId  true "In JSON with Service ID"
// @Success     200    {object} models.Service    "Service"
// @Failure     400    {object} models.PayloadErr "Error"
// @Router      /services/ [get]
func (h *Handler) GetService(c *fiber.Ctx) error {
	payload := models.PayloadId{}
	if err := c.BodyParser(&payload); err != nil {
		return returnBadRequest(err, c)
	}

	service, err := h.DB.GetService(payload.ID)
	if err != nil {
		return returnBadRequest(err, c)
	}

	return c.JSON(service)
}

// DeleteService deletes service by id
// @Description Delete service by given id
// @Summary     Delete service
// @Tags        Services
// @Accept      json
// @Produce     json
// @Param       inJSON body     models.PayloadId  true "In JSON with Service ID"
// @Success     200    {string} status            "OK"
// @Failure     400    {object} models.PayloadErr "Error"
// @Router      /services/ [delete]
func (h *Handler) DeleteService(c *fiber.Ctx) error {
	payload := models.PayloadId{}
	if err := c.BodyParser(&payload); err != nil {
		return returnBadRequest(err, c)
	}

	err := h.DB.DeleteService(payload.ID)
	if err != nil {
		return returnBadRequest(err, c)
	}

	return c.SendStatus(fiber.StatusOK)
}

// GetReport returns csv report file by given year and month
// @Description Get csv report file by given year and month
// @Summary     Get csv report file
// @Tags        Reports
// @Accept      json
// @Produce     plain
// @Param       year  path     integer           true "Year"
// @Param       month path     integer           true "Month"
// @Success     200   {string} string            "CSV file"
// @Failure     404   {object} models.PayloadErr "CSV file not found"
// @Failure     400   {object} models.PayloadErr "Error"
// @Router      /report/{year}/{month}/report.csv [get]
func (h *Handler) GetReport(c *fiber.Ctx) error {
	payload := models.PayloadDate{}
	if err := c.ParamsParser(&payload); err != nil {
		return returnBadRequest(err, c)
	}
	if payload.Month < 1 || payload.Month > 12 {
		return returnBadRequest(errors.New("handler: get report: wrong month input"), c)
	}

	filePath := utils.GetReportFilePath(payload.Year, payload.Month)
	if _, err := os.Stat(filePath); err == nil {
		return c.SendFile(filePath, false)
	} else if errors.Is(err, os.ErrNotExist) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": fmt.Errorf("handler: get report: report from %d.%d doesn't exist", payload.Month, payload.Year),
		})
	} else {
		return returnBadRequest(err, c)
	}
}

// CreateReport returns link to csv report file by given year and month
// @Description Get link to csv report file by given year and month
// @Summary     Get link to csv report file
// @Tags        Reports
// @Accept      json
// @Produce     json
// @Param       inJSON body     models.PayloadId   true "In JSON with Service ID"
// @Success     200    {object} models.PayloadLink "Link to CSV file"
// @Failure     400    {object} models.PayloadErr  "Error"
// @Router      /report/ [post]
func (h *Handler) CreateReport(c *fiber.Ctx) error {
	payload := models.PayloadDate{}
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
