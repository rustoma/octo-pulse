package app

import "github.com/rs/zerolog"

type Ctx struct {
	Logger *zerolog.Logger
}

func NewAppCtx(logger *zerolog.Logger) *Ctx {
	return &Ctx{
		Logger: logger,
	}
}
