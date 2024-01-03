package services

import (
	"github.com/rustoma/octo-pulse/internal/models"
	"github.com/rustoma/octo-pulse/internal/storage"
	"github.com/rustoma/octo-pulse/internal/validator"
)

type ScrapperService interface {
	GetQuestion(id int) (*models.Question, error)
	GetQuestions(filters ...*storage.GetQuestionsFilters) ([]*models.Question, error)
	UpdateQuestion(id int, question *models.Question) error
	GetQuestionCategories() ([]*models.QuestionCategory, error)
}

type scrapperService struct {
	scrapperStore     storage.ScrapperStore
	scrapperValidator validator.ScrapperValidatorer
}

type Category int

const (
	Budowlanka Category = iota + 1
	Gastronomia
)

func NewScrapperService(scrapperStore storage.ScrapperStore, scrapperValidator validator.ScrapperValidatorer) ScrapperService {
	return &scrapperService{
		scrapperStore:     scrapperStore,
		scrapperValidator: scrapperValidator,
	}
}

func (s *scrapperService) GetQuestion(id int) (*models.Question, error) {
	question, err := s.scrapperStore.GetQuestion(id)
	return question, err
}

func (s *scrapperService) GetQuestions(filters ...*storage.GetQuestionsFilters) ([]*models.Question, error) {
	questions, err := s.scrapperStore.GetQuestions(filters...)
	return questions, err
}

func (s *scrapperService) UpdateQuestion(id int, question *models.Question) error {
	err := s.scrapperValidator.Validate(question)

	if err != nil {
		return err
	}

	return s.scrapperStore.UpdateQuestion(id, question)
}

func (s *scrapperService) GetQuestionCategories() ([]*models.QuestionCategory, error) {
	questions, err := s.scrapperStore.GetQuestionCategories()
	return questions, err
}
