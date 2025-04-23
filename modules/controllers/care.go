package controllers

import (
	"Beside-Mom-BE/modules/entities"
	"Beside-Mom-BE/modules/usecases"
	"mime/multipart"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CareController struct {
	usecase usecases.CareUseCase
}

func NewCareController(usecase usecases.CareUseCase) *CareController {
	return &CareController{usecase: usecase}
}

func (c *CareController) CreateCareHandler(ctx *fiber.Ctx) error {
	form, err := ctx.MultipartForm()
	if err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Invalid form data",
			"result":      nil,
		})
	}

	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	title := form.Value["title"]
	if len(title) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Title is required",
			"result":      nil,
		})
	}

	desc := form.Value["desc"]
	if len(desc) == 0 {
		desc = []string{""}
	}

	typeValue := form.Value["type"]
	if len(typeValue) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Type is required",
			"result":      nil,
		})
	}

	care := entities.Care{
		ID:          uuid.New().String(),
		Type:        typeValue[0],
		Title:       title[0],
		Description: desc[0],
		UserID:      userID,
	}

	fileHeaders := form.File["banners"]
	var banner *multipart.FileHeader
	if len(fileHeaders) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Please add banner of the care",
			"result":      nil,
		})
	} else {
		banner = fileHeaders[0]
	}

	var createdCare *entities.Care
	switch typeValue[0] {
	case "image":
		files := form.File["link"]
		if len(files) == 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":      "Error",
				"status_code": fiber.StatusBadRequest,
				"message":     "No image files uploaded",
				"result":      nil,
			})
		}

		createdCare, err = c.usecase.CreateCarewithUploadImages(care, banner, files, ctx)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":      "Error",
				"status_code": fiber.StatusInternalServerError,
				"message":     err.Error(),
				"result":      nil,
			})
		}

	case "video":
		videoFile, _ := ctx.FormFile("link")
		videoLink := form.Value["link"]
		if videoFile != nil && len(videoLink) > 0 && videoLink[0] != "" {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":      "Error",
				"status_code": fiber.StatusBadRequest,
				"message":     "Cannot provide both video file and video link",
				"result":      nil,
			})
		}

		if videoFile != nil {
			file, err := videoFile.Open()
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":      "Error",
					"status_code": fiber.StatusInternalServerError,
					"message":     "Failed to open video file",
					"result":      nil,
				})
			}

			defer file.Close()
			createdCare, err = c.usecase.CreateCarewithUploadVideo(care, videoFile, file, banner, ctx)
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":      "Error",
					"status_code": fiber.StatusInternalServerError,
					"message":     err.Error(),
					"result":      nil,
				})
			}
		} else if len(videoLink) > 0 && videoLink[0] != "" {
			createdCare, err = c.usecase.CreateCarewithVideoLink(care, videoLink[0], banner, ctx)
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":      "Error",
					"status_code": fiber.StatusInternalServerError,
					"message":     err.Error(),
					"result":      nil,
				})
			}
		} else {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":      "Error",
				"status_code": fiber.StatusBadRequest,
				"message":     "No video file or link provided",
				"result":      nil,
			})
		}
	default:
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid care type",
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusCreated,
		"message":     "Care created successfully",
		"result":      createdCare,
	})
}

func (c *CareController) GetCareByID(ctx *fiber.Ctx) error {
	CareID := ctx.Params("id")
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	care, err := c.usecase.GetCareByID(CareID)
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
		"message":     "Care retrieved successfully",
		"result":      care,
	})
}

func (c *CareController) GetAllCareHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	cares, err := c.usecase.GetAllCare()
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
		"message":     "Care retrieved successfully",
		"result":      cares,
	})
}

func (c *CareController) UpdateCareHandler(ctx *fiber.Ctx) error {
	CareID := ctx.Params("id")
	form, err := ctx.MultipartForm()
	if err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Invalid form data",
			"result":      nil,
		})
	}

	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	title := form.Value["title"]
	if len(title) == 0 || strings.TrimSpace(title[0]) == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Title is required and cannot be empty",
			"result":      nil,
		})
	}

	desc := form.Value["desc"]
	if len(desc) == 0 {
		desc = []string{""}
	}

	typeValue := form.Value["type"]
	if len(typeValue) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Type is required",
			"result":      nil,
		})
	}

	care := entities.Care{
		Type:        typeValue[0],
		Title:       title[0],
		Description: desc[0],
	}

	banner, _ := ctx.FormFile("banners")
	var updatedCare *entities.Care
	switch typeValue[0] {
	case "image":
		existingAssets, err := c.usecase.GetCareByID(CareID)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":      "Error",
				"status_code": fiber.StatusInternalServerError,
				"message":     "Failed to retrieve existing images",
				"result":      nil,
			})
		}

		deleteAssets := form.Value["delete_assets"]
		files := form.File["link"]
		remainingImages := len(existingAssets.Assets) - len(deleteAssets) + len(files)
		if remainingImages < 1 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":      "Error",
				"status_code": fiber.StatusBadRequest,
				"message":     "At least one image must be present",
				"result":      nil,
			})
		}

		updatedCare, err = c.usecase.UpdateCarewithUploadImages(CareID, care, files, deleteAssets, banner, ctx)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":      "Error",
				"status_code": fiber.StatusInternalServerError,
				"message":     err.Error(),
				"result":      nil,
			})
		}
	case "video":
		videoFile, _ := ctx.FormFile("link")
		videoLink := form.Value["link"]
		if (len(videoLink) > 0 && videoLink[0] != "") && videoFile != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":      "Error",
				"status_code": fiber.StatusBadRequest,
				"message":     "Please provide either a video file or a video link, not both",
				"result":      nil,
			})
		}

		if videoFile != nil {
			if videoFile.Size > 2*1024*1024*1024 {
				return ctx.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{
					"status":      "Error",
					"status_code": fiber.StatusRequestEntityTooLarge,
					"message":     "Video file too large (max 2GB)",
					"result":      nil,
				})
			}

			file, err := videoFile.Open()
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":      "Error",
					"status_code": fiber.StatusInternalServerError,
					"message":     "Failed to open video file",
					"result":      nil,
				})
			}

			defer file.Close()
			updatedCare, err = c.usecase.UpdateCarewithUploadVideo(CareID, care, videoFile, file, banner, ctx)
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":      "Error",
					"status_code": fiber.StatusInternalServerError,
					"message":     err.Error(),
					"result":      nil,
				})
			}
		} else if len(videoLink) > 0 && videoLink[0] != "" {
			updatedCare, err = c.usecase.UpdateCarewithVideoLink(CareID, care, videoLink[0], banner, ctx)
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":      "Error",
					"status_code": fiber.StatusInternalServerError,
					"message":     err.Error(),
					"result":      nil,
				})
			}
		} else {
			updatedCare, err = c.usecase.UpdateCareByID(CareID, care)
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":      "Error",
					"status_code": fiber.StatusInternalServerError,
					"message":     err.Error(),
					"result":      nil,
				})
			}
		}
	default:
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid care type",
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusCreated,
		"message":     "Care updated successfully",
		"result":      updatedCare,
	})
}

func (c *CareController) DeleteCareCareHandler(ctx *fiber.Ctx) error {
	CareID := ctx.Params("id")

	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	if err := c.usecase.DeleteCareByID(CareID); err != nil {
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
		"message":     "Care deleted successfully",
	})
}
