package controllers

import (
	"Beside-Mom-BE/modules/entities"
	"Beside-Mom-BE/modules/usecases"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AppController struct {
	usecase usecases.AppUseCase
}

func NewAppController(usecase usecases.AppUseCase) *AppController {
	return &AppController{usecase: usecase}
}

func (c *AppController) CreateAppointmentHandler(ctx *fiber.Ctx) error {
	Id := ctx.Params("userID")
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	date := ctx.FormValue("date")
	startTime := ctx.FormValue("start_time")
	building := ctx.FormValue("building")
	requirement := ctx.FormValue("requirement")
	doctor := ctx.FormValue("doctor")
	if date == "" || startTime == "" || building == "" || doctor == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Missing required fields",
			"result":      nil,
		})
	}

	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid date format, expected YYYY-MM-DD",
			"result":      nil,
		})
	}

	parsedStartTime, err := time.Parse("15:04", startTime)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid start time format, expected HH:MM",
			"result":      nil,
		})
	}

	appointment := entities.Appointment{
		ID:          uuid.New().String(),
		UserID:      Id,
		Date:        parsedDate,
		StartTime:   parsedStartTime,
		Building:    building,
		Doctor:      doctor,
		Requirement: requirement,
		Status:      1,
	}

	createdApp, err := c.usecase.CreateAppointment(&appointment)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusInternalServerError,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusCreated,
		"message":     "Appointment created successfully",
		"result":      createdApp,
	})
}

func (c *AppController) GetAppByIDHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	data, err := c.usecase.GetAppByID(id)
	if err != nil {
		return ctx.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"status":      fiber.ErrNotFound.Message,
			"status_code": fiber.ErrNotFound.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "Appointment retrieved successfully",
		"result":      data,
	})
}

func (c *AppController) GetAppHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	role, ok := ctx.Locals("role").(string)
	if !ok || (role != "User" && role != "Admin") {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusForbidden,
			"message":     "Forbidden: Invalid role",
			"result":      nil,
		})
	}

	var data interface{}
	var err error
	if role == "Admin" {
		data, err = c.usecase.GetAllApp()
	} else {
		data, err = c.usecase.GetAppByUserID(userID)
	}

	if err != nil {
		return ctx.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"status":      fiber.ErrNotFound.Message,
			"status_code": fiber.ErrNotFound.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "Appointments retrieved successfully",
		"result":      data,
	})
}

func (c *AppController) GetAllAppUserIDHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	data, err := c.usecase.GetAppByUserID(id)
	if err != nil {
		return ctx.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"status":      fiber.ErrNotFound.Message,
			"status_code": fiber.ErrNotFound.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "Appointment retrieved successfully",
		"result":      data,
	})
}

func (c *AppController) GetAppInProgressUserIDHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	data, err := c.usecase.GetAppInProgressByUserID(id)
	if err != nil {
		return ctx.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"status":      fiber.ErrNotFound.Message,
			"status_code": fiber.ErrNotFound.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "Appointment retrieved successfully",
		"result":      data,
	})
}

func (c *AppController) UpdateAppByIDHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	date := ctx.FormValue("date")
	startTime := ctx.FormValue("start_time")
	building := ctx.FormValue("building")
	requirement := ctx.FormValue("requirement")
	doctor := ctx.FormValue("doctor")
	status := ctx.FormValue("status")
	if date == "" || startTime == "" || building == "" || doctor == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Missing required fields",
			"result":      nil,
		})
	}

	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid date format, expected YYYY-MM-DD",
			"result":      nil,
		})
	}

	parsedStartTime, err := time.Parse("15:04", startTime)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid start time format, expected HH:MM",
			"result":      nil,
		})
	}

	parsedStatus, err := strconv.Atoi(status)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid length value",
			"result":      nil,
		})
	}

	appointment := entities.Appointment{
		Date:        parsedDate,
		StartTime:   parsedStartTime,
		Building:    building,
		Doctor:      doctor,
		Requirement: requirement,
		Status:      parsedStatus,
	}

	updatedApp, err := c.usecase.UpdateAppByID(id, &appointment)
	if err != nil {
		return ctx.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"status":      fiber.ErrNotFound.Message,
			"status_code": fiber.ErrNotFound.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "Appointment update successfully",
		"result":      updatedApp,
	})
}

func (c *AppController) DeleteAppByIDHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	if err := c.usecase.DeleteAppByID(id); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusInternalServerError,
			"message":     "Something went wrong",
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "Appointment deleted successfully",
		"result":      nil,
	})
}
