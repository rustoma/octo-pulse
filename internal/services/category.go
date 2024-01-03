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
	AssignCategoryToDomain(categoryId int, domainId int) error
}

type categoryService struct {
	categoryStore          storage.CategoryStore
	categoriesDomainsStore storage.CategoriesDomainsStore
	categoryValidator      validator.CategoryValidatorer
}

func NewCategoryService(categoryStore storage.CategoryStore, categoriesDomainsStore storage.CategoriesDomainsStore, categoryValidator validator.CategoryValidatorer) CategoryService {
	return &categoryService{categoryStore: categoryStore, categoriesDomainsStore: categoriesDomainsStore, categoryValidator: categoryValidator}
}

func (s *categoryService) GetCategories(filters ...*storage.GetCategoriesFilters) ([]*models.Category, error) {
	return s.categoryStore.GetCategories(filters...)
}

func (s *categoryService) GetCategory(id int) (*models.Category, error) {
	return s.categoryStore.GetCategory(id)
}

func (s *categoryService) GetDomainCategories(domainId int) ([]*models.Category, error) {
	categoryIds, err := s.categoriesDomainsStore.GetDomainCategories(domainId)
	if err != nil {
		return nil, err
	}

	var categories []*models.Category
	for _, categoryId := range categoryIds {
		category, err := s.categoryStore.GetCategory(categoryId)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
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

func (s *categoryService) AssignCategoryToDomain(categoryId int, domainId int) error {
	return s.categoriesDomainsStore.AssignCategoryToDomain(categoryId, domainId)
}
