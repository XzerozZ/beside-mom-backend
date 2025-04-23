package usecases

import (
	"Beside-Mom-BE/modules/entities"
	"Beside-Mom-BE/modules/repositories"
)

type LikeUseCase interface {
	CreateLikes(like *entities.Likes) error
	GetLikeByUserID(userID string) ([]entities.Likes, error)
	CheckLike(userID string, videoID string) error
	DeleteLikeByID(userID string, videoID string) error
}

type LikeUseCaseImpl struct {
	repo repositories.LikesRepository
}

func NewLikeUseCase(repo repositories.LikesRepository) *LikeUseCaseImpl {
	return &LikeUseCaseImpl{repo: repo}
}

func (u *LikeUseCaseImpl) CreateLikes(like *entities.Likes) error {
	return u.repo.CreateLikes(like)
}

func (u *LikeUseCaseImpl) GetLikeByUserID(userID string) ([]entities.Likes, error) {
	return u.repo.GetLikeByUserID(userID)
}

func (u *LikeUseCaseImpl) CheckLike(userID string, videoID string) error {
	return u.repo.CheckLike(userID, videoID)
}

func (u *LikeUseCaseImpl) DeleteLikeByID(userID string, videoID string) error {
	return u.repo.DeleteLikeByID(userID, videoID)
}
