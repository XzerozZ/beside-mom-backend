package usecases

import (
	"Beside-Mom-BE/modules/entities"
	"Beside-Mom-BE/modules/repositories"
	"errors"

	"github.com/google/uuid"
)

type HistoryUseCase interface {
	CreateHistory(evaluateTimes int, kidID string, answers []bool) error
	GetHistoryOfEvaluate(times int, kidID string) ([]entities.History, error)
}

type HistoryUseCaseImpl struct {
	repo    repositories.HistoryRepository
	evaRepo repositories.EvaluateRepository
}

func NewHistoryUseCase(repo repositories.HistoryRepository, evaRepo repositories.EvaluateRepository) *HistoryUseCaseImpl {
	return &HistoryUseCaseImpl{
		repo:    repo,
		evaRepo: evaRepo,
	}
}

func (u *HistoryUseCaseImpl) CreateHistory(evaluateTimes int, kidID string, answers []bool) error {
	data, err := u.repo.GetLatestHistoryPerQuiz(evaluateTimes, kidID)
	if err != nil {
		return err
	}

	if len(answers) != len(data) {
		return errors.New("please answer all quizzes")
	}

	err = u.repo.DeleteHistoryWithTimes(evaluateTimes, kidID, 0)
	if err != nil {
		return err
	}

	allPass := true
	for i, d := range data {
		d.ID = uuid.New().String()
		d.Times = 1
		d.Status = true
		d.Answer = answers[i]
		if answers[i] {
			d.Solution = "ผ่าน"
		} else {
			d.Solution = "ไม่ผ่าน"
			allPass = false
		}

		err = u.repo.CreateHisotry(d)
		if err != nil {
			return err
		}
	}

	evalSolution := "ผ่านการประเมิน"
	if !allPass {
		evalSolution = "ไม่ผ่านการประเมินบางประการ"
	}

	err = u.evaRepo.UpdateEvaluate(evaluateTimes, kidID, evalSolution)
	if err != nil {
		return err
	}

	return nil
}

func (u *HistoryUseCaseImpl) GetHistoryOfEvaluate(times int, kidID string) ([]entities.History, error) {
	return u.repo.GetLatestHistoryPerQuiz(times, kidID)
}
