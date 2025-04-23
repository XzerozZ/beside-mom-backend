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
	GetAllEvaluate() ([]entities.Evaluate, error)
	UpdateEvaluate(evaluatedTimes int, kidID string, solution string) error
}

func (r *GormEvaluateRepository) GetEvaluateByID(id string) (*entities.Evaluate, error) {
	var eva entities.Evaluate
	if err := r.db.First(&eva, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &eva, nil
}

func (r *GormEvaluateRepository) GetAllEvaluate() ([]entities.Evaluate, error) {
	var evas []entities.Evaluate
	if err := r.db.Find(&evas).Error; err != nil {
		return nil, err
	}

	return evas, nil
}

func (r *GormEvaluateRepository) UpdateEvaluate(evaluatedTimes int, kidID string, solution string) error {
	return r.db.Model(&entities.Evaluate{}).Where("evaluated_times = ? AND kid_id = ?", evaluatedTimes, kidID).
		Updates(map[string]interface{}{
			"solution":     solution,
			"status":       true,
			"completed_at": time.Now(),
		}).Error
}
