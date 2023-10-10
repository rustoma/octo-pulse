package services

import (
	"github.com/rustoma/octo-pulse/internal/models"
	"github.com/rustoma/octo-pulse/internal/storage"
)

type ScrapperService interface {
	GetQuestion(id int) (*models.Question, error)
}

type scrapperService struct {
	scrapperStore storage.ScrapperStore
}

type Question struct {
	Id          int    `json:"id"`
	Question    string `json:"question"`
	Answer      string `json:"answer"`
	Href        string `json:"href"`
	PageContent string `json:"page_content"`
}

type Category int

const (
	Budowlanka Category = iota + 1
	Gastronomia
)

type GetQuestionResponse struct {
	Question Question `json:"question"`
}

func NewScrapperService(scrapperStore storage.ScrapperStore) ScrapperService {
	return &scrapperService{
		scrapperStore: scrapperStore,
	}
}

func (s *scrapperService) GetQuestion(id int) (*models.Question, error) {
	question, err := s.scrapperStore.GetQuestion(id)
	// req, err := http.NewRequest("GET", fmt.Sprintf("%s/question/%d", s.baseUrl, id), nil)
	// req.Header.Set("Authorization", os.Getenv("BOT_API_KEY"))
	// req.Header.Set("Content-Type", "application/json")

	// client := &http.Client{}
	// resp, err := client.Do(req)
	// if err != nil {
	// 	return nil, err
	// }
	// defer resp.Body.Close()

	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return nil, err
	// }

	// if resp.Status != "200" {
	// 	var responseError struct{ Message string }
	// 	err = json.Unmarshal(body, &responseError)
	// 	return nil, errors.New(responseError.Message)
	// }

	// var question GetQuestionResponse
	// err = json.Unmarshal(body, &question)
	// if err != nil {
	// 	return nil, err
	// }

	return question, err
}
