package services

import (
	"github.com/rustoma/octo-pulse/internal/models"
	"github.com/rustoma/octo-pulse/internal/storage"
)

type ImageService interface {
	GetImages(filters ...*storage.GetImagesFilters) ([]*models.Image, error)
}

type imageService struct {
	imageStore storage.ImageStorageStore
}

func NewImageService(imageStorageStore storage.ImageStorageStore) ImageService {
	return &imageService{imageStore: imageStorageStore}
}

func (s *imageService) GetImages(filters ...*storage.GetImagesFilters) ([]*models.Image, error) {
	return s.imageStore.GetImages(filters...)
}
