package usecases

import (
	"Beside-Mom-BE/configs"
	"Beside-Mom-BE/modules/entities"
	"Beside-Mom-BE/modules/repositories"
	"Beside-Mom-BE/pkg/utils"
	"errors"
	"mime/multipart"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	CreateUser(user *entities.User, image *multipart.FileHeader, ctx *fiber.Ctx) (*entities.User, error)
	GetMomByID(id string) (*entities.User, error)
	GetAllMom() ([]entities.User, error)
	UpdateUserByID(id string, image *multipart.FileHeader, ctx *fiber.Ctx) (*entities.User, error)
	DeleteUser(id string) error
}

type UserUseCaseImpl struct {
	repo repositories.UserRepository
	supa configs.Supabase
	mail configs.Mail
}

func NewUserUseCase(repo repositories.UserRepository, supa configs.Supabase, mail configs.Mail) *UserUseCaseImpl {
	return &UserUseCaseImpl{
		repo: repo,
		supa: supa,
		mail: mail,
	}
}

func (u *UserUseCaseImpl) CreateUser(user *entities.User, image *multipart.FileHeader, ctx *fiber.Ctx) (*entities.User, error) {
	normalizedEmail, err := utils.NormalizeEmail(user.Email)
	if err != nil {
		return nil, errors.New("invalid email format")
	}

	user.Email = normalizedEmail
	if _, err := u.repo.FindUserByEmail(user.Email); err == nil {
		return nil, errors.New("this email already have account")
	}

	role, err := u.repo.GetRoleByName("User")
	if err != nil {
		return nil, errors.New("role not found")
	}

	user.RoleID = role.ID
	password, err := utils.GeneratePassword(8)
	if err != nil {
		return nil, errors.New("can't generate password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user.Password = string(hashedPassword)
	if image != nil {
		fileName := uuid.New().String() + "_title.jpg"
		if err := ctx.SaveFile(image, "./uploads/"+fileName); err != nil {
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

		user.ImageLink = imageUrl
	}

	createdUser, err := u.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	if err := utils.SendPasswordMail("./assets/Passwordmail.html", *createdUser, password, u.mail); err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (u *UserUseCaseImpl) GetMomByID(id string) (*entities.User, error) {
	return u.repo.GetMomByID(id)
}

func (u *UserUseCaseImpl) GetAllMom() ([]entities.User, error) {
	return u.repo.GetAllMom()
}

func (u *UserUseCaseImpl) UpdateUserByID(id string, image *multipart.FileHeader, ctx *fiber.Ctx) (*entities.User, error) {
	existingUser, err := u.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	if image != nil {
		fileName := uuid.New().String() + "_title.jpg"
		if err := ctx.SaveFile(image, "./uploads/"+fileName); err != nil {
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

		existingUser.ImageLink = imageUrl
	}

	updatedUser, err := u.repo.UpdateUserByID(existingUser)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (u *UserUseCaseImpl) DeleteUser(id string) error {
	return u.repo.DeleteUser(id)
}
