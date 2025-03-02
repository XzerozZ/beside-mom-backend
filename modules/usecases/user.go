package usecases

import (
	"Beside-Mom-BE/configs"
	"Beside-Mom-BE/modules/entities"
	"Beside-Mom-BE/modules/repositories"
	"Beside-Mom-BE/pkg/utils"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	CreateUser(user *entities.User) (*entities.User, error)
	GetMomByID(id string) (*entities.User, error)
	GetAllMom() ([]entities.User, error)
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

func (u *UserUseCaseImpl) CreateUser(user *entities.User) (*entities.User, error) {
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
