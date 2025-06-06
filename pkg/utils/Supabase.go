package utils

import (
	"Beside-Mom-BE/configs"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	storage_go "github.com/supabase-community/storage-go"
)

func UploadImage(fileName string, dir string, config configs.Supabase) (string, error) {
	filePath := "./uploads/" + fileName
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}

	defer file.Close()
	if config.URL == "" || config.Key == "" || config.Bucket == "" {
		return "", fmt.Errorf("invalid Supabase config: URL='%s', Key='%s', Bucket='%s'", config.URL, config.Key, config.Bucket)
	}

	storageClient := storage_go.NewClient(config.URL, config.Key, nil)
	if storageClient == nil {
		return "", fmt.Errorf("failed to create storage client: invalid Supabase configuration")
	}

	options := storage_go.FileOptions{
		ContentType: func() *string { s := "image/jpeg"; return &s }(),
	}

	fileName = dir + fileName
	_, err = storageClient.UploadFile(config.Bucket, fileName, file, options)
	if err != nil {
		return "", fmt.Errorf("failed to upload file '%s' to bucket '%s': %w", fileName, config.Bucket, err)
	}

	url := fmt.Sprintf("%s/object/public/%s/%s", config.URL, config.Bucket, fileName)
	return url, nil
}

func UploadVideo(fileName string, file io.Reader, config configs.Supabase) (string, error) {
	if config.URL == "" || config.Key == "" || config.Bucket == "" {
		return "", fmt.Errorf("invalid Supabase config")
	}

	storageClient := storage_go.NewClient(config.URL, config.Key, nil)
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	if len(fileBytes) == 0 {
		return "", fmt.Errorf("file is empty")
	}

	fileReader := bytes.NewReader(fileBytes)
	options := storage_go.FileOptions{
		ContentType: stringPtr("video/mp4"),
	}

	_, err = storageClient.UploadFile(config.Bucket, fileName, fileReader, options)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	url := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", config.URL, config.Bucket, fileName)
	return url, nil
}

func DeleteImage(fileURL string, config configs.Supabase) error {
	if config.URL == "" || config.Key == "" || config.Bucket == "" {
		return fmt.Errorf("invalid Supabase config")
	}

	fileName := extractFileNameFromURL(fileURL, config.Bucket)
	storageClient := storage_go.NewClient(config.URL, config.Key, nil)

	_, err := storageClient.RemoveFile(config.Bucket, []string{fileName})
	if err != nil {
		return fmt.Errorf("failed to remove file: %w", err)
	}

	return nil
}

func extractFileNameFromURL(fullURL, bucketName string) string {
	prefix := "/storage/v1/object/public/" + bucketName + "/"
	index := strings.Index(fullURL, prefix)
	if index == -1 {
		return ""
	}
	return fullURL[index+len(prefix):]
}

func stringPtr(s string) *string {
	return &s
}
