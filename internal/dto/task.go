package dto

type GetTasksInfoRequest struct {
	Queue   string   `json:"queue"`
	TaskIds []string `json:"taskIds"`
}
