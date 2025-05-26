package controllers

import (
	"Beside-Mom-BE/modules/entities"
	"Beside-Mom-BE/modules/usecases"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type KidController struct {
	usecase usecases.KidUseCase
}

func NewKidController(usecase usecases.KidUseCase) *KidController {
	return &KidController{usecase: usecase}
}

func (c *KidController) CreateKidHandler(ctx *fiber.Ctx) error {
	momID := ctx.Params("id")
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
			"message":     "invalid form data",
			"result":      nil,
		})
	}

	fileHeaders := form.File["images"]
	birthWeight, _ := strconv.ParseFloat(form.Value["birthweight"][0], 64)
	birthLength, _ := strconv.ParseFloat(form.Value["birthlength"][0], 64)
	birthDate, err := time.Parse("2006-01-02", form.Value["birthdate"][0])
	if err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Invalid birthdate format. Use YYYY-MM-DD",
			"result":      nil,
		})
	}

	kid := &entities.Kid{
		ID:          uuid.New().String(),
		Firstname:   form.Value["firstname"][0],
		Lastname:    form.Value["lastname"][0],
		Username:    form.Value["username"][0],
		Sex:         form.Value["sex"][0],
		BirthDate:   birthDate,
		BloodType:   form.Value["bloodtype"][0],
		BirthWeight: birthWeight,
		BirthLength: birthLength,
		Note:        form.Value["note"][0],
		UserID:      momID,
	}

	kid, err = c.usecase.CreateKid(kid, fileHeaders[0], ctx)
	if err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     "Failed to create kid: " + err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "Kid added successfully",
		"result":      kid,
	})
}

func (c *KidController) GetKidByIDHandler(ctx *fiber.Ctx) error {
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

	role, ok := ctx.Locals("role").(string)
	if !ok || role == "" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusForbidden,
			"message":     "Forbidden: Missing role",
			"result":      nil,
		})
	}

	var data interface{}
	var err error

	data, err = c.usecase.GetKidByID(kidID)
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
		"message":     "Kid retrieved successfully",
		"result":      data,
	})
}

func (c *KidController) UpdateKidByIDHandler(ctx *fiber.Ctx) error {
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

	images, _ := ctx.FormFile("images")
	birthWeight, _ := strconv.ParseFloat(ctx.FormValue("birthweight"), 64)
	birthLength, _ := strconv.ParseFloat(ctx.FormValue("birthlength"), 64)
	birthDate, err := time.Parse("2006-01-02", ctx.FormValue("birthdate"))
	if err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Invalid birthdate format. Use YYYY-MM-DD",
			"result":      nil,
		})
	}

	kid := &entities.Kid{
		ID:          uuid.New().String(),
		Firstname:   ctx.FormValue("firstname"),
		Lastname:    ctx.FormValue("lastname"),
		Username:    ctx.FormValue("username"),
		Sex:         ctx.FormValue("sex"),
		BirthDate:   birthDate,
		BloodType:   ctx.FormValue("bloodtypet"),
		BirthWeight: birthWeight,
		BirthLength: birthLength,
		Note:        ctx.FormValue("note"),
	}
	updatedKid, err := c.usecase.UpdateKidByID(kidID, kid, images, ctx)
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
		"message":     "User retrieved successfully",
		"result":      updatedKid,
	})
}
