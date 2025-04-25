package usecases

import (
	"Beside-Mom-BE/modules/entities"
	"Beside-Mom-BE/modules/repositories"
	"errors"

	"github.com/google/uuid"
)

type HistoryUseCase interface {
	CreateHistory(evaluateTimes int, kidID string, answers []bool) error
	GetHistoryOfEvaluate(times int, kidID string) (map[int]entities.GroupedHistory, error)
	GetLatestHistoryOfEvaluate(times int, kidID string) ([]entities.History, error)
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
		if answers[i] {
			d.Solution = "ผ่าน"
		} else {
			d.Solution = "ไม่ผ่าน"
			allPass = false
		}

		history := entities.History{
			ID:             uuid.New().String(),
			QuizID:         d.QuizID,
			Answer:         answers[i],
			Status:         true,
			Solution:       d.Solution,
			EvaluatedTimes: d.EvaluatedTimes,
			Times:          d.Times + 1,
			KidID:          d.KidID,
		}

		err = u.repo.CreateHisotry(history)
		if err != nil {
			return err
		}
	}

	evalSolution := "ผ่านการประเมิน"
	Times := data[0].Times + 1
	if !allPass {
		evalSolution = "ไม่ผ่านการประเมินบางประการ"
	}

	err = u.evaRepo.UpdateEvaluate(evaluateTimes, kidID, evalSolution, Times)
	if err != nil {
		return err
	}

	return nil
}

func (u *HistoryUseCaseImpl) GetHistoryOfEvaluate(times int, kidID string) (map[int]entities.GroupedHistory, error) {
	return u.repo.GetHistoryPerQuizByTimez(times, kidID)
}

func (u *HistoryUseCaseImpl) GetLatestHistoryOfEvaluate(times int, kidID string) ([]entities.History, error) {
	return u.repo.GetLatestHistoryPerQuiz(times, kidID)
}
