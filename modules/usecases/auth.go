package usecases

import (
	"Beside-Mom-BE/configs"
	"Beside-Mom-BE/modules/entities"
	"Beside-Mom-BE/modules/repositories"
	"Beside-Mom-BE/pkg/utils"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase interface {
	RegisterAdmin(user *entities.User) (*entities.User, error)
	Login(email, password string) (string, *entities.User, error)
	ForgotPassword(email string) error
	VerifyOTP(email, otpCode string) error
	ChangedPassword(email, newPassword string) error
}

type AuthUseCaseImpl struct {
	repo      repositories.UserRepository
	jwtSecret string
	mail      configs.Mail
}

func NewAuthUseCase(repo repositories.UserRepository, jwt configs.JWT, mail configs.Mail) *AuthUseCaseImpl {
	return &AuthUseCaseImpl{
		repo:      repo,
		jwtSecret: jwt.Secret,
		mail:      mail,
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

func (u *AuthUseCaseImpl) ForgotPassword(email string) error {
	user, err := u.repo.FindUserByEmail(email)
	if err != nil {
		return errors.New("invalid email")
	}

	userID := user.ID
	otpCode, err := utils.GenerateRandomOTP(6, true)
	if err != nil {
		return err
	}

	expiresAt := time.Now().Add(5 * time.Minute)
	otp, err := u.repo.GetOTPByUserID(userID)
	if err == nil && otp != nil {
		if err := u.repo.DeleteOTP(userID); err != nil {
			return err
		}
	}

	newOTP := &entities.OTP{
		UserID:    userID,
		OTP:       otpCode,
		ExpiresAt: expiresAt,
	}

	if err := u.repo.CreateOTP(newOTP); err != nil {
		return err
	}

	if err := utils.SendMail("./assets/OTPmail.html", user, otpCode, u.mail); err != nil {
		return err
	}

	return nil
}

func (u *AuthUseCaseImpl) VerifyOTP(email, otpCode string) error {
	user, err := u.repo.FindUserByEmail(email)
	if err != nil {
		return err
	}

	otp, err := u.repo.GetOTPByUserID(user.ID)
	if err != nil {
		return err
	}

	if time.Now().After(otp.ExpiresAt) {
		return errors.New("OTP is expired")
	}

	if otp.OTP != otpCode {
		return errors.New("OTP is incorrect")
	}

	if err := u.repo.DeleteOTP(user.ID); err != nil {
		return err
	}

	return nil
}

func (u *AuthUseCaseImpl) ChangedPassword(email, newPassword string) error {
	user, err := u.repo.FindUserByEmail(email)
	if err != nil {
		return err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(newPassword)); err == nil {
		return errors.New("new password cannot be the same as the old password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	_, err = u.repo.UpdateUserByID(&user)
	if err != nil {
		return err
	}

	return nil
}
