package services

import (
	"github.com/rustoma/octo-pulse/internal/models"
	"github.com/rustoma/octo-pulse/internal/storage"
)

type CategoryService interface {
	GetCategories() ([]*models.Category, error)
	GetCategory(id int) (*models.Category, error)
	GetDomainCategories(domainId int) ([]*models.Category, error)
}

type categoryService struct {
	categoryStore          storage.CategoryStore
	categoriesDomainsStore storage.CategoriesDomainsStore
}

func NewCategoryService(categoryStore storage.CategoryStore, categoriesDomainsStore storage.CategoriesDomainsStore) CategoryService {
	return &categoryService{categoryStore: categoryStore, categoriesDomainsStore: categoriesDomainsStore}
}

func (s *categoryService) GetCategories() ([]*models.Category, error) {
	return s.categoryStore.GetCategories()
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
