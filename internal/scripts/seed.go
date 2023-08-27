package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	lr "github.com/rustoma/octo-pulse/internal/logger"
)

func main() {
	//Init logger
	logger, logFile, err := lr.NewLogger()
	defer logFile.Close()

	if err != nil {
		log.Fatal(err)
	}

	//Init .env
	err = godotenv.Load()
	if err != nil {
		logger.Fatal().Msg("Error loading .env file")
	}

	//Init DB
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("SEED_DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		logger.Fatal().Err(err).Msg("")
	}

	logger.Info().Msg("Connected to the DB")
	defer dbpool.Close()
}
