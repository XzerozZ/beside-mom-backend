package repositories

import (
	"Beside-Mom-BE/modules/entities"
	"time"

	"gorm.io/gorm"
)

type GormHistoryRepository struct {
	db *gorm.DB
}

func NewGormHistoryRepository(db *gorm.DB) *GormHistoryRepository {
	return &GormHistoryRepository{db: db}
}

type HistoryRepository interface {
	CreateHistory(history entities.History) error
	GetHistoryPerQuizGroupedByCategoryAndTimes(times int, kidID string) (map[string]map[int]entities.GroupedHistory, error)
	GetLatestHistoryPerEvaluate(times int, kidID string) ([]entities.History, error)
	GetLatestHistoryPerQuiz(times int, cate int, kidID string) ([]entities.History, error)
	GetHistoryResult(evaluatedTimes int, kidID string) (map[int]entities.GroupedHistory, error)
	DeleteHistoryWithTimes(evaluatedTimes int, kidID string, times int, cate int) error
}

func (r *GormHistoryRepository) CreateHistory(history entities.History) error {
	return r.db.Create(&history).Error
}

func (r *GormHistoryRepository) GetLatestHistoryPerQuiz(times int, cate int, kidID string) ([]entities.History, error) {
	var histories []entities.History
	subQuery := r.db.
		Table("histories").
		Select("histories.quiz_id, MAX(histories.created_at) AS max_created_at").
		Joins("JOIN quizzes ON quizzes.id = histories.quiz_id").
		Where("histories.kid_id = ? AND histories.evaluated_times = ? AND quizzes.category_id = ?", kidID, times, cate).
		Group("histories.quiz_id")

	err := r.db.
		Table("histories").
		Joins("JOIN quizzes ON quizzes.id = histories.quiz_id").
		Joins("JOIN (?) AS latest ON histories.quiz_id = latest.quiz_id AND histories.created_at = latest.max_created_at", subQuery).
		Where("histories.kid_id = ? AND histories.evaluated_times = ? AND quizzes.category_id = ?", kidID, times, cate).
		Preload("Quiz.Category").Preload("Quiz.Period").
		Find(&histories).Error

	if err != nil {
		return nil, err
	}

	return histories, nil
}

func (r *GormHistoryRepository) GetLatestHistoryPerEvaluate(times int, kidID string) ([]entities.History, error) {
	var histories []entities.History
	subQuery := r.db.
		Table("histories").
		Select("histories.quiz_id, MAX(histories.created_at) AS max_created_at").
		Joins("JOIN quizzes ON quizzes.id = histories.quiz_id").
		Where("histories.kid_id = ? AND histories.evaluated_times = ? ", kidID, times).
		Group("histories.quiz_id")

	err := r.db.
		Table("histories").
		Joins("JOIN quizzes ON quizzes.id = histories.quiz_id").
		Joins("JOIN (?) AS latest ON histories.quiz_id = latest.quiz_id AND histories.created_at = latest.max_created_at", subQuery).
		Where("histories.kid_id = ? AND histories.evaluated_times = ?", kidID, times).
		Preload("Quiz.Category").Preload("Quiz.Period").
		Find(&histories).Error

	if err != nil {
		return nil, err
	}

	return histories, nil
}

func (r *GormHistoryRepository) GetHistoryPerQuizGroupedByCategoryAndTimes(times int, kidID string) (map[string]map[int]entities.GroupedHistory, error) {
	var histories []entities.History

	if err := r.db.
		Joins("JOIN quizzes ON quizzes.id = histories.quiz_id").
		Joins("JOIN categories ON categories.id = quizzes.category_id").
		Preload("Quiz.Category").
		Where("histories.evaluated_times = ? AND histories.kid_id = ?", times, kidID).
		Find(&histories).Error; err != nil {
		return nil, err
	}

	result := make(map[string]map[int]entities.GroupedHistory)

	for _, h := range histories {
		categoryName := h.Quiz.Category.Category
		if _, ok := result[categoryName]; !ok {
			result[categoryName] = make(map[int]entities.GroupedHistory)
		}

		g := result[categoryName][h.Times]
		g.Histories = append(g.Histories, h)
		result[categoryName][h.Times] = g
	}

	for _, groupByTimes := range result {
		for times, g := range groupByTimes {
			if times == 0 {
				g.DoneAt = nil
			} else {
				var latest time.Time
				for _, h := range g.Histories {
					if h.CreatedAt.After(latest) {
						latest = h.CreatedAt
					}
				}
				g.DoneAt = &latest
			}

			allPassed := true
			hasDone := false
			for _, h := range g.Histories {
				if h.Times > 0 {
					hasDone = true
					if !h.Answer {
						allPassed = false
						break
					}
				}
			}

			if !hasDone {
				g.Solution = "รอประเมิน"
			} else if allPassed {
				g.Solution = "ผ่าน"
			} else {
				g.Solution = "ไม่ผ่าน"
			}

			groupByTimes[times] = g
		}
	}

	return result, nil
}

func (r *GormHistoryRepository) GetHistoryResult(evaluatedTimes int, kidID string) (map[int]entities.GroupedHistory, error) {
	var histories []entities.History

	if err := r.db.
		Joins("JOIN quizzes ON quizzes.id = histories.quiz_id").
		Joins("JOIN categories ON categories.id = quizzes.category_id").
		Preload("Quiz.Category").Preload("Quiz.Period").
		Where("histories.evaluated_times = ? AND histories.kid_id = ?", evaluatedTimes, kidID).
		Find(&histories).Error; err != nil {
		return nil, err
	}

	latestHistories := make(map[int]entities.History)
	for _, h := range histories {
		if existing, found := latestHistories[h.QuizID]; !found || h.Times > existing.Times || h.CreatedAt.After(existing.CreatedAt) {
			latestHistories[h.QuizID] = h
		}
	}

	grouped := make(map[int][]entities.History)
	for _, h := range latestHistories {
		categoryID := h.Quiz.CategoryID
		grouped[categoryID] = append(grouped[categoryID], h)
	}

	result := make(map[int]entities.GroupedHistory)
	for categoryID, hs := range grouped {
		group := entities.GroupedHistory{
			Histories: hs,
		}

		allPassed := true
		var latestTime *time.Time
		for i, h := range hs {
			if !h.Answer {
				allPassed = false
			}
			if i == 0 || h.CreatedAt.After(*latestTime) {
				t := h.CreatedAt
				latestTime = &t
			}
		}

		group.Solution = "ผ่าน"
		if !allPassed {
			group.Solution = "ไม่ผ่าน"
		}
		group.DoneAt = latestTime

		result[categoryID] = group
	}

	return result, nil
}

func (r *GormHistoryRepository) DeleteHistoryWithTimes(evaluatedTimes int, kidID string, times int, cate int) error {
	var quizIDs []int
	if err := r.db.Model(&entities.Quiz{}).
		Where("category_id = ? AND period_id = ?", cate, evaluatedTimes).
		Pluck("id", &quizIDs).Error; err != nil {
		return err
	}

	if len(quizIDs) == 0 {
		return nil
	}

	return r.db.
		Where("evaluated_times = ? AND kid_id = ? AND times = ? AND quiz_id IN ?", evaluatedTimes, kidID, times, quizIDs).
		Delete(&entities.History{}).Error
}
