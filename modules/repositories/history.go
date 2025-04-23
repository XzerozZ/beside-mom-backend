package repositories

import (
	"Beside-Mom-BE/modules/entities"

	"gorm.io/gorm"
)

type GormHistoryRepository struct {
	db *gorm.DB
}

func NewGormHistoryRepository(db *gorm.DB) *GormHistoryRepository {
	return &GormHistoryRepository{db: db}
}

type HistoryRepository interface {
	CreateHisotry(history entities.History) error
	GetHistoryOfEvaluate(times int, kidID string) ([]entities.History, error)
	GetLatestHistoryPerQuiz(times int, kidID string) ([]entities.History, error)
	DeleteHistoryWithTimes(evaluatedTimes int, kidID string, times int) error
}

func (r *GormHistoryRepository) CreateHisotry(history entities.History) error {
	return r.db.Create(&history).Error
}

func (r *GormHistoryRepository) GetHistoryOfEvaluate(times int, kidID string) ([]entities.History, error) {
	var histories []entities.History
	if err := r.db.Where("evaluated_times = ? AND kid_id = ?", times, kidID).Find(&histories).Error; err != nil {
		return nil, err
	}

	return histories, nil
}

func (r *GormHistoryRepository) GetLatestHistoryPerQuiz(times int, kidID string) ([]entities.History, error) {
	var histories []entities.History
	var quizCount int64
	if err := r.db.Model(&entities.Quiz{}).Count(&quizCount).Error; err != nil {
		return nil, err
	}

	subQuery := r.db.
		Table("histories").
		Select("quiz_id, MAX(created_at) AS max_created_at").
		Where("kid_id = ? AND evaluated_times = ?", kidID, times).
		Group("quiz_id")

	if err := r.db.
		Table("histories").
		Joins("JOIN (?) AS latest ON histories.quiz_id = latest.quiz_id AND histories.created_at = latest.max_created_at", subQuery).
		Where("histories.kid_id = ? AND histories.evaluated_times = ?", kidID, times).
		Limit(int(quizCount)).
		Find(&histories).Error; err != nil {
		return nil, err
	}

	return histories, nil
}

func (r *GormHistoryRepository) DeleteHistoryWithTimes(evaluatedTimes int, kidID string, times int) error {
	return r.db.Where("evaluated_times = ? AND kid_id = ? AND times = ?", evaluatedTimes, kidID, times).Delete(&entities.History{}).Error
}
