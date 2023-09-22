package services

import (
	"fmt"
	"net/http"

	a "github.com/rustoma/octo-pulse/internal/ai"
	"github.com/rustoma/octo-pulse/internal/api"
	"github.com/rustoma/octo-pulse/internal/models"
	"github.com/rustoma/octo-pulse/internal/storage"
)

type ArticleService interface {
	GenerateDescription(articleId int) (string, error)
	UpdateArticle(articleId int, article *models.Article) (int, error)
	GetArticle(articleId int) (*models.Article, error)
}

type articleService struct {
	articleStore storage.ArticleStore
	ai           *a.AI
}

func NewArticleService(articleStore storage.ArticleStore, ai *a.AI) ArticleService {
	return &articleService{articleStore: articleStore, ai: ai}
}

func (s *articleService) GenerateDescription(articleId int) (string, error) {

	article, err := s.articleStore.GetArticle(articleId)

	if err != nil {
		return "", api.Error{Err: err.Error(), Status: http.StatusInternalServerError}
	}

	if article == nil {
		return "", api.Error{Err: fmt.Sprintf("article with %d not found", articleId), Status: http.StatusBadRequest}
	}

	description, err := s.ai.ChatGPT.GenerateArticleDescription()

	if err != nil {
		return "", api.Error{Err: err.Error(), Status: http.StatusInternalServerError}
	}

	return description, nil
}

func (s *articleService) UpdateArticle(articleId int, article *models.Article) (int, error) {
	return s.articleStore.UpdateArticle(articleId, article)
}

func (s *articleService) GetArticle(articleId int) (*models.Article, error) {
	return s.articleStore.GetArticle(articleId)
}
