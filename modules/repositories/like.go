package repositories

import (
	"Beside-Mom-BE/modules/entities"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type GormLikesRepository struct {
	db *gorm.DB
}

func NewGormLikesRepository(db *gorm.DB) *GormLikesRepository {
	return &GormLikesRepository{db: db}
}

type LikesRepository interface {
	CreateLikes(like *entities.Likes) error
	GetLikeByUserID(userID string) ([]entities.Likes, error)
	CheckLike(userID string, videoID string) error
	CountLikeVideoByVideoID(videoID string) (int, error)
	DeleteLikeByID(userID string, videoID string) error
	DeleteLikeByVideoID(videoID string) error
}

func (r *GormLikesRepository) CreateLikes(like *entities.Likes) error {
	var user entities.User
	var video entities.Video
	if err := r.db.First(&user, "id = ?", like.UserID).Error; err != nil {
		return fmt.Errorf("user_id not found: %v", err)
	}

	if err := r.db.First(&video, "id = ?", like.VideoID).Error; err != nil {
		return fmt.Errorf("video_id not found: %v", err)
	}

	if err := r.db.Create(&like).Error; err != nil {
		return err
	}

	return nil
}

func (r *GormLikesRepository) GetLikeByUserID(userID string) ([]entities.Likes, error) {
	var likes []entities.Likes
	if err := r.db.Preload("Video").Where("user_id = ?", userID).First(&likes).Error; err != nil {
		return nil, err
	}

	return likes, nil
}

func (r *GormLikesRepository) CheckLike(userID string, videoID string) error {
	var like entities.Likes
	if err := r.db.Where("user_id = ? AND video_id = ?", userID, videoID).First(&like).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("not liked video")
		}

		return err
	}

	return nil
}

func (r *GormLikesRepository) DeleteLikeByID(userID string, videoID string) error {
	return r.db.Where("user_id = ? AND video_id = ?", userID, videoID).Delete(&entities.Likes{}).Error
}

func (r *GormLikesRepository) CountLikeVideoByVideoID(videoID string) (int, error) {
	var count int64
	if err := r.db.Model(&entities.Likes{}).Where("video_id = ?", videoID).Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

func (r *GormLikesRepository) DeleteLikeByVideoID(videoID string) error {
	return r.db.Where("video_id = ?", videoID).Delete(&entities.Likes{}).Error
}
