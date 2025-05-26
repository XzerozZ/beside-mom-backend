package controllers

import (
	"Beside-Mom-BE/modules/entities"
	"Beside-Mom-BE/modules/usecases"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	usecase usecases.AuthUseCase
}

func NewAuthController(usecase usecases.AuthUseCase) *AuthController {
	return &AuthController{usecase: usecase}
}

func (c *AuthController) RegisterHandler(ctx *fiber.Ctx) error {
	var req struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	if req.Email == "" || req.Password == "" || req.Firstname == "" || req.Lastname == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Email, Password, Firstname and Lastname is required",
			"result":      nil,
		})
	}

	user := &entities.User{
		Email:     req.Email,
		Password:  req.Password,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
	}

	data, err := c.usecase.RegisterAdmin(user)
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
		"message":     "user created successfully",
		"result":      data,
	})
}

func (c *AuthController) LoginHandler(ctx *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	if req.Email == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Email is missing",
			"result":      nil,
		})
	}

	if req.Password == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Password is missing",
			"result":      nil,
		})
	}

	token, user, err := c.usecase.Login(req.Email, req.Password)
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
		"message":     "Login successful",
		"result": fiber.Map{
			"token": token,
			"u_id":  user.ID,
			"name":  user.Firstname + " " + user.Lastname,
			"role":  user.Role.RoleName,
		},
	})
}

func (c *AuthController) ForgotPasswordHandler(ctx *fiber.Ctx) error {
	type ForgotPasswordRequest struct {
		Email string `json:"email" validate:"required,email"`
	}

	var req ForgotPasswordRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	if req.Email == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Email is missing",
			"result":      nil,
		})
	}

	err := c.usecase.ForgotPassword(req.Email)
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
		"message":     "Sent OTP successfully",
	})
}

func (c *AuthController) VerifyOTPHandler(ctx *fiber.Ctx) error {
	type OTPRequest struct {
		Email string `json:"email" validate:"required,email"`
		OTP   string `json:"otp"`
	}

	var req OTPRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	if req.Email == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Email is missing",
			"result":      nil,
		})
	}

	if req.OTP == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "OTP is missing",
			"result":      nil,
		})
	}

	err := c.usecase.VerifyOTP(req.Email, req.OTP)
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
		"message":     "OTP is correct",
	})
}

func (c *AuthController) ChangedPasswordHandler(ctx *fiber.Ctx) error {
	type OTPRequest struct {
		Email       string `json:"email" validate:"required,email"`
		NewPassword string `json:"newPassword"`
	}

	var req OTPRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	if req.Email == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Email is missing",
			"result":      nil,
		})
	}

	if req.NewPassword == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "NewPassword is missing",
			"result":      nil,
		})
	}

	err := c.usecase.ChangedPassword(req.Email, req.NewPassword)
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
		"message":     "changed password successfully",
	})
}
