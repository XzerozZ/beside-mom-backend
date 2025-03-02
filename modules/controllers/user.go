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

	user := &entities.User{
		ID:        uuid.New().String(),
		Firstname: form.Value["firstname"][0],
		Lastname:  form.Value["lastname"][0],
		Email:     form.Value["email"][0],
	}

	data, err := c.usecase.CreateUser(user)
	if err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	var kids []*entities.Kid
	fileHeaders := form.File["images"]
	for i := 1; i < len(form.Value["firstname"]); i++ {
		birthWeight, _ := strconv.ParseFloat(form.Value["birthweight"][i-1], 64)
		birthLength, _ := strconv.ParseFloat(form.Value["birthlength"][i-1], 64)
		birthDate, err := time.Parse("2006-01-02", form.Value["birthdate"][i-1])
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
			Firstname:   form.Value["firstname"][i],
			Lastname:    form.Value["lastname"][i],
			Username:    form.Value["username"][i-1],
			Sex:         form.Value["sex"][i-1],
			BirthDate:   birthDate,
			BloodType:   form.Value["bloodtype"][i-1],
			BirthWeight: birthWeight,
			BirthLength: birthLength,
			Note:        form.Value["note"][i-1],
			UserID:      user.ID,
		}

		var image *multipart.FileHeader
		if len(fileHeaders) > i-1 {
			image = fileHeaders[i-1]
		}

		createdKid, err := c.kidusecase.CreateKid(kid, image, ctx)
		if err != nil {
			return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
				"status":      fiber.ErrInternalServerError.Message,
				"status_code": fiber.ErrInternalServerError.Code,
				"message":     "Failed to create kid: " + err.Error(),
				"result":      nil,
			})
		}

		kids = append(kids, createdKid)
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "User created successfully",
		"result": fiber.Map{
			"Mom":  data,
			"Kids": kids,
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
