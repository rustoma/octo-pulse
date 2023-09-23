package tasks

import (
	"os"

	"github.com/hibiken/asynq"
)

type Inspectorer interface {
	GetTasksInfo(queue string, taskIds []string) map[string]*TaskInfo
}

type Inspector struct{}

func NewTaskInspector() Inspectorer {
	return &Inspector{}
}

type TaskInfo struct {
	ID       string `json:"id"`
	Queue    string `json:"queue"`
	MaxRetry int    `json:"maxRetry"`
	Retried  int    `json:"retried"`
	State    string `json:"state"`
}

func (t *Inspector) GetTasksInfo(queue string, taskIds []string) map[string]*TaskInfo {
	inspector := asynq.NewInspector(asynq.RedisClientOpt{Addr: os.Getenv("REDIS_ADDR"), Password: os.Getenv("REDIS_PASSWORD")})
	defer inspector.Close()

	taskInfos := make(map[string]*TaskInfo)

	for _, taskId := range taskIds {
		info, err := inspector.GetTaskInfo(queue, taskId)

		if err != nil {
			logger.Err(err).Send()
		}

		taskInfos[taskId] = &TaskInfo{
			ID:       info.ID,
			Queue:    info.Queue,
			MaxRetry: info.MaxRetry,
			Retried:  info.Retried,
			State:    info.State.String(),
		}
	}

	return taskInfos
}
