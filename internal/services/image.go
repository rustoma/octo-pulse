package services

import (
	"errors"
	"github.com/gosimple/slug"
	e "github.com/rustoma/octo-pulse/internal/errors"
	"github.com/rustoma/octo-pulse/internal/models"
	"github.com/rustoma/octo-pulse/internal/storage"
	"github.com/rustoma/octo-pulse/internal/validator"
	"image"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"time"
)

type ImageService interface {
	GetImages(filters ...*storage.GetImagesFilters) ([]*models.Image, error)
	GetImage(id int) (*models.Image, error)
	GetImageCategories() ([]*models.ImageCategory, error)
	UploadImage(image multipart.File, handler *multipart.FileHeader, imageCategory int) (int, error)
	CreateImageCategory(category *models.ImageCategory) (int, error)
	UpdateImageCategory(id int, category *models.ImageCategory) (int, error)
	GetImageCategory(id int) (*models.ImageCategory, error)
}

type imageService struct {
	imageStore             storage.ImageStorageStore
	imageCategoryStore     storage.ImageCategoryStore
	imageCategoryValidator validator.ImageCategoryValidatorer
}

func NewImageService(imageStorageStore storage.ImageStorageStore, imageCategoryStore storage.ImageCategoryStore, imageCategoryValidator validator.ImageCategoryValidatorer) ImageService {
	return &imageService{imageStore: imageStorageStore, imageCategoryStore: imageCategoryStore, imageCategoryValidator: imageCategoryValidator}
}

func (s *imageService) GetImages(filters ...*storage.GetImagesFilters) ([]*models.Image, error) {
	return s.imageStore.GetImages(filters...)
}

func (s *imageService) GetImage(id int) (*models.Image, error) {
	return s.imageStore.GetImage(id)
}

func (s *imageService) GetImageCategory(id int) (*models.ImageCategory, error) {
	return s.imageCategoryStore.GetCategory(id)
}

func (s *imageService) GetImageCategories() ([]*models.ImageCategory, error) {
	return s.imageCategoryStore.GetCategories()
}

func (s *imageService) CreateImageCategory(category *models.ImageCategory) (int, error) {
	err := s.imageCategoryValidator.Validate(category)
	if err != nil {
		logger.Err(err).Send()
		return 0, err
	}

	return s.imageCategoryStore.InsertCategory(category)
}

func (s *imageService) UploadImage(file multipart.File, handler *multipart.FileHeader, imageCategory int) (int, error) {

	tempFile, err := os.CreateTemp(filepath.Join(os.Getenv("PATH_TO_ASSETS"), "images", "uploaded"), "*-"+handler.Filename)
	if err != nil {
		logger.Err(err).Send()
		return 0, errors.New("cannot create temporary file when uploading image")
	}
	defer tempFile.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		logger.Err(err).Send()
		return 0, errors.New("cannot read file when uploading image")
	}

	_, err = tempFile.Write(fileBytes)
	if err != nil {
		logger.Err(err).Send()
		return 0, errors.New("cannot save uploading image into temporary file")
	}

	filenameBase := path.Base(handler.Filename)
	extension := path.Ext(handler.Filename)
	filenameWithoutExt := filenameBase[:len(filenameBase)-len(extension)]

	imageFile, err := os.Open(tempFile.Name())
	if err != nil {
		logger.Err(err).Send()
		return 0, errors.New("cannot open temporary file")
	}
	defer imageFile.Close()

	imgCfg, _, err := image.DecodeConfig(imageFile)
	if err != nil {
		logger.Err(err).Send()
		return 0, e.BadRequest{Err: "Bad image file type"}
	}

	img := models.Image{
		Name:       slug.Make(filenameWithoutExt),
		Path:       filepath.Join("/", tempFile.Name()),
		Size:       int(handler.Size),
		Type:       handler.Header.Get("Content-Type"),
		Width:      imgCfg.Width,
		Height:     imgCfg.Height,
		Alt:        slug.Make(filenameWithoutExt),
		CategoryId: imageCategory,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}

	return s.imageStore.InsertImage(&img)
}

func (s *imageService) UpdateImageCategory(id int, category *models.ImageCategory) (int, error) {
	err := s.imageCategoryValidator.Validate(category)
	if err != nil {
		return 0, err
	}

	return s.imageCategoryStore.UpdateImageCategory(id, category)
}
