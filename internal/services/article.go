package services

import (
	a "github.com/rustoma/octo-pulse/internal/ai"
	"github.com/rustoma/octo-pulse/internal/models"
	"github.com/rustoma/octo-pulse/internal/storage"
	"github.com/rustoma/octo-pulse/internal/validator"
)

type ArticleService interface {
	GenerateDescription() (string, error)
	UpdateArticle(articleId int, article *models.Article) (int, error)
	GetArticle(id int) (*models.Article, error)
	GetArticles() ([]*models.Article, error)
	CreateArticle(article *models.Article) (int, error)
	DeleteArticle(id int) (int, error)
}

type articleService struct {
	articleStore     storage.ArticleStore
	articleValidator validator.ArticleValidatorer
	ai               *a.AI
}

func NewArticleService(articleStore storage.ArticleStore, articleValidator validator.ArticleValidatorer, ai *a.AI) ArticleService {
	return &articleService{articleStore: articleStore, articleValidator: articleValidator, ai: ai}
}

func (s *articleService) CreateArticle(article *models.Article) (int, error) {
	err := s.articleValidator.Validate(article)
	if err != nil {
		return 0, err
	}

	return s.articleStore.InsertArticle(article)
}

func (s *articleService) DeleteArticle(id int) (int, error) {
	return s.articleStore.DeleteArticle(id)
}

func (s *articleService) GenerateDescription() (string, error) {

	description, err := s.ai.ChatGPT.GenerateArticleDescription()

	if err != nil {
		return "", err
	}

	return description, nil
}

func (s *articleService) UpdateArticle(articleId int, article *models.Article) (int, error) {
	err := s.articleValidator.Validate(article)

	if err != nil {
		return 0, err
	}

	return s.articleStore.UpdateArticle(articleId, article)
}

func (s *articleService) GetArticle(id int) (*models.Article, error) {
	return s.articleStore.GetArticle(id)
}

func (s *articleService) GetArticles() ([]*models.Article, error) {
	return s.articleStore.GetArticles()
}
