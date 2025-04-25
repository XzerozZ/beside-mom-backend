package repositories

import (
	"Beside-Mom-BE/modules/entities"

	"gorm.io/gorm"
)

type GormAppRepository struct {
	db *gorm.DB
}

func NewGormAppRepository(db *gorm.DB) *GormAppRepository {
	return &GormAppRepository{db: db}
}

type AppRepository interface {
	CreateAppointment(app *entities.Appointment) (*entities.Appointment, error)
	GetAppByID(id string) (*entities.Appointment, error)
	GetAppByUserID(userID string) ([]entities.Appointment, error)
	GetAppInProgressByUserID(userID string) ([]entities.Appointment, error)
	GetAllApp() ([]entities.Appointment, error)
	UpdateAppByID(app *entities.Appointment) (*entities.Appointment, error)
	DeleteAppByID(id string) error
}

func (r *GormAppRepository) CreateAppointment(app *entities.Appointment) (*entities.Appointment, error) {
	if err := r.db.Create(&app).Error; err != nil {
		return nil, err
	}

	return r.GetAppByID(app.ID)
}

func (r *GormAppRepository) GetAppByID(id string) (*entities.Appointment, error) {
	var app entities.Appointment
	if err := r.db.Preload("User").Where("id = ?", id).First(&app).Error; err != nil {
		return nil, err
	}

	return &app, nil
}

func (r *GormAppRepository) GetAppByUserID(userID string) ([]entities.Appointment, error) {
	var apps []entities.Appointment
	if err := r.db.Preload("User").Where("user_id = ?", userID).Find(&apps).Error; err != nil {
		return nil, err
	}

	return apps, nil
}

func (r *GormAppRepository) GetAllApp() ([]entities.Appointment, error) {
	var apps []entities.Appointment
	if err := r.db.Preload("User").Find(&apps).Error; err != nil {
		return nil, err
	}

	return apps, nil
}

func (r *GormAppRepository) GetAppInProgressByUserID(userID string) ([]entities.Appointment, error) {
	var apps []entities.Appointment
	if err := r.db.Preload("User").Where("status = ? ", 1).Find(&apps).Error; err != nil {
		return nil, err
	}

	return apps, nil
}

func (r *GormAppRepository) UpdateAppByID(app *entities.Appointment) (*entities.Appointment, error) {
	if err := r.db.Save(&app).Error; err != nil {
		return nil, err
	}

	return r.GetAppByID(app.ID)
}

func (r *GormAppRepository) DeleteAppByID(id string) error {
	return r.db.Where("id = ?", id).Delete(&entities.Appointment{}).Error
}
