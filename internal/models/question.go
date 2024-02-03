package models

type Question struct {
	Id           int                    `json:"id"`
	Question     string                 `json:"question"`
	Answer       string                 `json:"answer"`
	Href         string                 `json:"href"`
	Fetched      int                    `json:"fetched"`
	CategoryId   int                    `json:"categoryId"`
	PageContents []*QuestionPageContent `json:"pageContents"`
}

type QuestionSource struct {
	Id         int    `json:"id"`
	QuestionId int    `json:"questionId"`
	Href       string `json:"href"`
}

type QuestionPageContent struct {
	SourceId             int    `json:"sourceId"`
	QuestionId           int    `json:"questionId"`
	Href                 string `json:"href"`
	PageContent          string `json:"pageContent"`
	PageContentProcessed string `json:"pageContentProcessed"`
}

type QuestionCategory struct {
	IdCategory  int    `json:"idCategory"`
	Name        string `json:"name"`
	Language    string `json:"language"`
	DateCreated string `json:"dateCreated"`
}
