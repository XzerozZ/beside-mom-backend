package repositories

import (
	"Beside-Mom-BE/modules/entities"

	"gorm.io/gorm"
)

type GormQuestionRepository struct {
	db *gorm.DB
}

func NewGormQuestionRepository(db *gorm.DB) *GormQuestionRepository {
	return &GormQuestionRepository{db: db}
}

type QuestionRepository interface {
	CreateQuestion(question *entities.Question) (*entities.Question, error)
	GetQuestionByID(id string) (*entities.Question, error)
	GetAllQuestion() ([]entities.Question, error)
	UpdateQuestionByID(question *entities.Question) (*entities.Question, error)
	DeleteQuestionByID(id string) error
}

func (r *GormQuestionRepository) CreateQuestion(question *entities.Question) (*entities.Question, error) {
	if err := r.db.Create(&question).Error; err != nil {
		return nil, err
	}

	return r.GetQuestionByID(question.ID)
}

func (r *GormQuestionRepository) GetQuestionByID(id string) (*entities.Question, error) {
	var question entities.Question
	if err := r.db.First(&question, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &question, nil
}

func (r *GormQuestionRepository) GetAllQuestion() ([]entities.Question, error) {
	var questions []entities.Question
	if err := r.db.Find(&questions).Error; err != nil {
		return nil, err
	}

	return questions, nil
}

func (r *GormQuestionRepository) UpdateQuestionByID(question *entities.Question) (*entities.Question, error) {
	if err := r.db.Save(&question).Error; err != nil {
		return nil, err
	}

	return r.GetQuestionByID(question.ID)
}

func (r *GormQuestionRepository) DeleteQuestionByID(id string) error {
	return r.db.Where("id = ?", id).Delete(&entities.Question{}).Error
}
