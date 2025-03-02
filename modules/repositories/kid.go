package repositories

import (
	"Beside-Mom-BE/modules/entities"

	"gorm.io/gorm"
)

type GormKidsRepository struct {
	db *gorm.DB
}

func NewGormKidsRepository(db *gorm.DB) *GormKidsRepository {
	return &GormKidsRepository{db: db}
}

type KidsRepository interface {
	CreateKid(kid *entities.Kid) (*entities.Kid, error)
	GetKidByID(id string) (*entities.Kid, error)
	GetKidByUserID(userID string) ([]entities.Kid, error)
	UpdateKidByID(kid *entities.Kid) (*entities.Kid, error)
	DeleteKidByID(id string) error
}

func (r *GormKidsRepository) CreateKid(kid *entities.Kid) (*entities.Kid, error) {
	if err := r.db.Create(&kid).Error; err != nil {
		return nil, err
	}

	return r.GetKidByID(kid.ID)
}

func (r *GormKidsRepository) GetKidByID(id string) (*entities.Kid, error) {
	var kid entities.Kid
	if err := r.db.Where("id = ?", id).First(&kid).Error; err != nil {
		return nil, err
	}

	return &kid, nil
}

func (r *GormKidsRepository) GetKidByUserID(userID string) ([]entities.Kid, error) {
	var kids []entities.Kid
	if err := r.db.Where("user_id = ?", userID).Find(&kids).Error; err != nil {
		return nil, err
	}

	return kids, nil
}

func (r *GormKidsRepository) UpdateKidByID(kid *entities.Kid) (*entities.Kid, error) {
	if err := r.db.Save(&kid).Error; err != nil {
		return nil, err
	}

	return r.GetKidByID(kid.ID)
}

func (r *GormKidsRepository) DeleteKidByID(id string) error {
	return r.db.Where("id = ?", id).Delete(&entities.Kid{}).Error
}
