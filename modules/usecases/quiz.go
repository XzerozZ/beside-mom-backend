package usecases

import (
	"Beside-Mom-BE/configs"
	"Beside-Mom-BE/modules/entities"
	"Beside-Mom-BE/modules/repositories"
	"Beside-Mom-BE/pkg/utils"
	"mime/multipart"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type QuizUseCase interface {
	CreateQuiz(quiz *entities.Quiz, banner *multipart.FileHeader, ctx *fiber.Ctx) (*entities.Quiz, error)
	GetQuizByID(id int) (*entities.Quiz, error)
	GetAllQuiz() ([]entities.Quiz, error)
	UpdateQuizByID(id int, quiz *entities.Quiz, banner *multipart.FileHeader, ctx *fiber.Ctx) (*entities.Quiz, error)
	DeleteQuizByID(id int) error
}

type QuizUseCaseImpl struct {
	repo repositories.QuizRepository
	supa configs.Supabase
}

func NewQuizUseCase(repo repositories.QuizRepository, supa configs.Supabase) *QuizUseCaseImpl {
	return &QuizUseCaseImpl{
		repo: repo,
		supa: supa,
	}
}

func (u *QuizUseCaseImpl) CreateQuiz(quiz *entities.Quiz, banner *multipart.FileHeader, ctx *fiber.Ctx) (*entities.Quiz, error) {
	if banner != nil {
		fileName := uuid.New().String() + "_title.jpg"
		if err := ctx.SaveFile(banner, "./uploads/"+fileName); err != nil {
			return nil, err
		}

		imageUrl, err := utils.UploadImage(fileName, "", u.supa)
		if err != nil {
			os.Remove("./uploads/" + fileName)
			return nil, err
		}

		if err := os.Remove("./uploads/" + fileName); err != nil {
			return nil, err
		}

		quiz.Banner = imageUrl
	}

	return u.repo.CreateQuiz(quiz)
}

func (u *QuizUseCaseImpl) GetQuizByID(id int) (*entities.Quiz, error) {
	return u.repo.GetQuizByID(id)
}

func (u *QuizUseCaseImpl) GetAllQuiz() ([]entities.Quiz, error) {
	return u.repo.GetAllQuiz()
}

func (u *QuizUseCaseImpl) UpdateQuizByID(id int, quiz *entities.Quiz, banner *multipart.FileHeader, ctx *fiber.Ctx) (*entities.Quiz, error) {
	existingQuiz, err := u.repo.GetQuizByID(id)
	if err != nil {
		return nil, err
	}

	existingQuiz.Question = quiz.Question
	if banner != nil {
		fileName := uuid.New().String() + "_title.jpg"
		if err := ctx.SaveFile(banner, "./uploads/"+fileName); err != nil {
			return nil, err
		}

		imageUrl, err := utils.UploadImage(fileName, "", u.supa)
		if err != nil {
			os.Remove("./uploads/" + fileName)
			return nil, err
		}

		if err := os.Remove("./uploads/" + fileName); err != nil {
			return nil, err
		}

		existingQuiz.Banner = imageUrl
	}

	updatedQuiz, err := u.repo.UpdateQuizByID(existingQuiz)
	if err != nil {
		return nil, err
	}

	return updatedQuiz, nil
}

func (u *QuizUseCaseImpl) DeleteQuizByID(id int) error {
	return u.repo.DeleteQuizByID(id)
}
