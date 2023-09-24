package services

import (
	a "github.com/rustoma/octo-pulse/internal/ai"
	"github.com/rustoma/octo-pulse/internal/models"
	"github.com/rustoma/octo-pulse/internal/storage"
)

type ArticleService interface {
	GenerateDescription() (string, error)
	UpdateArticle(articleId int, article *models.Article) (int, error)
	GetArticle(articleId int) (*models.Article, error)
	GetArticles() ([]*models.Article, error)
}

type articleService struct {
	articleStore storage.ArticleStore
	ai           *a.AI
}

func NewArticleService(articleStore storage.ArticleStore, ai *a.AI) ArticleService {
	return &articleService{articleStore: articleStore, ai: ai}
}

func (s *articleService) GenerateDescription() (string, error) {

	description, err := s.ai.ChatGPT.GenerateArticleDescription()

	if err != nil {
		return "", err
	}

	return description, nil
}

func (s *articleService) UpdateArticle(articleId int, article *models.Article) (int, error) {
	return s.articleStore.UpdateArticle(articleId, article)
}

func (s *articleService) GetArticle(articleId int) (*models.Article, error) {
	return s.articleStore.GetArticle(articleId)
}

func (s *articleService) GetArticles() ([]*models.Article, error) {
	return s.articleStore.GetArticles()
}
