package tasks

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"

	"github.com/rustoma/octo-pulse/internal/ai"
	lr "github.com/rustoma/octo-pulse/internal/logger"
	"github.com/rustoma/octo-pulse/internal/models"
	"github.com/rustoma/octo-pulse/internal/services"
)

var logger *zerolog.Logger

type Tasks struct {
	Article  ArticleTasker
	Scrapper ScrapperTasker
}

func NewTasks(
	articleService services.ArticleService,
	domainService services.DomainService,
	scrapperService services.ScrapperService,
	categoryService services.CategoryService,
	ai *ai.AI) *Tasks {
	scrapperTasks := NewScrapperTasks(scrapperService)

	return &Tasks{
		Article:  NewArticleTasks(articleService, domainService, scrapperService, categoryService, ai, scrapperTasks),
		Scrapper: scrapperTasks,
	}
}

type ArticleTasker interface {
	NewGenerateDescriptionTask(pageId int, questionId int) error
	HandleGenerateDescription(ctx context.Context, task *asynq.Task) error
	NewGenerateArticlesTask(domainId int, numberOfArticlesToCreate int, questionCategoryId int) error
	HandleGenerateArticles(ctx context.Context, task *asynq.Task) error
}

type ScrapperTasker interface {
	NewUpdateQuestionTask(id int, question *models.Question) error
	HandleUpdateQuestionTask(ctx context.Context, task *asynq.Task) error
}

func init() {
	//Init logger
	l, logFile := lr.NewLogger()
	defer logFile.Close()
	logger = l
}
