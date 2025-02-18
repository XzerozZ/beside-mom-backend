package repositories

import (
	"Beside-Mom-BE/modules/entities"

	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

type UserRepository interface {
	CreateUser(user *entities.User) (*entities.User, error)
	GetUserByID(id string) (*entities.User, error)
	FindUserByEmail(email string) (entities.User, error)
	UpdateUserByID(user *entities.User) (*entities.User, error)
	DeleteUserByID(userID string) error

	GetRoleByName(name string) (entities.Role, error)
}

func (r *GormUserRepository) CreateUser(user *entities.User) (*entities.User, error) {
	if err := r.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return r.GetUserByID(user.ID)
}

func (r *GormUserRepository) GetUserByID(id string) (*entities.User, error) {
	var user entities.User
	err := r.db.Preload("Role").Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *GormUserRepository) FindUserByEmail(email string) (entities.User, error) {
	var user entities.User
	err := r.db.Preload("Role").Where("email = ?", email).First(&user).Error
	if err != nil {
		return entities.User{}, err
	}

	return user, nil
}

func (r *GormUserRepository) UpdateUserByID(user *entities.User) (*entities.User, error) {
	if err := r.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return r.GetUserByID(user.ID)
}

func (r *GormUserRepository) DeleteUserByID(userID string) error {
	return r.db.Where("id = ?", userID).Delete(&entities.User{}).Error
}

func (r *GormUserRepository) GetRoleByName(name string) (entities.Role, error) {
	var role entities.Role
	err := r.db.Where("role_name = ?", name).First(&role).Error
	if err != nil {
		return entities.Role{}, err
	}

	return role, nil
}
