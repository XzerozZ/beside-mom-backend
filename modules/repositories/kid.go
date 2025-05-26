package repositories

import (
	"Beside-Mom-BE/modules/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormKidsRepository struct {
	db *gorm.DB
}

func NewGormKidsRepository(db *gorm.DB) *GormKidsRepository {
	return &GormKidsRepository{db: db}
}

type KidsRepository interface {
	CreateKid(kid *entities.Kid) (*entities.Kid, error)
	GetKidByID(id string) (*entities.Kid, error)
	GetKidByIDForUser(id string) (*entities.Kid, error)
	UpdateKidByID(kid *entities.Kid) (*entities.Kid, error)
	DeleteKidByID(id string) error
}

func (r *GormKidsRepository) CreateKid(kid *entities.Kid) (*entities.Kid, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&kid).Error; err != nil {
			return err
		}

		var quizzes []entities.Quiz
		if err := tx.Find(&quizzes).Error; err != nil {
			return err
		}

		var periods []entities.Period
		if err := tx.Find(&periods).Error; err != nil {
			return err
		}

		for i, period := range periods {
			eval := entities.Evaluate{
				ID:             uuid.New().String(),
				Status:         false,
				Solution:       "รอประเมิน",
				EvaluatedTimes: i + 1,
				PeriodID:       period.ID,
				KidID:          kid.ID,
			}

			if err := tx.Create(&eval).Error; err != nil {
				return err
			}

			for _, quiz := range quizzes {
				if quiz.PeriodID != period.ID {
					continue
				}

				history := entities.History{
					ID:             uuid.New().String(),
					QuizID:         quiz.ID,
					Answer:         false,
					Status:         false,
					EvaluatedTimes: eval.EvaluatedTimes,
					Times:          0,
					KidID:          kid.ID,
				}

				if err := tx.Create(&history).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return r.GetKidByID(kid.ID)
}

func (r *GormKidsRepository) GetKidByID(id string) (*entities.Kid, error) {
	var kid entities.Kid
	if err := r.db.Preload("Growth", func(db *gorm.DB) *gorm.DB {
		return db.Order("months")
	}).Where("id = ?", id).First(&kid).Error; err != nil {
		return nil, err
	}

	return &kid, nil
}

func (r *GormKidsRepository) GetKidByIDForUser(id string) (*entities.Kid, error) {
	var kid entities.Kid
	if err := r.db.Where("id = ?", id).First(&kid).Error; err != nil {
		return nil, err
	}

	var growth *entities.Growth
	if err := r.db.Where("kid_id = ?", kid.ID).Order("created_at desc").First(&growth).Error; err != nil {
		return nil, err
	}

	kid.BirthLength = growth.Length
	kid.BirthWeight = growth.Weight
	return &kid, nil
}

func (r *GormKidsRepository) UpdateKidByID(kid *entities.Kid) (*entities.Kid, error) {
	if err := r.db.Save(&kid).Error; err != nil {
		return nil, err
	}

	return r.GetKidByID(kid.ID)
}

func (r *GormKidsRepository) DeleteKidByID(id string) error {
	return r.db.Where("id = ?", id).Delete(&entities.Kid{}).Error
}
