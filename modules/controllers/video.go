package controllers

import (
	"Beside-Mom-BE/modules/entities"
	"Beside-Mom-BE/modules/usecases"
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type VideoController struct {
	usecase usecases.VideoUseCase
}

func NewVideoController(usecase usecases.VideoUseCase) *VideoController {
	return &VideoController{usecase: usecase}
}

func (c *VideoController) CreateVideoHandler(ctx *fiber.Ctx) error {
	form, err := ctx.MultipartForm()
	if err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "invalid form data",
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

	videoLink := form.Value["video_link"]
	videoFile, _ := ctx.FormFile("video_link")
	fileHeaders := form.File["banners"]
	var banner *multipart.FileHeader
	if len(fileHeaders) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Please add banner of the video",
			"result":      nil,
		})
	} else {
		banner = fileHeaders[0]
	}

	if (len(videoLink) > 0 && videoLink[0] != "") && videoFile != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Please provide either a video file or a video link, not both",
			"result":      nil,
		})
	}

	if (len(videoLink) == 0 || videoLink[0] == "") && videoFile == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Please provide either a video file or a video link",
			"result":      nil,
		})
	}

	video := entities.Video{
		ID:          uuid.New().String(),
		Title:       title[0],
		Description: desc[0],
		UserID:      userID,
	}

	if len(videoLink) > 0 && videoLink[0] != "" {
		video.Link = videoLink[0]
		data, err := c.usecase.CreateVideowithLink(&video, banner, ctx)
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
			"message":     "Video created successfully with link",
			"result":      data,
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
		data, err := c.usecase.CreateVideo(&video, videoFile, file, banner, ctx)
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
			"message":     "Video created successfully with file",
			"result":      data,
		})
	}

	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"status":      "Error",
		"status_code": fiber.StatusInternalServerError,
		"message":     "Unexpected error in video processing",
		"result":      nil,
	})
}

func (c *VideoController) GetAllVideoHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	data, err := c.usecase.GetAllVideo()
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
		"message":     "Video retrieved successfully",
		"result":      data,
	})
}

func (c *VideoController) GetVideoByIDHandler(ctx *fiber.Ctx) error {
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

	role, ok := ctx.Locals("role").(string)
	if !ok || (role != "User" && role != "Admin") {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusForbidden,
			"message":     "Forbidden: Invalid role",
			"result":      nil,
		})
	}

	if role == "User" {
		if err := c.usecase.IncreaseView(id); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":      "Error",
				"status_code": fiber.StatusInternalServerError,
				"message":     "Failed to increase view count",
				"result":      nil,
			})
		}
	}

	video, err := c.usecase.GetVideoByID(id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusNotFound,
			"message":     "Video not found",
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "Video retrieved successfully",
		"result":      video,
	})
}

func (c *VideoController) UpdateVideoHandler(ctx *fiber.Ctx) error {
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

	form, err := ctx.MultipartForm()
	if err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Invalid form data",
			"result":      nil,
		})
	}

	videoUpdate := &entities.Video{
		Title:       form.Value["title"][0],
		Description: form.Value["desc"][0],
	}

	videoLink := form.Value["video_link"]
	videoFile, _ := ctx.FormFile("video_link")
	banner, _ := ctx.FormFile("banners")
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
		data, err := c.usecase.UpdateVideo(id, videoUpdate, videoFile, file, banner, ctx)
		if err != nil {
			if err.Error() == "unauthorized: user does not own this video" {
				return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"status":      "Error",
					"status_code": fiber.StatusForbidden,
					"message":     err.Error(),
					"result":      nil,
				})
			}
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
			"message":     "Video updated successfully with new file",
			"result":      data,
		})
	} else {
		if len(videoLink) > 0 {
			videoUpdate.Link = videoLink[0]
			data, err := c.usecase.UpdateVideowithLink(id, videoUpdate, banner, ctx)
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
				"message":     "Video updated successfully with new link",
				"result":      data,
			})
		} else {
			data, err := c.usecase.UpdateVideowithLink(id, videoUpdate, banner, ctx)
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
				"message":     "Video updated successfully with new link",
				"result":      data,
			})
		}
	}
}

func (c *VideoController) DeleteVideoByIDHandler(ctx *fiber.Ctx) error {
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

	if err := c.usecase.DeleteVideoByID(id); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusInternalServerError,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "Video deleted successfully",
		"result":      nil,
	})
}
