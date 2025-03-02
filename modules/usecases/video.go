package usecases

import (
	"Beside-Mom-BE/configs"
	"Beside-Mom-BE/modules/entities"
	"Beside-Mom-BE/modules/repositories"
	"Beside-Mom-BE/pkg/utils"
	"io"
	"mime/multipart"

	"github.com/google/uuid"
)

type VideoUseCase interface {
	CreateVideo(video *entities.Video, videoFile *multipart.FileHeader, file io.Reader) (*entities.Video, error)
	CreateVideowithLink(video *entities.Video) (*entities.Video, error)
	GetAllVideo() ([]map[string]interface{}, error)
	GetVideoByID(id string) (map[string]interface{}, error)
	IncreaseView(id string) error
	UpdateVideo(id string, userID string, video *entities.Video, videoFile *multipart.FileHeader, file io.Reader) (*entities.Video, error)
	UpdateVideowithLink(id string, userID string, video *entities.Video) (*entities.Video, error)
	DeleteVideoByID(id string) error
}

type VideoUseCaseImpl struct {
	repo     repositories.VideoRepository
	likerepo repositories.LikesRepository
	supa     configs.Supabase
}

func NewVideoUseCase(repo repositories.VideoRepository, likerepo repositories.LikesRepository, supa configs.Supabase) *VideoUseCaseImpl {
	return &VideoUseCaseImpl{
		repo:     repo,
		likerepo: likerepo,
		supa:     supa,
	}
}

func (u *VideoUseCaseImpl) CreateVideo(video *entities.Video, videoFile *multipart.FileHeader, file io.Reader) (*entities.Video, error) {
	if videoFile != nil {
		fileName := uuid.New().String() + "_video.mp4"
		videoUrl, err := utils.UploadVideo(fileName, file, u.supa)
		if err != nil {
			return nil, err
		}

		video.Link = videoUrl
	}

	createdVideo, err := u.repo.CreateVideo(video)
	if err != nil {
		return nil, err
	}

	return createdVideo, nil
}

func (u *VideoUseCaseImpl) CreateVideowithLink(video *entities.Video) (*entities.Video, error) {
	return u.repo.CreateVideo(video)
}

func (u *VideoUseCaseImpl) GetAllVideo() ([]map[string]interface{}, error) {
	videos, err := u.repo.GetAllVideo()
	if err != nil {
		return nil, err
	}

	var videoList []map[string]interface{}
	for _, video := range videos {
		countLike, err := u.likerepo.CountLikeVideoByVideoID(video.ID)
		if err != nil {
			return nil, err
		}

		videoData := map[string]interface{}{
			"id":          video.ID,
			"title":       video.Title,
			"description": video.Description,
			"link":        video.Link,
			"view":        video.View,
			"count_like":  countLike,
		}

		videoList = append(videoList, videoData)
	}

	return videoList, nil
}

func (u *VideoUseCaseImpl) GetVideoByID(id string) (map[string]interface{}, error) {
	video, err := u.repo.GetVideoByID(id)
	if err != nil {
		return nil, err
	}

	countLike, err := u.likerepo.CountLikeVideoByVideoID(id)
	if err != nil {
		return nil, err
	}

	videoData := map[string]interface{}{
		"id":          video.ID,
		"title":       video.Title,
		"description": video.Description,
		"link":        video.Link,
		"view":        video.View,
		"count_like":  countLike,
	}

	return videoData, nil
}

func (u *VideoUseCaseImpl) IncreaseView(id string) error {
	video, err := u.repo.GetVideoByID(id)
	if err != nil {
		return err
	}

	video.View = video.View + 1
	_, err = u.repo.UpdateVideoByID(video)
	if err != nil {
		return err
	}

	return nil
}

func (u *VideoUseCaseImpl) UpdateVideo(id string, userID string, video *entities.Video, videoFile *multipart.FileHeader, file io.Reader) (*entities.Video, error) {
	existingVideo, err := u.repo.GetVideoByID(id)
	if err != nil {
		return nil, err
	}

	existingVideo.Title = video.Title
	existingVideo.Description = video.Description
	if videoFile != nil {
		fileName := uuid.New().String() + "_video.mp4"
		videoUrl, err := utils.UploadVideo(fileName, file, u.supa)
		if err != nil {
			return nil, err
		}

		existingVideo.Link = videoUrl
	}

	updatedVideo, err := u.repo.UpdateVideoByID(existingVideo)
	if err != nil {
		return nil, err
	}

	return updatedVideo, nil
}

func (u *VideoUseCaseImpl) UpdateVideowithLink(id string, userID string, video *entities.Video) (*entities.Video, error) {
	existingVideo, err := u.repo.GetVideoByID(id)
	if err != nil {
		return nil, err
	}

	existingVideo.Title = video.Title
	existingVideo.Description = video.Description
	if video.Link != "" {
		existingVideo.Link = video.Link
	}

	return u.repo.UpdateVideoByID(existingVideo)
}

func (u *VideoUseCaseImpl) DeleteVideoByID(id string) error {
	if err := u.likerepo.DeleteLikeByVideoID(id); err != nil {
		return err
	}

	if err := u.repo.DeleteVideoByID(id); err != nil {
		return err
	}

	return nil
}
