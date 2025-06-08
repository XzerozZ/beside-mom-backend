package usecases

import (
	"Beside-Mom-BE/configs"
	"Beside-Mom-BE/modules/entities"
	"Beside-Mom-BE/modules/repositories"
	"Beside-Mom-BE/pkg/utils"
	"mime/multipart"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type KidUseCase interface {
	CreateKid(kid *entities.Kid, image *multipart.FileHeader, ctx *fiber.Ctx) (*entities.Kid, error)
	GetKidByID(id string) (map[string]interface{}, error)
	UpdateKidByID(id string, kid *entities.Kid, image *multipart.FileHeader, ctx *fiber.Ctx) (*entities.Kid, error)
}

type KidUseCaseImpl struct {
	repo repositories.KidsRepository
	supa configs.Supabase
}

func NewKidUseCase(repo repositories.KidsRepository, supa configs.Supabase) *KidUseCaseImpl {
	return &KidUseCaseImpl{
		repo: repo,
		supa: supa,
	}
}

func (u *KidUseCaseImpl) CreateKid(kid *entities.Kid, image *multipart.FileHeader, ctx *fiber.Ctx) (*entities.Kid, error) {
	if image != nil {
		fileName := uuid.New().String() + "_title.jpg"
		if err := ctx.SaveFile(image, "./uploads/"+fileName); err != nil {
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

		kid.ImageLink = imageUrl
	}

	createdKid, err := u.repo.CreateKid(kid)
	if err != nil {
		return nil, err
	}

	return createdKid, nil
}

func (u *KidUseCaseImpl) GetKidByID(id string) (map[string]interface{}, error) {
	kid, err := u.repo.GetKidByID(id)
	if err != nil {
		return nil, err
	}

	years, weeks, days, err := utils.CalculateAgeDetailed(kid.BirthDate)
	if err != nil {
		return nil, err
	}

	kidData := map[string]interface{}{
		"id":          kid.ID,
		"firstname":   kid.Firstname,
		"lastname":    kid.Lastname,
		"username":    kid.Username,
		"sex":         kid.Sex,
		"blood":       kid.BloodType,
		"rh":          kid.RHType,
		"imagelink":   kid.ImageLink,
		"birthdate":   kid.BirthDate,
		"birthweight": kid.BirthWeight,
		"birthlength": kid.BirthLength,
		"note":        kid.Note,
		"growth":      kid.Growth,
		"days":        days,
		"weeks":       weeks,
		"years":       years,
	}

	return kidData, nil
}

func (u *KidUseCaseImpl) GetKidByIDForUser(id string) (map[string]interface{}, error) {
	kid, err := u.repo.GetKidByIDForUser(id)
	if err != nil {
		return nil, err
	}

	years, weeks, days, err := utils.CalculateAgeDetailed(kid.BirthDate)
	if err != nil {
		return nil, err
	}

	kidData := map[string]interface{}{
		"id":          kid.ID,
		"firstname":   kid.Firstname,
		"lastname":    kid.Lastname,
		"username":    kid.Username,
		"sex":         kid.Sex,
		"blood":       kid.BloodType,
		"rh":          kid.RHType,
		"imagelink":   kid.ImageLink,
		"birthdate":   kid.BirthDate,
		"birthweight": kid.BirthWeight,
		"birthlength": kid.BirthLength,
		"note":        kid.Note,
		"days":        days,
		"weeks":       weeks,
		"years":       years,
	}

	return kidData, nil
}

func (u *KidUseCaseImpl) UpdateKidByID(id string, kid *entities.Kid, image *multipart.FileHeader, ctx *fiber.Ctx) (*entities.Kid, error) {
	existingKid, err := u.repo.GetKidByID(id)
	if err != nil {
		return nil, err
	}

	if image != nil {
		fileName := uuid.New().String() + "_title.jpg"
		if err := ctx.SaveFile(image, "./uploads/"+fileName); err != nil {
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

		if err := utils.DeleteImage(existingKid.ImageLink, u.supa); err != nil {
			return nil, err
		}

		existingKid.ImageLink = imageUrl
	}

	existingKid.Firstname = kid.Firstname
	existingKid.Lastname = kid.Lastname
	existingKid.Username = kid.Username
	existingKid.BirthDate = kid.BirthDate
	existingKid.BirthLength = kid.BirthLength
	existingKid.BirthWeight = kid.BirthWeight
	existingKid.BloodType = kid.BloodType
	existingKid.RHType = kid.RHType
	existingKid.Note = kid.Note
	existingKid.Sex = kid.Sex
	updatedKid, err := u.repo.UpdateKidByID(existingKid)
	if err != nil {
		return nil, err
	}

	return updatedKid, nil
}
