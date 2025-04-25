package usecases

import (
	"Beside-Mom-BE/modules/entities"
	"Beside-Mom-BE/modules/repositories"
)

type EvaluateUseCaseImpl struct {
	repo repositories.EvaluateRepository
}

type EvaluateUseCase interface {
	GetAllEvaluate(kidID string) ([]entities.Evaluate, error)
}

func NewEvaluateUseCase(repo repositories.EvaluateRepository) *EvaluateUseCaseImpl {
	return &EvaluateUseCaseImpl{repo: repo}
}

func (u *EvaluateUseCaseImpl) GetAllEvaluate(kidID string) ([]entities.Evaluate, error) {
	return u.repo.GetAllEvaluate(kidID)
}
