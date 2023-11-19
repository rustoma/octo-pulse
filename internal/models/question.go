package models

type Question struct {
	Id          int               `json:"id"`
	Question    string            `json:"question"`
	Answer      string            `json:"answer"`
	Href        string            `json:"href"`
	PageContent string            `json:"pageContent"`
	Fetched     int               `json:"fetched"`
	CategoryId  int               `json:"categoryId"`
	Sources     []*QuestionSource `json:"sources"`
}

type QuestionSource struct {
	Id          int    `json:"id"`
	QuestionId  int    `json:"questionId"`
	Href        string `json:"href"`
	PageContent string `json:"pageContent"`
}
