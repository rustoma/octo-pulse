package services

import (
	"github.com/rustoma/octo-pulse/internal/models"
	"github.com/rustoma/octo-pulse/internal/storage"
)

type AuthorService interface {
	GetAuthors() ([]*models.Author, error)
}

type authorService struct {
	authorStore storage.AuthorStore
}

func NewAuthorService(authorStore storage.AuthorStore) AuthorService {
	return &authorService{authorStore: authorStore}
}

func (s *authorService) GetAuthors() ([]*models.Author, error) {
	return s.authorStore.GetAuthors()
}
