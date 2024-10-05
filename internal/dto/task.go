package dto

type GetTasksInfoRequest struct {
	Queue   string   `json:"queue"`
	TaskIds []string `json:"taskIds"`
}

type GenerateArticlesRequest struct {
	DomainId           int    `json:"domainId"`
	NumberOfArticles   int    `json:"numberOfArticles"`
	QuestionCategoryId int    `json:"questionCategoryId"`
	ImagesCategory     int    `json:"imagesCategory"`
	Lang               string `json:"lang"`
}

type GenerateDescriptionRequest struct {
	QuestionId int    `json:"questionId"`
	Lang       string `json:"lang"`
}
