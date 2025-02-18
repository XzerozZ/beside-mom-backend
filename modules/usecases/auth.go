package usecases

import (
	"Beside-Mom-BE/configs"
	"Beside-Mom-BE/modules/entities"
	"Beside-Mom-BE/modules/repositories"
	"Beside-Mom-BE/pkg/utils"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase interface {
	RegisterAdmin(user *entities.User) (*entities.User, error)
	Login(email, password string) (string, *entities.User, error)
}

type AuthUseCaseImpl struct {
	repo      repositories.UserRepository
	jwtSecret string
}

func NewAuthUseCase(repo repositories.UserRepository, jwt configs.JWT) *AuthUseCaseImpl {
	return &AuthUseCaseImpl{
		repo:      repo,
		jwtSecret: jwt.Secret,
	}
}

func (u *AuthUseCaseImpl) RegisterAdmin(user *entities.User) (*entities.User, error) {
	normalizedEmail, err := utils.NormalizeEmail(user.Email)
	if err != nil {
		return nil, errors.New("invalid email format")
	}

	user.Email = normalizedEmail
	if _, err := u.repo.FindUserByEmail(user.Email); err == nil {
		return nil, errors.New("this email already have account")
	}

	role, err := u.repo.GetRoleByName("Admin")
	if err != nil {
		return nil, errors.New("role not found")
	}

	user.ID = uuid.New().String()
	user.RoleID = role.ID
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user.Password = string(hashedPassword)
	createdUser, err := u.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (u *AuthUseCaseImpl) Login(email, password string) (string, *entities.User, error) {
	normalizedEmail, err := utils.NormalizeEmail(email)
	if err != nil {
		return "", nil, errors.New("invalid email format")
	}

	email = normalizedEmail
	user, err := u.repo.FindUserByEmail(email)
	if err != nil {
		return "", nil, errors.New("invalid email")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, errors.New("invalid password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role.RoleName,
	})

	tokenString, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return "", nil, err
	}

	return tokenString, &user, nil
}
