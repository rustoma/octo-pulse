package services

import (
	"github.com/gosimple/slug"
	"github.com/rustoma/octo-pulse/internal/models"
	"github.com/rustoma/octo-pulse/internal/storage"
	"github.com/rustoma/octo-pulse/internal/validator"
)

type CategoryService interface {
	GetCategories(filters ...*storage.GetCategoriesFilters) ([]*models.Category, error)
	GetCategory(id int) (*models.Category, error)
	GetDomainCategories(domainId int) ([]*models.Category, error)
	CreateCategory(category *models.Category) (int, error)
	UpdateCategory(id int, category *models.Category) (int, error)
}

type categoryService struct {
	categoryStore     storage.CategoryStore
	categoryValidator validator.CategoryValidatorer
}

func NewCategoryService(categoryStore storage.CategoryStore, categoryValidator validator.CategoryValidatorer) CategoryService {
	return &categoryService{categoryStore: categoryStore, categoryValidator: categoryValidator}
}

func (s *categoryService) GetCategories(filters ...*storage.GetCategoriesFilters) ([]*models.Category, error) {
	return s.categoryStore.GetCategories(filters...)
}

func (s *categoryService) GetCategory(id int) (*models.Category, error) {
	return s.categoryStore.GetCategory(id)
}

func (s *categoryService) GetDomainCategories(domainId int) ([]*models.Category, error) {
	var filters storage.GetCategoriesFilters
	filters.DomainId = domainId

	return s.categoryStore.GetCategories(&filters)
}

func (s *categoryService) CreateCategory(category *models.Category) (int, error) {
	category.Slug = slug.Make(category.Name)

	err := s.categoryValidator.Validate(category)
	if err != nil {
		logger.Err(err).Send()
		return 0, err
	}

	return s.categoryStore.InsertCategory(category)
}

func (s *categoryService) UpdateCategory(id int, category *models.Category) (int, error) {
	category.Slug = slug.Make(category.Name)

	err := s.categoryValidator.Validate(category)
	if err != nil {
		return 0, err
	}

	return s.categoryStore.UpdateCategory(id, category)
}
