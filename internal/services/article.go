package services

import (
	"fmt"
	"github.com/gosimple/slug"
	a "github.com/rustoma/octo-pulse/internal/ai"
	"github.com/rustoma/octo-pulse/internal/dto"
	"github.com/rustoma/octo-pulse/internal/models"
	"github.com/rustoma/octo-pulse/internal/storage"
	"github.com/rustoma/octo-pulse/internal/utils"
	"github.com/rustoma/octo-pulse/internal/validator"
	"regexp"
	"strings"
)

type ArticleService interface {
	GenerateDescription(question *models.Question) (string, error)
	UpdateArticle(articleId int, article *models.Article) (int, error)
	GetArticle(id int) (*models.Article, error)
	GetArticles(filters ...*storage.GetArticlesFilters) ([]*dto.Article, error)
	CreateArticle(article *models.Article) (int, error)
	DeleteArticle(id int) (int, error)
	RemoveDuplicateHeadingsFromArticle(articleId int) error
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
	article.Slug = slug.Make(article.Title)
	readingTime := utils.CalculateReadTime(article.Body)
	article.ReadingTime = &readingTime

	err := s.articleValidator.Validate(article)
	if err != nil {
		return 0, err
	}

	return s.articleStore.InsertArticle(article)
}

func (s *articleService) DeleteArticle(id int) (int, error) {
	return s.articleStore.DeleteArticle(id)
}

func (s *articleService) GenerateDescription(question *models.Question) (string, error) {

	description, err := s.ai.ChatGPT.GenerateArticleDescription(question)

	if err != nil {
		return "", err
	}

	return description, nil
}

func (s *articleService) UpdateArticle(articleId int, article *models.Article) (int, error) {
	article.Slug = slug.Make(article.Title)
	readingTime := utils.CalculateReadTime(article.Body)
	article.ReadingTime = &readingTime

	err := s.articleValidator.Validate(article)
	if err != nil {
		return 0, err
	}

	return s.articleStore.UpdateArticle(articleId, article)
}

func (s *articleService) GetArticle(id int) (*models.Article, error) {
	return s.articleStore.GetArticle(id)
}

func (s *articleService) GetArticles(filters ...*storage.GetArticlesFilters) ([]*dto.Article, error) {
	return s.articleStore.GetArticles(filters...)
}

func (s *articleService) RemoveDuplicateHeadingsFromArticle(articleId int) error {

	article, err := s.articleStore.GetArticle(articleId)
	if err != nil {
		return err
	}

	htmlString := article.Body

	// Define a regular expression pattern for detecting headings
	headingPattern := `(<h[1-6][^>]*>(.*?)<\/h[1-6]>)`

	// Find all matches of headings in the HTML string
	re := regexp.MustCompile(headingPattern)
	matches := re.FindAllString(utils.RemoveMultipleSpaces(htmlString), -1)

	headingMap := make(map[string]bool)

	for _, match := range matches {
		if headingMap[match] {

			logger.Info().Interface("Heading to remove: ", match).Send()
			pattern := match + `(.*?)(<h[1-6][^>]*>|$)`
			re := regexp.MustCompile(pattern)
			matches := re.FindAllString(utils.RemoveMultipleSpaces(htmlString), -1)
			modifiedHTML := utils.RemoveMultipleSpaces(htmlString)
			if len(matches) >= 1 {
				// Define the number of characters you want to extract from the end <h2> or <h3> it is 4 characters
				numCharacters := 4
				// Calculate the starting index for slicing
				startIndex := len(matches[0]) - numCharacters
				// Check if startIndex is valid
				if startIndex >= 0 {
					lastCharacters := matches[0][startIndex:]
					logger.Info().Interface("HTML with heading to remove: ", modifiedHTML).Send()
					logger.Info().Interface("match to replace: ", matches[0]).Send()
					logger.Info().Interface("lastCharacters: ", lastCharacters).Send()
					modifiedHTML = strings.Replace(modifiedHTML, matches[0], lastCharacters, 1)
				} else {
					fmt.Println("Input string is too short.")
				}

				logger.Info().Interface("modifiedHTML: ", modifiedHTML).Send()
				article.Body = modifiedHTML

				_, err := s.UpdateArticle(article.ID, article)
				if err != nil {
					return err
				}

				err = s.RemoveDuplicateHeadingsFromArticle(articleId)
				if err != nil {
					return err
				}
			}
		} else {
			headingMap[match] = true
		}
	}

	return nil
}
