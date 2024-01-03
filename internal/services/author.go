package services

import (
	"github.com/rustoma/octo-pulse/internal/models"
	"github.com/rustoma/octo-pulse/internal/storage"
	"github.com/rustoma/octo-pulse/internal/validator"
)

type AuthorService interface {
	GetAuthors() ([]*models.Author, error)
	CreateAuthor(author *models.Author) (int, error)
	UpdateAuthor(id int, author *models.Author) (int, error)
}

type authorService struct {
	authorStore     storage.AuthorStore
	authorValidator validator.AuthorValidatorer
}

func NewAuthorService(authorStore storage.AuthorStore, authorValidator validator.AuthorValidatorer) AuthorService {
	return &authorService{authorStore: authorStore, authorValidator: authorValidator}
}

func (s *authorService) GetAuthors() ([]*models.Author, error) {
	return s.authorStore.GetAuthors()
}

func (s *authorService) CreateAuthor(author *models.Author) (int, error) {
	err := s.authorValidator.Validate(author)
	if err != nil {
		logger.Err(err).Send()
		return 0, err
	}

	return s.authorStore.InsertAuthor(author)
}

func (s *authorService) UpdateAuthor(id int, author *models.Author) (int, error) {
	err := s.authorValidator.Validate(author)
	if err != nil {
		return 0, err
	}

	return s.authorStore.UpdateAuthor(id, author)
}
