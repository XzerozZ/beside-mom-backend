package usecases

import (
	"Beside-Mom-BE/modules/entities"
	"Beside-Mom-BE/modules/repositories"

	"github.com/google/uuid"
)

type QuestionUseCase interface {
	CreateQuestion(question *entities.Question) (*entities.Question, error)
	GetQuestionByID(id string) (*entities.Question, error)
	GetAllQuestion() ([]entities.Question, error)
	UpdateQuestionByID(id string, question *entities.Question) (*entities.Question, error)
	DeleteQuestionByID(id string) error
}

type QuestionUseCaseImpl struct {
	repo repositories.QuestionRepository
}

func NewQuestionUseCase(repo repositories.QuestionRepository) *QuestionUseCaseImpl {
	return &QuestionUseCaseImpl{repo: repo}
}

func (u *QuestionUseCaseImpl) CreateQuestion(question *entities.Question) (*entities.Question, error) {
	question.ID = uuid.New().String()
	return u.repo.CreateQuestion(question)
}

func (u *QuestionUseCaseImpl) GetQuestionByID(id string) (*entities.Question, error) {
	return u.repo.GetQuestionByID(id)
}

func (u *QuestionUseCaseImpl) GetAllQuestion() ([]entities.Question, error) {
	return u.repo.GetAllQuestion()
}

func (u *QuestionUseCaseImpl) UpdateQuestionByID(id string, question *entities.Question) (*entities.Question, error) {
	existingQuestion, err := u.repo.GetQuestionByID(id)
	if err != nil {
		return nil, err
	}

	existingQuestion.Question = question.Question
	existingQuestion.Answer = question.Answer
	updatedQuestion, err := u.repo.UpdateQuestionByID(existingQuestion)
	if err != nil {
		return nil, err
	}

	return updatedQuestion, nil
}

func (u *QuestionUseCaseImpl) DeleteQuestionByID(id string) error {
	return u.repo.DeleteQuestionByID(id)
}
