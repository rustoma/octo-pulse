package dto

type GetTasksInfoRequest struct {
	Queue   string   `json:"queue"`
	TaskIds []string `json:"taskIds"`
}

type GenerateArticlesRequest struct {
	DomainId           int `json:"domainId"`
	NumberOfArtilces   int `json:"numberOfArtilces"`
	QuestionCategoryId int `json:"questionCategoryId"`
}
