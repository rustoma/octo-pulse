package services

import (
	"github.com/rustoma/octo-pulse/internal/models"
	"github.com/rustoma/octo-pulse/internal/storage"
)

type BasicPageService interface {
	InsertBasicPage(page *models.BasicPage) (int, error)
	GetBasicPage(id int) (*models.BasicPage, error)
	GetBasicPageBySlug(slug string, filters ...*storage.GetBasicPageBySlugFilters) (*models.BasicPage, error)
	GetBasicPages(filters ...*storage.GetBasicPagesFilters) ([]*models.BasicPage, error)
}

type basicPageService struct {
	basicPageStore storage.BasicPageStore
}

func NewBasicPageService(basicPageStore storage.BasicPageStore) BasicPageService {
	return &basicPageService{basicPageStore: basicPageStore}
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

func (s *basicPageService) InsertBasicPage(page *models.BasicPage) (int, error) {
	return s.basicPageStore.InsertBasicPage(page)
}
