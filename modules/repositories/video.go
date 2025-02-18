package repositories

import (
	"Beside-Mom-BE/modules/entities"

	"gorm.io/gorm"
)

type GormVideoRepository struct {
	db *gorm.DB
}

func NewGormVideoRepository(db *gorm.DB) *GormVideoRepository {
	return &GormVideoRepository{db: db}
}

type VideoRepository interface {
	CreateVideo(video *entities.Video) (*entities.Video, error)
	GetVideoByID(id string) (*entities.Video, error)
	GetAllVideo() ([]entities.Video, error)
	UpdateVideoByID(video *entities.Video) (*entities.Video, error)
	DeleteVideoByID(id string) error
}

func (r *GormVideoRepository) CreateVideo(video *entities.Video) (*entities.Video, error) {
	if err := r.db.Create(&video).Error; err != nil {
		return nil, err
	}

	return r.GetVideoByID(video.ID)
}

func (r *GormVideoRepository) GetVideoByID(id string) (*entities.Video, error) {
	var video entities.Video
	if err := r.db.First(&video, id).Error; err != nil {
		return nil, err
	}

	return &video, nil
}

func (r *GormVideoRepository) GetAllVideo() ([]entities.Video, error) {
	var video []entities.Video
	if err := r.db.Where("id != ?", "00001").Find(&video).Error; err != nil {
		return nil, err
	}

	return video, nil
}

func (r *GormVideoRepository) UpdateVideoByID(video *entities.Video) (*entities.Video, error) {
	if err := r.db.Save(&video).Error; err != nil {
		return nil, err
	}

	return r.GetVideoByID(video.ID)
}

func (r *GormVideoRepository) DeleteVideoByID(id string) error {
	return r.db.Where("id = ?", id).Delete(&entities.Video{}).Error
}
