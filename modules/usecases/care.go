package usecases

import (
	"Beside-Mom-BE/configs"
	"Beside-Mom-BE/modules/entities"
	"Beside-Mom-BE/modules/repositories"
	"Beside-Mom-BE/pkg/utils"
	"io"
	"mime/multipart"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CareUseCase interface {
	CreateCarewithUploadVideo(care entities.Care, videoFile *multipart.FileHeader, file io.Reader, banner *multipart.FileHeader, ctx *fiber.Ctx) (*entities.Care, error)
	CreateCarewithVideoLink(care entities.Care, link string, banner *multipart.FileHeader, ctx *fiber.Ctx) (*entities.Care, error)
	CreateCarewithUploadImages(care entities.Care, banner *multipart.FileHeader, files []*multipart.FileHeader, ctx *fiber.Ctx) (*entities.Care, error)
	GetCareByID(id string) (*entities.Care, error)
	GetAllCare() ([]entities.Care, error)
	UpdateCareByID(id string, care entities.Care) (*entities.Care, error)
	UpdateCarewithUploadVideo(id string, care entities.Care, videoFile *multipart.FileHeader, file io.Reader, banner *multipart.FileHeader, ctx *fiber.Ctx) (*entities.Care, error)
	UpdateCarewithVideoLink(id string, care entities.Care, link string, banner *multipart.FileHeader, ctx *fiber.Ctx) (*entities.Care, error)
	UpdateCarewithUploadImages(id string, care entities.Care, files []*multipart.FileHeader, assetsToDelete []string, banner *multipart.FileHeader, ctx *fiber.Ctx) (*entities.Care, error)
	DeleteCareByID(id string) error
}

type CareUseCaseImpl struct {
	repo repositories.CareRepository
	supa configs.Supabase
}

func NewCareUseCase(repo repositories.CareRepository, supa configs.Supabase) *CareUseCaseImpl {
	return &CareUseCaseImpl{
		repo: repo,
		supa: supa,
	}
}

func (u *CareUseCaseImpl) CreateCarewithUploadVideo(care entities.Care, videoFile *multipart.FileHeader, file io.Reader, banner *multipart.FileHeader, ctx *fiber.Ctx) (*entities.Care, error) {
	assets := []entities.Asset{}
	if videoFile != nil {
		fileName := uuid.New().String() + "_video.mp4"
		videoUrl, err := utils.UploadVideo(fileName, file, u.supa)
		if err != nil {
			return nil, err
		}

		assets = append(assets, entities.Asset{
			ID:   uuid.New().String(),
			Link: videoUrl,
		})
	}

	if banner != nil {
		fileName := uuid.New().String() + "_title.jpg"
		if err := ctx.SaveFile(banner, "./uploads/"+fileName); err != nil {
			return nil, err
		}

		imageUrl, err := utils.UploadImage(fileName, "", u.supa)
		if err != nil {
			os.Remove("./uploads/" + fileName)
			return nil, err
		}

		if err := os.Remove("./uploads/" + fileName); err != nil {
			return nil, err
		}

		care.Banner = imageUrl
	}

	createdCare, err := u.repo.CreateCare(&care, assets)
	if err != nil {
		return nil, err
	}

	return createdCare, nil
}

func (u *CareUseCaseImpl) CreateCarewithVideoLink(care entities.Care, link string, banner *multipart.FileHeader, ctx *fiber.Ctx) (*entities.Care, error) {
	assets := []entities.Asset{
		{
			ID:   uuid.New().String(),
			Link: link,
		},
	}

	if banner != nil {
		fileName := uuid.New().String() + "_title.jpg"
		if err := ctx.SaveFile(banner, "./uploads/"+fileName); err != nil {
			return nil, err
		}

		imageUrl, err := utils.UploadImage(fileName, "", u.supa)
		if err != nil {
			os.Remove("./uploads/" + fileName)
			return nil, err
		}

		if err := os.Remove("./uploads/" + fileName); err != nil {
			return nil, err
		}

		care.Banner = imageUrl
	}

	createdCare, err := u.repo.CreateCare(&care, assets)
	if err != nil {
		return nil, err
	}

	return createdCare, nil
}

func (u *CareUseCaseImpl) CreateCarewithUploadImages(care entities.Care, banner *multipart.FileHeader, files []*multipart.FileHeader, ctx *fiber.Ctx) (*entities.Care, error) {
	assets := []entities.Asset{}
	for _, file := range files {
		fileName := uuid.New().String() + ".jpg"
		if err := ctx.SaveFile(file, "./uploads/"+fileName); err != nil {
			return nil, err
		}

		imageUrl, err := utils.UploadImage(fileName, "", u.supa)
		if err != nil {
			return nil, err
		}

		err = os.Remove("./uploads/" + fileName)
		if err != nil {
			return nil, err
		}

		assets = append(assets, entities.Asset{
			ID:   uuid.New().String(),
			Link: imageUrl,
		})
	}

	if banner != nil {
		fileName := uuid.New().String() + "_title.jpg"
		if err := ctx.SaveFile(banner, "./uploads/"+fileName); err != nil {
			return nil, err
		}

		imageUrl, err := utils.UploadImage(fileName, "", u.supa)
		if err != nil {
			os.Remove("./uploads/" + fileName)
			return nil, err
		}

		if err := os.Remove("./uploads/" + fileName); err != nil {
			return nil, err
		}

		care.Banner = imageUrl
	}

	createdCare, err := u.repo.CreateCare(&care, assets)
	if err != nil {
		return nil, err
	}

	return createdCare, nil
}

func (u *CareUseCaseImpl) GetCareByID(id string) (*entities.Care, error) {
	return u.repo.GetCareByID(id)
}

func (u *CareUseCaseImpl) GetAllCare() ([]entities.Care, error) {
	return u.repo.GetAllCare()
}

func (u *CareUseCaseImpl) UpdateCareByID(id string, care entities.Care) (*entities.Care, error) {
	existingCare, err := u.repo.GetCareByID(id)
	if err != nil {
		return nil, err
	}

	existingCare.Title = care.Title
	existingCare.Description = care.Description
	updatedCare, err := u.repo.UpdateCare(existingCare)
	if err != nil {
		return nil, err
	}

	return updatedCare, nil
}

func (u *CareUseCaseImpl) UpdateCarewithUploadVideo(id string, care entities.Care, videoFile *multipart.FileHeader, file io.Reader, banner *multipart.FileHeader, ctx *fiber.Ctx) (*entities.Care, error) {
	existingCare, err := u.repo.GetCareByID(id)
	if err != nil {
		return nil, err
	}

	existingCare.Title = care.Title
	existingCare.Description = care.Description
	var assets []entities.Asset
	if len(existingCare.Assets) > 0 {
		if err := u.repo.RemoveAssets(id, &existingCare.Assets[0].ID); err != nil {
			return nil, err
		}
	}

	if banner != nil {
		fileName := uuid.New().String() + "_title.jpg"
		if err := ctx.SaveFile(banner, "./uploads/"+fileName); err != nil {
			return nil, err
		}

		imageUrl, err := utils.UploadImage(fileName, "", u.supa)
		if err != nil {
			os.Remove("./uploads/" + fileName)
			return nil, err
		}

		if err := os.Remove("./uploads/" + fileName); err != nil {
			return nil, err
		}

		if err := utils.DeleteImage(existingCare.Banner, u.supa); err != nil {
			return nil, err
		}

		existingCare.Banner = imageUrl
	}

	fileName := uuid.New().String() + "_video.mp4"
	videoUrl, err := utils.UploadVideo(fileName, file, u.supa)
	if err != nil {
		return nil, err
	}

	assets = append(assets, entities.Asset{
		ID:   uuid.New().String(),
		Link: videoUrl,
	})

	existingCare.Assets = assets
	_, err = u.repo.AddAssets(id, assets)
	if err != nil {
		return nil, err
	}

	updatedCare, err := u.repo.UpdateCare(existingCare)
	if err != nil {
		return nil, err
	}

	return updatedCare, nil
}

func (u *CareUseCaseImpl) UpdateCarewithVideoLink(id string, care entities.Care, link string, banner *multipart.FileHeader, ctx *fiber.Ctx) (*entities.Care, error) {
	existingCare, err := u.repo.GetCareByID(id)
	if err != nil {
		return nil, err
	}

	if len(existingCare.Assets) > 0 {
		if err := u.repo.RemoveAssets(id, &existingCare.Assets[0].ID); err != nil {
			return nil, err
		}
	}

	existingCare.Title = care.Title
	existingCare.Description = care.Description
	assets := []entities.Asset{
		{
			ID:   uuid.New().String(),
			Link: link,
		},
	}

	if banner != nil {
		fileName := uuid.New().String() + "_title.jpg"
		if err := ctx.SaveFile(banner, "./uploads/"+fileName); err != nil {
			return nil, err
		}

		imageUrl, err := utils.UploadImage(fileName, "", u.supa)
		if err != nil {
			os.Remove("./uploads/" + fileName)
			return nil, err
		}

		if err := os.Remove("./uploads/" + fileName); err != nil {
			return nil, err
		}

		if err := utils.DeleteImage(existingCare.Banner, u.supa); err != nil {
			return nil, err
		}

		existingCare.Banner = imageUrl
	}

	existingCare.Assets = assets
	_, err = u.repo.AddAssets(id, assets)
	if err != nil {
		return nil, err
	}

	updatedCare, err := u.repo.UpdateCare(existingCare)
	if err != nil {
		return nil, err
	}

	return updatedCare, nil
}

func (u *CareUseCaseImpl) UpdateCarewithUploadImages(id string, care entities.Care, files []*multipart.FileHeader, assetsToDelete []string, banner *multipart.FileHeader, ctx *fiber.Ctx) (*entities.Care, error) {
	existingCare, err := u.repo.GetCareByID(id)
	if err != nil {
		return nil, err
	}

	existingCare.Title = care.Title
	existingCare.Description = care.Description
	if len(assetsToDelete) > 0 {
		for _, assetID := range assetsToDelete {
			if err := u.repo.RemoveAssets(id, &assetID); err != nil {
				return nil, err
			}
		}
	}

	var assets []entities.Asset
	if len(files) > 0 {
		for _, file := range files {
			fileName := uuid.New().String() + ".jpg"
			if err := ctx.SaveFile(file, "./uploads/"+fileName); err != nil {
				return nil, err
			}

			imageUrl, err := utils.UploadImage(fileName, "", u.supa)
			if err != nil {
				return nil, err
			}

			err = os.Remove("./uploads/" + fileName)
			if err != nil {
				return nil, err
			}

			assets = append(assets, entities.Asset{
				ID:   uuid.New().String(),
				Link: imageUrl,
			})
		}
	}

	existingCare.Assets = assets
	if len(assets) > 0 {
		existingCare.Assets = append(existingCare.Assets, assets...)
		_, err = u.repo.AddAssets(id, assets)
		if err != nil {
			return nil, err
		}
	}

	if banner != nil {
		fileName := uuid.New().String() + "_title.jpg"
		if err := ctx.SaveFile(banner, "./uploads/"+fileName); err != nil {
			return nil, err
		}

		imageUrl, err := utils.UploadImage(fileName, "", u.supa)
		if err != nil {
			os.Remove("./uploads/" + fileName)
			return nil, err
		}

		if err := os.Remove("./uploads/" + fileName); err != nil {
			return nil, err
		}

		if err := utils.DeleteImage(existingCare.Banner, u.supa); err != nil {
			return nil, err
		}

		existingCare.Banner = imageUrl
	}

	updatedCare, err := u.repo.UpdateCare(existingCare)
	if err != nil {
		return nil, err
	}

	return updatedCare, nil
}

func (u *CareUseCaseImpl) DeleteCareByID(id string) error {
	existingCare, err := u.repo.GetCareByID(id)
	if err != nil {
		return err
	}

	if err := utils.DeleteImage(existingCare.Banner, u.supa); err != nil {
		return err
	}

	return u.repo.DeleteCare(id)
}
