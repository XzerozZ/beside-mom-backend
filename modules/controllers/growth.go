package controllers

import (
	"Beside-Mom-BE/modules/entities"
	"Beside-Mom-BE/modules/usecases"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type GrowthController struct {
	usecase usecases.GrowthUseCase
}

func NewGrowthController(usecase usecases.GrowthUseCase) *GrowthController {
	return &GrowthController{usecase: usecase}
}

func (c *GrowthController) CreateGrowthHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	kidID := ctx.Params("id")
	lengthStr := ctx.FormValue("length")
	weightStr := ctx.FormValue("weight")
	if lengthStr == "" || weightStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Length and weight are required",
			"result":      nil,
		})
	}

	length, err := strconv.ParseFloat(lengthStr, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid length value",
			"result":      nil,
		})
	}

	weight, err := strconv.ParseFloat(weightStr, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid weight value",
			"result":      nil,
		})
	}

	growth := &entities.Growth{
		ID:     uuid.New().String(),
		Length: length,
		Weight: weight,
		KidID:  kidID,
	}

	growth, err = c.usecase.CreateGrowth(kidID, growth)
	if err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "Growth data created successfully",
		"result":      growth,
	})
}

func (c *GrowthController) GetSummary(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	kidID := ctx.Params("id")
	growth, err := c.usecase.GetSummary(kidID)
	if err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "Growth data retreive successfully",
		"result":      growth,
	})
}

func (c *GrowthController) GetAllGrowth(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	kidID := ctx.Params("id")
	growth, err := c.usecase.GetAllGrowth(kidID)
	if err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "Growth data retrieved successfully",
		"result":      growth,
	})
}

func (c *GrowthController) UpdateGrowthByID(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	id := ctx.Params("id")
	lengthStr := ctx.FormValue("length")
	weightStr := ctx.FormValue("weight")
	if lengthStr == "" || weightStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Length and weight are required",
			"result":      nil,
		})
	}

	date, err := time.Parse("2006-01-02", ctx.FormValue("date"))
	if err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Invalid birthdate format. Use YYYY-MM-DD",
			"result":      nil,
		})
	}

	length, err := strconv.ParseFloat(lengthStr, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid length value",
			"result":      nil,
		})
	}

	weight, err := strconv.ParseFloat(weightStr, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid weight value",
			"result":      nil,
		})
	}

	growth := &entities.Growth{
		Length:    length,
		Weight:    weight,
		UpdatedAt: date,
	}

	growth, err = c.usecase.UpdateGrowthByID(id, growth)
	if err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "Growth data created successfully",
		"result":      growth,
	})
}
