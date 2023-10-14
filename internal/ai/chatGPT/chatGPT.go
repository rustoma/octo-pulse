package chatgpt

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"regexp"
	"strconv"

	"github.com/rs/zerolog"
	lr "github.com/rustoma/octo-pulse/internal/logger"
	"github.com/rustoma/octo-pulse/internal/models"
	"github.com/sashabaranov/go-openai"
)

var logger *zerolog.Logger

type ChatGPTer interface {
	GenerateArticleDescription() (string, error)
	AssignToCategory(categories []*models.Category, question *models.Question) (int, error)
}

type chatGPT struct {
	retriesLimit int
	Client       *openai.Client
}

func NewChatGPT() ChatGPTer {
	client := openai.NewClient(os.Getenv("AI_KEY"))

	return &chatGPT{
		retriesLimit: 2,
		Client:       client,
	}
}

func (c *chatGPT) AssignToCategory(categories []*models.Category, question *models.Question) (int, error) {

	categoriesJSON, err := json.Marshal(categories)
	if err != nil {
		return 0, err
	}

	messages := []openai.ChatCompletionMessage{
		{
			Role: openai.ChatMessageRoleUser,
			Content: "Tytuł artykułu to: " + question.Question + "\n" +
				"Opis artykułu: " + question.Answear + "\n\n" +

				"Dostępne kategorie: " + string(categoriesJSON) + "\n\n" +

				"Przypasuj tytuł artykułu do jednej z podanych kategorii. Zwróć jedynie id kategorii. \n\n" +

				"Odpowiedź według zaleceń: \n\n" +

				"- zwróć jedynie id kategorii do której pasuje tytuł \n" +
				"- jeżeli tytuł nie pasuje do żadnej kategorii zwróć 0 \n" +
				"- id kategorii zwróć pomiędzy trzema myślnikami \n\n" +

				"Przykład poprawnej odpowiedzi: ---133---",
		},
	}

	logger.Info().Interface("message: ", messages).Send()

	resp, err := c.ask(messages)

	logger.Info().Interface("respond: ", resp).Send()
	if err != nil {
		return 0, err
	}

	re := regexp.MustCompile(`---(\d+)---`)
	match := re.FindStringSubmatch(resp)
	if match == nil {
		return 0, errors.New("No integer found in the respond from AssignToCategory")
	}

	categoryId, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, err
	}
	return categoryId, nil
}

func (c *chatGPT) GenerateArticleDescription() (string, error) {

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: "Wygeneruj losowy tekst",
		},
	}

	resp, err := c.ask(messages)

	if err != nil {
		return "", err
	}

	return resp, nil
}

func (ai *chatGPT) newChatCompletion(messages []openai.ChatCompletionMessage) (openai.ChatCompletionResponse, error) {

	resp, err := ai.Client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo16K,
			Messages:    messages,
			Temperature: 0,
		},
	)

	return resp, err
}

func (ai *chatGPT) retry(messages []openai.ChatCompletionMessage, retriesLimit int) (openai.ChatCompletionResponse, error) {
	var (
		retries     = 0
		isAskFailed = true
	)

	for {
		logger.Info().Msg("Re-attempting...")
		retries += 1
		logger.Info().Msgf("Retry number: %d", retries)

		resp, err := ai.newChatCompletion(messages)

		if err != nil {
			logger.Err(err).Send()
		} else {
			isAskFailed = false
		}

		if retries > retriesLimit || !isAskFailed {
			logger.Info().Msg("Repeated query successful")
			return resp, err
		}
	}

}

func (ai *chatGPT) ask(messages []openai.ChatCompletionMessage) (string, error) {
	resp, err := ai.newChatCompletion(messages)

	if err != nil {
		logger.Err(err).Send()
		resp, err = ai.retry(messages, ai.retriesLimit)
	}

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func init() {
	l, logFile := lr.NewLogger()
	defer logFile.Close()
	logger = l
}
