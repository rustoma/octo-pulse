package tasks

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"

	lr "github.com/rustoma/octo-pulse/internal/logger"
	"github.com/rustoma/octo-pulse/internal/services"
)

var logger *zerolog.Logger

type Tasks struct {
	Article ArticleTasker
}

func NewTasks(articleService services.ArticleService) *Tasks {
	return &Tasks{Article: NewArticleTasks(articleService)}
}

type ArticleTasker interface {
	NewGenerateDescriptionTask(pageId int) error
	HandleGenerateDescription(ctx context.Context, task *asynq.Task) error
}

func init() {
	//Init logger
	l, logFile := lr.NewLogger()
	defer logFile.Close()
	logger = l
}
