package services

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	lr "github.com/rustoma/octo-pulse/internal/logger"
	"os"
	"path/filepath"
)

var logger *zerolog.Logger

func init() {
	//Init logger
	l, logFile := lr.NewLogger()
	defer logFile.Close()
	logger = l

	//Init .env
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logger.Fatal().Msg("Error loading .env file")
	}

	if err := godotenv.Load(filepath.Join(dir, ".env")); err != nil {
		logger.Fatal().Msg("Error loading .env file")
	}
}
