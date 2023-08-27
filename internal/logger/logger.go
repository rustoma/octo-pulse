package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rustoma/octo-pulse/internal/utils"
)

func NewLogger() (*zerolog.Logger, *os.File, error) {
	var (
		logger zerolog.Logger
	)

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	file, err := os.OpenFile(
		"myapp.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)

	if err != nil {
		return nil, nil, err
	}

	if utils.IsProdDev() {
		logger = zerolog.New(file).With().Timestamp().Logger()
	} else {
		logger = log.Logger
	}

	return &logger, file, nil
}
