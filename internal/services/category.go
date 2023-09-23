package services

import (
	"github.com/rustoma/octo-pulse/internal/models"
	"github.com/rustoma/octo-pulse/internal/storage"
)

type CategoryService interface {
	GetCategories() ([]*models.Category, error)
}

type categoryService struct {
	categoryStore storage.CategoryStore
}

func NewCategoryService(categoryStore storage.CategoryStore) CategoryService {
	return &categoryService{categoryStore: categoryStore}
}

func (s *categoryService) GetCategories() ([]*models.Category, error) {
	return s.categoryStore.GetCategories()
}
