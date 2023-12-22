package controllers

import (
	"github.com/rs/zerolog"
	lr "github.com/rustoma/octo-pulse/internal/logger"
)

var logger *zerolog.Logger

func init() {
	//Init logger
	l, lF := lr.NewLogger()

	logger = l
	defer lF.Close()
}
