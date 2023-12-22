package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gosimple/slug"
	"github.com/rustoma/octo-pulse/internal/utils"
	"os"
	"time"

	"github.com/hibiken/asynq"
	"github.com/rustoma/octo-pulse/internal/ai"
	"github.com/rustoma/octo-pulse/internal/models"
	"github.com/rustoma/octo-pulse/internal/services"
	"github.com/rustoma/octo-pulse/internal/storage"
)

const (
	TypeArticleGenerateDescription = "article:generateDescription"
	TypeArticleGenerateArticles    = "article:generateArticles"
)

type articleTasks struct {
	articleService  services.ArticleService
	domainService   services.DomainService
	scrapperService services.ScrapperService
	categoryService services.CategoryService
	ai              *ai.AI
	inspector       *asynq.Inspector
	scrapperTasks   scrapperTasks
}

func NewArticleTasks(
	articleService services.ArticleService,
	domainService services.DomainService,
	scrapperService services.ScrapperService,
	categoryService services.CategoryService,
	ai *ai.AI,
	scrapperTasks scrapperTasks,
) articleTasks {
	return articleTasks{
		articleService:  articleService,
		domainService:   domainService,
		scrapperService: scrapperService,
		categoryService: categoryService,
		ai:              ai,
		scrapperTasks:   scrapperTasks,
	}
}

type DescriptionTaskPayload struct {
	ArticleId  int
	QuestionId int
}

type GenerateArticlesTaskPayload struct {
	DomainId                 int
	NumberOfArticlesToCreate int
	QuestionCategoryId       int
}

func (t articleTasks) NewGenerateArticlesTask(domainId int, numberOfArticlesToCreate int, questionCategoryId int) error {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: os.Getenv("REDIS_ADDR"), Password: os.Getenv("REDIS_PASSWORD")})
	defer client.Close()

	payload, err := json.Marshal(GenerateArticlesTaskPayload{
		DomainId:                 domainId,
		NumberOfArticlesToCreate: numberOfArticlesToCreate,
		QuestionCategoryId:       questionCategoryId})

	if err != nil {
		return err
	}

	task := asynq.NewTask(TypeArticleGenerateArticles, payload)
	info, err := client.Enqueue(task, asynq.MaxRetry(2), asynq.Timeout(2*time.Hour))

	logger.Info().Msgf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	if err != nil {
		return err
	}

	return nil
}

func (t articleTasks) NewGenerateDescriptionTask(articleId int, questionId int) error {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: os.Getenv("REDIS_ADDR"), Password: os.Getenv("REDIS_PASSWORD")})
	defer client.Close()

	payload, err := json.Marshal(DescriptionTaskPayload{ArticleId: articleId, QuestionId: questionId})
	if err != nil {
		return err
	}

	task := asynq.NewTask(TypeArticleGenerateDescription, payload)
	info, err := client.Enqueue(task, asynq.MaxRetry(2), asynq.Timeout(2*time.Hour))

	logger.Info().Msgf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	if err != nil {
		return err
	}

	return nil
}

func (t articleTasks) HandleGenerateDescription(ctx context.Context, task *asynq.Task) error {
	var payload DescriptionTaskPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	article, err := t.articleService.GetArticle(payload.ArticleId)
	question, err := t.scrapperService.GetQuestion(payload.QuestionId)

	if err != nil {
		return err
	}

	if article == nil {
		return fmt.Errorf("article with %d not found", payload.ArticleId)
	}

	if question == nil {
		return fmt.Errorf("question with %d not found", payload.QuestionId)
	}

	description, err := t.articleService.GenerateDescription(question)

	if err != nil {
		return err
	}

	article.Body = description

	readingTime := utils.CalculateReadTime(description)
	article.ReadingTime = &readingTime

	_, err = t.articleService.UpdateArticle(payload.ArticleId, article)

	if err != nil {
		return err
	}

	return nil
}

func (t articleTasks) HandleGenerateArticles(ctx context.Context, task *asynq.Task) error {
	var payload GenerateArticlesTaskPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	logger.Info().Interface("payload", payload).Send()

	questions, err := t.scrapperService.GetQuestions(&storage.GetQuestionsFilters{CategoryId: payload.QuestionCategoryId})

	if err != nil {
		return err
	}

	domainCategories, err := t.categoryService.GetDomainCategories(payload.DomainId)

	if err != nil {
		return err
	}

	logger.Info().Interface("Domain categories", domainCategories).Send()

	createdArticles := 0
	for _, question := range questions {

		if createdArticles == payload.NumberOfArticlesToCreate {
			break
		}

		//TODO: add validation for question source total length

		catgoryId, err := t.ai.ChatGPT.AssignToCategory(domainCategories, question)
		if err != nil {
			return err
		}

		if catgoryId == 0 {
			logger.Info().Msg("There is no category that fits")
			continue
		}

		logger.Info().Interface("Assigned to category", catgoryId).Send()

		article := &models.Article{
			Title:      question.Question,
			Slug:       slug.Make(question.Question),
			Body:       "",
			Thumbnail:  nil,
			CategoryId: catgoryId,
			AuthorId:   1,
			DomainId:   payload.DomainId,
			Featured:   false,
			CreatedAt:  time.Now().UTC(),
			UpdatedAt:  time.Now().UTC(),
		}

		articleId, err := t.articleService.CreateArticle(article)
		if err != nil {
			return err
		}

		logger.Info().Interface("CreatedArticle ID", articleId).Send()

		//Increase number of created articles
		createdArticles++

		err = t.scrapperTasks.NewUpdateQuestionTask(question.Id, question)
		if err != nil {
			_, _ = t.articleService.DeleteArticle(articleId)
			return err
		}

		//Generate Description For article
		_ = t.NewGenerateDescriptionTask(articleId, question.Id)
	}

	return nil
}
