package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/hibiken/asynq"
	"github.com/rustoma/octo-pulse/internal/models"
	"github.com/rustoma/octo-pulse/internal/services"
)

const (
	TypeScrapperUpdateQuestion = "scrapper:updateQuestion"
)

type scrapperTasks struct {
	scrapperService services.ScrapperService
}

func NewScrapperTasks(scrapperService services.ScrapperService) scrapperTasks {
	return scrapperTasks{
		scrapperService: scrapperService,
	}
}

type UpdateQuestionTaskPayload struct {
	Id       int
	Question *models.Question
}

func (t scrapperTasks) NewUpdateQuestionTask(id int, question *models.Question) error {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: os.Getenv("REDIS_ADDR"), Password: os.Getenv("REDIS_PASSWORD")})
	defer client.Close()

	payload, err := json.Marshal(UpdateQuestionTaskPayload{
		Id:       id,
		Question: question,
	})

	if err != nil {
		return err
	}

	task := asynq.NewTask(TypeScrapperUpdateQuestion, payload)
	info, err := client.Enqueue(task, asynq.MaxRetry(20))
	logger.Info().Msgf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	if err != nil {
		return err
	}

	return nil
}

func (t scrapperTasks) HandleUpdateQuestionTask(ctx context.Context, task *asynq.Task) error {
	var payload UpdateQuestionTaskPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	payload.Question.Fetched = 1
	err := t.scrapperService.UpdateQuestion(payload.Id, payload.Question)
	if err != nil {
		return err
	}
	logger.Info().Interface("Updated question ID", payload.Question.Id).Send()

	return nil
}
