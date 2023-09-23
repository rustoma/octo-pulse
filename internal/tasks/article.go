package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/hibiken/asynq"
	"github.com/rustoma/octo-pulse/internal/services"
)

const (
	TypeArticleGenerateDescription = "article:generateDescription"
)

type articleTasks struct {
	articleService services.ArticleService
	inspector      *asynq.Inspector
}

func NewArticleTasks(articleService services.ArticleService) articleTasks {
	return articleTasks{
		articleService: articleService,
	}
}

type DescriptionTaskPayload struct {
	ArticleId int
}

func (t articleTasks) NewGenerateDescriptionTask(articleId int) error {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: os.Getenv("REDIS_ADDR"), Password: os.Getenv("REDIS_PASSWORD")})
	defer client.Close()

	payload, err := json.Marshal(DescriptionTaskPayload{ArticleId: articleId})
	if err != nil {
		return err
	}

	task := asynq.NewTask(TypeArticleGenerateDescription, payload)
	info, err := client.Enqueue(task, asynq.MaxRetry(5))

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

	if err != nil {
		return err
	}

	if article == nil {
		return fmt.Errorf("article with %d not found", payload.ArticleId)
	}

	description, err := t.articleService.GenerateDescription(payload.ArticleId)

	if err != nil {
		return err
	}

	article.Description = description

	_, err = t.articleService.UpdateArticle(payload.ArticleId, article)

	if err != nil {
		return err
	}

	return nil
}
