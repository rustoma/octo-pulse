package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	lr "github.com/rustoma/octo-pulse/internal/logger"
)

var (
	logger zerolog.Logger
)

func main() {
	//Init logger
	logger, logFile, err := lr.NewLogger()
	defer logFile.Close()

	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	//Init .env
	err = godotenv.Load(filepath.Join(".", ".env"))
	if err != nil {
		logger.Fatal().Msg("Error loading .env file")
	}

	//Init DB
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		logger.Fatal().Err(err).Msg("")
	}
	logger.Info().Msg("Connected to the DB")
	defer dbpool.Close()

}
