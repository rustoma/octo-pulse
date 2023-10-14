package models

type Question struct {
	Id          int    `json:"id"`
	Question    string `json:"question"`
	Answear     string `json:"answear"`
	Href        string `json:"href"`
	PageContent string `json:"pageContent"`
	Fetched     int    `json:"fetched"`
	CategoryId  int    `json:"categoryId"`
}
