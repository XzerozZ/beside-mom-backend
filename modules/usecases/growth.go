package usecases

import (
	"Beside-Mom-BE/modules/entities"
	"Beside-Mom-BE/modules/repositories"
	"Beside-Mom-BE/pkg/utils"
	"errors"
	"time"

	"gorm.io/gorm"
)

type GrowthUseCaseImpl struct {
	repo    repositories.GrowthRepository
	kidRepo repositories.KidsRepository
}

type GrowthUseCase interface {
	CreateGrowth(kidID string, growth *entities.Growth, date time.Time) (*entities.Growth, error)
	GetSummary(kidID string) ([]map[string]interface{}, error)
	GetAllGrowth(kidID string) ([]entities.Growth, error)
	UpdateGrowthByID(id string, growth *entities.Growth) (*entities.Growth, error)
}

func NewGrowthUseCase(repo repositories.GrowthRepository, kidRepo repositories.KidsRepository) *GrowthUseCaseImpl {
	return &GrowthUseCaseImpl{
		repo:    repo,
		kidRepo: kidRepo,
	}
}

func (u *GrowthUseCaseImpl) CreateGrowth(kidID string, growth *entities.Growth, date time.Time) (*entities.Growth, error) {
	if growth.Length <= 0 || growth.Weight <= 0 {
		return nil, errors.New("length and weight must be positive numbers")
	}

	kid, err := u.kidRepo.GetKidByID(kidID)
	if err != nil {
		return nil, err
	}

	months, err := utils.CompareAgeKid(kid.BirthDate, date)
	if err != nil {
		return nil, err
	}

	growth.Months = months
	growth.CreatedAt = date
	growth.UpdatedAt = date
	existingGrowth, err := u.repo.GetLatestGrowthByKidID(kidID, months)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		createdGrowth, err := u.repo.CreateGrowth(growth)
		if err != nil {
			return nil, err
		}

		return createdGrowth, nil
	} else if err != nil {
		return nil, err
	}

	existingGrowth.Length = growth.Length
	existingGrowth.Weight = growth.Weight
	existingGrowth.CreatedAt = date
	updatedGrowth, err := u.repo.UpdateGrowth(existingGrowth)
	if err != nil {
		return nil, err
	}

	return updatedGrowth, nil
}

func (u *GrowthUseCaseImpl) GetSummary(kidID string) ([]map[string]interface{}, error) {
	return u.repo.GetSummary(kidID)
}

func (u *GrowthUseCaseImpl) GetAllGrowth(kidID string) ([]entities.Growth, error) {
	return u.repo.GetAllGrowth(kidID)
}

func (u *GrowthUseCaseImpl) UpdateGrowthByID(id string, growth *entities.Growth) (*entities.Growth, error) {
	existingGrowth, err := u.repo.GetGrowthByID(id)
	if err != nil {
		return nil, err
	}

	existingGrowth.Length = growth.Length
	existingGrowth.Weight = growth.Weight
	existingGrowth.UpdatedAt = growth.CreatedAt
	updatedGrowth, err := u.repo.UpdateGrowth(existingGrowth)
	if err != nil {
		return nil, err
	}

	return updatedGrowth, nil
}
