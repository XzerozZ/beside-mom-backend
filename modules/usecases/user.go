package usecases

import (
	"Beside-Mom-BE/configs"
	"Beside-Mom-BE/modules/entities"
	"Beside-Mom-BE/modules/repositories"
	"Beside-Mom-BE/pkg/utils"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	CreateUser(user *entities.User, image *multipart.FileHeader, ctx *fiber.Ctx) (*entities.User, error)
	Chat(meassage string) (map[string]interface{}, error)
	GetMomByID(id string) (*entities.User, error)
	GetAllMom() ([]entities.User, error)
	UpdateUserByIDForUser(id string, image *multipart.FileHeader, ctx *fiber.Ctx) (*entities.User, error)
	UpdateUserByIDForAdmin(id string, user *entities.User, image *multipart.FileHeader, ctx *fiber.Ctx) (*entities.User, error)
	DeleteUser(id string) error
}

type UserUseCaseImpl struct {
	repo repositories.UserRepository
	supa configs.Supabase
	mail configs.Mail
	chat configs.Chat
}

func NewUserUseCase(repo repositories.UserRepository, supa configs.Supabase, mail configs.Mail, chat configs.Chat) *UserUseCaseImpl {
	return &UserUseCaseImpl{
		repo: repo,
		supa: supa,
		mail: mail,
		chat: chat,
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
	if user.PID == "" {
		latestPID, err := u.repo.FindLatestUnnamedPID()
		if err != nil {
			return nil, err
		}

		seq := 0
		if latestPID != "" {
			if _, err := fmt.Sscanf(latestPID, "Unnamed-Case-%03d", &seq); err == nil {
				seq++
			} else {
				seq = 1
			}
		}

		seq += 1
		user.PID = fmt.Sprintf("Unnamed-Case-%03d", seq)
	}

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

func (u *UserUseCaseImpl) UpdateUserByIDForUser(id string, image *multipart.FileHeader, ctx *fiber.Ctx) (*entities.User, error) {
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

		if err := utils.DeleteImage(existingUser.ImageLink, u.supa); err != nil {
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

func (u *UserUseCaseImpl) UpdateUserByIDForAdmin(id string, user *entities.User, image *multipart.FileHeader, ctx *fiber.Ctx) (*entities.User, error) {
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

		if err := utils.DeleteImage(existingUser.ImageLink, u.supa); err != nil {
			return nil, err
		}

		existingUser.ImageLink = imageUrl
	}

	existingUser.Email = user.Email
	existingUser.PID = user.PID
	existingUser.Firstname = user.Firstname
	existingUser.Lastname = user.Lastname
	updatedUser, err := u.repo.UpdateUserByID(existingUser)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (u *UserUseCaseImpl) DeleteUser(id string) error {
	existingUser, err := u.repo.GetUserByID(id)
	if err != nil {
		return err
	}

	if err := utils.DeleteImage(existingUser.ImageLink, u.supa); err != nil {
		return err
	}

	return u.repo.DeleteUser(id)
}

func (u *UserUseCaseImpl) Chat(message string) (map[string]interface{}, error) {
	requestBody := map[string]interface{}{
		"message":    message,
		"max_tokens": 512,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	urlStr := fmt.Sprintf("%s/chat", u.chat.URL)
	resp, err := http.Post(urlStr, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get response: %s | Body: %s", resp.Status, string(body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	responseText, ok := result["response"].(string)
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	finalResult := map[string]interface{}{
		"sender":   "chat",
		"response": responseText,
		"sent_at":  time.Now(),
	}

	return finalResult, nil
}
