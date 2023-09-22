package chatgpt

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	lr "github.com/rustoma/octo-pulse/internal/logger"
	"github.com/sashabaranov/go-openai"
)

var logger *zerolog.Logger

type ChatGPTer interface {
	GenerateArticleDescription() (string, error)
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