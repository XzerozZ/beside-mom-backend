package controllers

import (
	"Beside-Mom-BE/modules/usecases"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type HistoryController struct {
	usecase usecases.HistoryUseCase
}

func NewHistoryController(usecase usecases.HistoryUseCase) *HistoryController {
	return &HistoryController{usecase: usecase}
}

func (c *HistoryController) CreateHistoryHandler(ctx *fiber.Ctx) error {
	idParam := ctx.Params("times")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid Evaluate Times",
			"result":      nil,
		})
	}

	kidID := ctx.Params("id")
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Invalid form data",
			"result":      nil,
		})
	}

	rawAnswers := form.Value["answer"]
	if len(rawAnswers) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Missing answer[] in form-data",
			"result":      nil,
		})
	}

	answers := make([]bool, 0, len(rawAnswers))
	for _, val := range rawAnswers {
		parsed, err := strconv.ParseBool(val)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":      "Error",
				"status_code": fiber.StatusBadRequest,
				"message":     fmt.Sprintf("Invalid answer value: %s", val),
				"result":      nil,
			})
		}

		answers = append(answers, parsed)
	}

	err = c.usecase.CreateHistory(id, kidID, answers)
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
		"message":     "History created successfully",
	})
}

func (c *HistoryController) GetHistoryHandler(ctx *fiber.Ctx) error {
	idParam := ctx.Params("times")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid Evaluate Time",
			"result":      nil,
		})
	}

	kidID := ctx.Params("id")
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	data, err := c.usecase.GetHistoryOfEvaluate(id, kidID)
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
		"message":     "History retrieved successfully",
		"result":      data,
	})
}

func (c *HistoryController) GetLatestHistoryHandler(ctx *fiber.Ctx) error {
	idParam := ctx.Params("times")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid Evaluate Time",
			"result":      nil,
		})
	}

	kidID := ctx.Params("id")
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	data, err := c.usecase.GetLatestHistoryOfEvaluate(id, kidID)
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
		"message":     "History retrieved successfully",
		"result":      data,
	})
}
