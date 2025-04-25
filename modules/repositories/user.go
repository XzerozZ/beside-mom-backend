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
	GetMomByID(id string) (*entities.User, error)
	GetAllMom() ([]entities.User, error)

	CreateOTP(otp *entities.OTP) error
	GetOTPByUserID(userID string) (*entities.OTP, error)
	DeleteOTP(userID string) error
}

func (r *GormUserRepository) CreateUser(user *entities.User) (*entities.User, error) {
	if err := r.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return r.GetUserByID(user.ID)
}

func (r *GormUserRepository) GetUserByID(id string) (*entities.User, error) {
	var user entities.User
	if err := r.db.Preload("Role").Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *GormUserRepository) FindUserByEmail(email string) (entities.User, error) {
	var user entities.User
	if err := r.db.Preload("Role").Where("email = ?", email).First(&user).Error; err != nil {
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
	if err := r.db.Where("role_name = ?", name).First(&role).Error; err != nil {
		return entities.Role{}, err
	}

	return role, nil
}

func (r *GormUserRepository) GetMomByID(id string) (*entities.User, error) {
	var user entities.User
	if err := r.db.Preload("Role").Preload("Kid.Growth", func(db *gorm.DB) *gorm.DB {
		return db.Order("months")
	}).Where("id = ? AND role_id = ?", id, 2).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *GormUserRepository) GetAllMom() ([]entities.User, error) {
	var users []entities.User
	if err := r.db.Where("role_id = ?", 2).Preload("Role").Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *GormUserRepository) CreateOTP(otp *entities.OTP) error {
	if err := r.db.Create(otp).Error; err != nil {
		return err
	}

	return nil
}

func (r *GormUserRepository) GetOTPByUserID(userID string) (*entities.OTP, error) {
	var otp entities.OTP
	if err := r.db.Where("user_id = ?", userID).First(&otp).Error; err != nil {
		return nil, err
	}

	return &otp, nil
}

func (r *GormUserRepository) DeleteOTP(userID string) error {
	if err := r.db.Delete(&entities.OTP{}, "user_id = ?", userID).Error; err != nil {
		return err
	}

	return nil
}
