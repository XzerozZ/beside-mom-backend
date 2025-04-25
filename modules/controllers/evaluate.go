package controllers

import (
	"Beside-Mom-BE/modules/usecases"

	"github.com/gofiber/fiber/v2"
)

type EvaluateController struct {
	usecase usecases.EvaluateUseCase
}

func NewEvaluateController(usecase usecases.EvaluateUseCase) *EvaluateController {
	return &EvaluateController{usecase: usecase}
}

func (c *EvaluateController) GetAllEvaluateHandler(ctx *fiber.Ctx) error {
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

	data, err := c.usecase.GetAllEvaluate(kidID)
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
		"message":     "Evaluate retrieved successfully",
		"result":      data,
	})
}
