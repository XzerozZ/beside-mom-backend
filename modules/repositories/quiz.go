package repositories

import (
	"Beside-Mom-BE/modules/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormQuizRepository struct {
	db *gorm.DB
}

func NewGormQuizRepository(db *gorm.DB) *GormQuizRepository {
	return &GormQuizRepository{db: db}
}

type QuizRepository interface {
	CreateQuiz(quiz *entities.Quiz) (*entities.Quiz, error)
	GetQuizByID(id int) (*entities.Quiz, error)
	GetAllQuiz() ([]entities.Quiz, error)
	GetQuizByIDandPeriod(id int, period int, cate int) (*entities.Quiz, error)
	GetQuizByCategoryandPeriod(period int, cate int) ([]entities.Quiz, error)
	UpdateQuizByID(quiz *entities.Quiz) (*entities.Quiz, error)
	DeleteQuizByID(id int) error
}

func (r *GormQuizRepository) CreateQuiz(quiz *entities.Quiz) (*entities.Quiz, error) {
	if err := r.db.Create(&quiz).Error; err != nil {
		return nil, err
	}

	var evaluates []entities.Evaluate
	if err := r.db.Where("status = ? AND solution = ? AND period_id = ?", false, "รอประเมิน", quiz.PeriodID).Find(&evaluates).Error; err != nil {
		return nil, err
	}

	for _, eval := range evaluates {
		history := entities.History{
			ID:             uuid.New().String(),
			QuizID:         quiz.ID,
			Answer:         false,
			Status:         false,
			EvaluatedTimes: eval.EvaluatedTimes,
			Times:          0,
			KidID:          eval.KidID,
		}

		if err := r.db.Create(&history).Error; err != nil {
			return nil, err
		}
	}

	return r.GetQuizByID(quiz.ID)
}

func (r *GormQuizRepository) GetQuizByID(id int) (*entities.Quiz, error) {
	var quiz entities.Quiz
	if err := r.db.Preload("Category").Preload("Period").First(&quiz, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &quiz, nil
}

func (r *GormQuizRepository) GetAllQuiz() ([]entities.Quiz, error) {
	var quizs []entities.Quiz
	if err := r.db.Preload("Category").Preload("Period").Order("id").Find(&quizs).Error; err != nil {
		return nil, err
	}

	return quizs, nil
}

func (r *GormQuizRepository) GetQuizByIDandPeriod(id int, period int, cate int) (*entities.Quiz, error) {
	var quiz *entities.Quiz
	if err := r.db.Where("period_id = ? AND category_id = ?", period, cate).Order("id").Preload("Category").Preload("Period").First(&quiz, id).Error; err != nil {
		return nil, err
	}

	return quiz, nil
}

func (r *GormQuizRepository) GetQuizByCategoryandPeriod(period int, cate int) ([]entities.Quiz, error) {
	var quiz []entities.Quiz
	if err := r.db.Where("period_id = ? AND category_id = ?", period, cate).Order("id").Preload("Category").Preload("Period").Find(&quiz).Error; err != nil {
		return nil, err
	}

	return quiz, nil
}

func (r *GormQuizRepository) UpdateQuizByID(quiz *entities.Quiz) (*entities.Quiz, error) {
	if err := r.db.
		Model(&entities.Quiz{}).
		Where("id = ?", quiz.ID).
		Updates(map[string]interface{}{
			"question":    quiz.Question,
			"description": quiz.Description,
			"solution":    quiz.Solution,
			"suggestion":  quiz.Suggestion,
			"category_id": quiz.CategoryID,
			"period_id":   quiz.PeriodID,
			"banner":      quiz.Banner,
		}).Error; err != nil {
		return nil, err
	}

	return r.GetQuizByID(quiz.ID)
}

func (r *GormQuizRepository) DeleteQuizByID(id int) error {
	return r.db.Where("id = ?", id).Delete(&entities.Quiz{}).Error
}
