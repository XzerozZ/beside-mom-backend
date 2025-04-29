package controllers

import (
	"Beside-Mom-BE/modules/entities"
	"Beside-Mom-BE/modules/usecases"
	"mime/multipart"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserController struct {
	usecase    usecases.UserUseCase
	kidusecase usecases.KidUseCase
}

func NewUserController(usecase usecases.UserUseCase, kidusecase usecases.KidUseCase) *UserController {
	return &UserController{
		usecase:    usecase,
		kidusecase: kidusecase,
	}
}

func (c *UserController) CreateUserandKidsHandler(ctx *fiber.Ctx) error {
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
	user := &entities.User{
		ID:        uuid.New().String(),
		Firstname: form.Value["firstname"][0],
		Lastname:  form.Value["lastname"][0],
		Email:     form.Value["email"][0],
	}

	var image *multipart.FileHeader
	if len(fileHeaders) > 0 {
		image = fileHeaders[0]
	}

	data, err := c.usecase.CreateUser(user, image, ctx)
	if err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	if len(form.Value["firstname"]) <= 1 {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "At least one kid is required",
			"result":      nil,
		})
	}

	if len(form.Value["firstname"]) > 2 {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Only one kid is allowed",
			"result":      nil,
		})
	}

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
		Firstname:   form.Value["firstname"][1],
		Lastname:    form.Value["lastname"][1],
		Username:    form.Value["username"][0],
		Sex:         form.Value["sex"][0],
		BirthDate:   birthDate,
		BloodType:   form.Value["bloodtype"][0],
		BirthWeight: birthWeight,
		BirthLength: birthLength,
		Note:        form.Value["note"][0],
		UserID:      user.ID,
	}

	kid, err = c.kidusecase.CreateKid(kid, fileHeaders[1], ctx)
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
		"message":     "User created successfully",
		"result": fiber.Map{
			"Mom":  data,
			"Kids": kid,
		},
	})
}

func (c *UserController) GetMomByIDHandler(ctx *fiber.Ctx) error {
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

	mom, err := c.usecase.GetMomByID(momID)
	if err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "Mom retrieved successfully",
		"result":      mom,
	})
}

func (c *UserController) GetAllMomHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	data, err := c.usecase.GetAllMom()
	if err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "Mom retrieved successfully",
		"result":      data,
	})
}

func (c *UserController) UpdateUserByIDHandler(ctx *fiber.Ctx) error {
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

	images, _ := ctx.FormFile("images")
	updatedUser, err := c.usecase.UpdateUserByID(momID, images, ctx)
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
		"result":      updatedUser,
	})
}

func (c *UserController) DeleteUserHandler(ctx *fiber.Ctx) error {
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

	err := c.usecase.DeleteUser(momID)
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
		"message":     "User deleted successfully",
	})
}
