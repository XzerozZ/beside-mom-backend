package controllers

import (
	"Beside-Mom-BE/modules/entities"
	"Beside-Mom-BE/modules/usecases"
	"fmt"
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
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid form data: " + err.Error(),
			"result":      nil,
		})
	}

	if len(form.Value["firstname"]) < 2 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Both user and kid firstname are required",
			"result":      nil,
		})
	}

	if len(form.Value["lastname"]) < 2 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Both user and kid lastname are required",
			"result":      nil,
		})
	}

	requiredFields := []string{"pid", "email", "username", "sex", "birthdate", "bloodtype", "beforebirth", "birthweight", "birthlength"}
	for _, field := range requiredFields {
		if len(form.Value[field]) == 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":      "Error",
				"status_code": fiber.StatusBadRequest,
				"message":     fmt.Sprintf("Missing required field: %s", field),
				"result":      nil,
			})
		}
	}

	user := &entities.User{
		ID:        uuid.New().String(),
		PID:       form.Value["pid"][0],
		Firstname: form.Value["firstname"][0],
		Lastname:  form.Value["lastname"][0],
		Email:     form.Value["email"][0],
	}

	fileHeaders := form.File["images"]
	var userImage, kidImage *multipart.FileHeader

	if len(fileHeaders) > 0 {
		userImage = fileHeaders[0]
	}

	if len(fileHeaders) > 1 {
		kidImage = fileHeaders[1]
	}

	userData, err := c.usecase.CreateUser(user, userImage, ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusInternalServerError,
			"message":     "Failed to create user: " + err.Error(),
			"result":      nil,
		})
	}

	beforebirth, err := strconv.ParseInt(form.Value["beforebirth"][0], 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid before birth format",
			"result":      nil,
		})
	}

	birthWeight, err := strconv.ParseFloat(form.Value["birthweight"][0], 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid birth weight format",
			"result":      nil,
		})
	}

	birthLength, err := strconv.ParseFloat(form.Value["birthlength"][0], 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid birth length format",
			"result":      nil,
		})
	}

	birthDate, err := time.Parse("2006-01-02", form.Value["birthdate"][0])
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusBadRequest,
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
		BeforeBirth: int(beforebirth),
		BirthWeight: birthWeight,
		BirthLength: birthLength,
		UserID:      user.ID,
	}

	if len(form.Value["rh"]) > 0 {
		kid.RHType = form.Value["rh"][0]
	}

	if len(form.Value["note"]) > 0 {
		kid.Note = form.Value["note"][0]
	}

	kidData, err := c.kidusecase.CreateKid(kid, kidImage, ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusInternalServerError,
			"message":     "Failed to create kid: " + err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusCreated,
		"message":     "User and kid created successfully",
		"result": fiber.Map{
			"user": userData,
			"kid":  kidData,
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

func (c *UserController) UpdateUserByIDForUserHandler(ctx *fiber.Ctx) error {
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
	updatedUser, err := c.usecase.UpdateUserByIDForUser(userID, images, ctx)
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

func (c *UserController) UpdateUserByIDForAdminHandler(ctx *fiber.Ctx) error {
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
	requiredFields := []string{"email", "pid"}
	for _, field := range requiredFields {
		if len(form.Value[field]) == 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":      "Error",
				"status_code": fiber.StatusBadRequest,
				"message":     fmt.Sprintf("Missing required field: %s", field),
				"result":      nil,
			})
		}
	}

	user := &entities.User{
		ID:        uuid.New().String(),
		PID:       form.Value["pid"][0],
		Firstname: form.Value["firstname"][0],
		Lastname:  form.Value["lastname"][0],
		Email:     form.Value["email"][0],
	}

	var image *multipart.FileHeader
	if len(fileHeaders) > 0 {
		image = fileHeaders[0]
	}

	data, err := c.usecase.UpdateUserByIDForAdmin(momID, user, image, ctx)
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
		"message":     "User updated successfully",
		"result":      data,
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

func (c *UserController) ChatBotHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	message := ctx.FormValue("message")
	result, err := c.usecase.Chat(message)
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
		"message":     result,
	})
}
