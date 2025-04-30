package repositories

import (
	"Beside-Mom-BE/modules/entities"
	"time"

	"gorm.io/gorm"
)

type GormEvaluateRepository struct {
	db *gorm.DB
}

func NewGormEvaluateRepository(db *gorm.DB) *GormEvaluateRepository {
	return &GormEvaluateRepository{db: db}
}

type EvaluateRepository interface {
	GetEvaluateByID(id string) (*entities.Evaluate, error)
	GetAllEvaluate(kidID string) ([]entities.Evaluate, error)
	UpdateEvaluate(evaluatedTimes int, kidID string, solution string, status bool) error
}

func (r *GormEvaluateRepository) GetEvaluateByID(id string) (*entities.Evaluate, error) {
	var eva entities.Evaluate
	if err := r.db.First(&eva, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &eva, nil
}

func (r *GormEvaluateRepository) GetAllEvaluate(kidID string) ([]entities.Evaluate, error) {
	var evaluate []entities.Evaluate
	if err := r.db.Where("kid_id = ?", kidID).Order("evaluated_times").Find(&evaluate).Error; err != nil {
		return nil, err
	}

	return evaluate, nil
}

func (r *GormEvaluateRepository) UpdateEvaluate(evaluatedTimes int, kidID string, solution string, status bool) error {
	return r.db.Model(&entities.Evaluate{}).Where("evaluated_times = ? AND kid_id = ?", evaluatedTimes, kidID).
		Updates(map[string]interface{}{
			"solution":     solution,
			"status":       status,
			"completed_at": time.Now(),
		}).Error
}
