package usecases

import (
	"Beside-Mom-BE/modules/entities"
	"Beside-Mom-BE/modules/repositories"
	"errors"

	"github.com/google/uuid"
)

type HistoryUseCase interface {
	CreateHistoryInPeriodandHistory(evaluateTimes int, cate int, kidID string, answers []bool) error
	GetHistoryOfEvaluate(times int, kidID string) (map[string]map[int]entities.GroupedHistory, error)
	GetLatestHistoryOfEvaluate(times int, kidID string, cate int) ([]entities.History, error)
	GetHistoryResult(evaluatedTimes int, kidID string) (map[int]entities.GroupedHistory, error)
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

func (u *HistoryUseCaseImpl) CreateHistoryInPeriodandHistory(evaluateTimes int, cate int, kidID string, answers []bool) error {
	data, err := u.repo.GetLatestHistoryPerQuiz(evaluateTimes, cate, kidID)
	if err != nil {
		return err
	}

	if len(answers) != len(data) {
		return errors.New("please answer all quizzes")
	}

	err = u.repo.DeleteHistoryWithTimes(evaluateTimes, kidID, 0, cate)
	if err != nil {
		return err
	}

	for i, d := range data {
		history := entities.History{
			ID:             uuid.New().String(),
			QuizID:         d.QuizID,
			Answer:         answers[i],
			Status:         true,
			EvaluatedTimes: d.EvaluatedTimes,
			Times:          d.Times + 1,
			KidID:          d.KidID,
		}

		err = u.repo.CreateHistory(history)
		if err != nil {
			return err
		}
	}

	latestHistories, err := u.repo.GetLatestHistoryPerEvaluate(evaluateTimes, kidID)
	if err != nil {
		return err
	}

	incomplete := false
	failed := false
	for _, h := range latestHistories {
		if !h.Status {
			incomplete = true
			break
		}

		if !h.Answer {
			failed = true
		}
	}

	var result string
	var status bool
	if incomplete {
		result = "กำลังประเมิน"
		status = false
	} else if failed {
		result = "ไม่ผ่านการประเมินบางประการ"
		status = true
	} else {
		result = "ผ่านการประเมิน"
		status = true
	}

	err = u.evaRepo.UpdateEvaluate(evaluateTimes, kidID, result, status)
	if err != nil {
		return err
	}

	return nil
}

func (u *HistoryUseCaseImpl) GetHistoryOfEvaluate(times int, kidID string) (map[string]map[int]entities.GroupedHistory, error) {
	return u.repo.GetHistoryPerQuizGroupedByCategoryAndTimes(times, kidID)
}

func (u *HistoryUseCaseImpl) GetLatestHistoryOfEvaluate(times int, kidID string, cate int) ([]entities.History, error) {
	return u.repo.GetLatestHistoryPerQuiz(times, cate, kidID)
}

func (u *HistoryUseCaseImpl) GetHistoryResult(evaluatedTimes int, kidID string) (map[int]entities.GroupedHistory, error) {
	return u.repo.GetHistoryResult(evaluatedTimes, kidID)
}
