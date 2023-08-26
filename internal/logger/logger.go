package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rustoma/octo-pulse/internal/utils"
)

var (
	Logger zerolog.Logger
)

func InitLogger() (*os.File, error) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	file, err := os.OpenFile(
		"myapp.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)

	if err != nil {
		return nil, err
	}

	if utils.IsProdDev() {
		Logger = zerolog.New(file).With().Timestamp().Logger()
	} else {
		Logger = log.Logger
	}

	return file, nil
}
