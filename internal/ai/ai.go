package ai

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	chatgpt "github.com/rustoma/octo-pulse/internal/ai/chatGPT"
	lr "github.com/rustoma/octo-pulse/internal/logger"
)

var logger *zerolog.Logger

type AI struct {
	ChatGPT chatgpt.ChatGPTer
}

func NewAI() *AI {
	return &AI{
		ChatGPT: chatgpt.NewChatGPT(),
	}
}

func init() {
	//Init logger
	l, logFile := lr.NewLogger()
	defer logFile.Close()
	logger = l

	//ENV init
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logger.Fatal().Msg("Error loading .env file")
	}

	if err := godotenv.Load(filepath.Join(dir, ".env")); err != nil {
		logger.Fatal().Msg("Error loading .env file")
	}
}
