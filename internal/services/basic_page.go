package services

import (
	"github.com/gosimple/slug"
	"github.com/rustoma/octo-pulse/internal/models"
	"github.com/rustoma/octo-pulse/internal/storage"
	"github.com/rustoma/octo-pulse/internal/validator"
)

type BasicPageService interface {
	CreateBasicPage(page *models.BasicPage) (int, error)
	GetBasicPage(id int) (*models.BasicPage, error)
	GetBasicPageBySlug(slug string, filters ...*storage.GetBasicPageBySlugFilters) (*models.BasicPage, error)
	GetBasicPages(filters ...*storage.GetBasicPagesFilters) ([]*models.BasicPage, error)
	UpdateBasicPage(id int, basicPage *models.BasicPage) (int, error)
}

type basicPageService struct {
	basicPageStore     storage.BasicPageStore
	basicPageValidator validator.BasicPageValidatorer
}

func NewBasicPageService(basicPageStore storage.BasicPageStore, basicPageValidator validator.BasicPageValidatorer) BasicPageService {
	return &basicPageService{basicPageStore: basicPageStore, basicPageValidator: basicPageValidator}
}

func (s *basicPageService) GetBasicPages(filters ...*storage.GetBasicPagesFilters) ([]*models.BasicPage, error) {
	return s.basicPageStore.GetBasicPages(filters...)
}

func (s *basicPageService) GetBasicPage(id int) (*models.BasicPage, error) {
	return s.basicPageStore.GetBasicPage(id)
}

func (s *basicPageService) GetBasicPageBySlug(slug string, filters ...*storage.GetBasicPageBySlugFilters) (*models.BasicPage, error) {
	return s.basicPageStore.GetBasicPageBySlug(slug, filters...)
}

func (s *basicPageService) CreateBasicPage(page *models.BasicPage) (int, error) {
	page.Slug = slug.Make(page.Title)

	err := s.basicPageValidator.Validate(page)
	if err != nil {
		logger.Err(err).Send()
		return 0, err
	}

	return s.basicPageStore.InsertBasicPage(page)
}

func (s *basicPageService) UpdateBasicPage(id int, basicPage *models.BasicPage) (int, error) {
	basicPage.Slug = slug.Make(basicPage.Title)

	err := s.basicPageValidator.Validate(basicPage)
	if err != nil {
		return 0, err
	}

	return s.basicPageStore.UpdateBasicPage(id, basicPage)
}
