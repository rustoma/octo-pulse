package db

import (
	"context"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	lr "github.com/rustoma/octo-pulse/internal/logger"
)

var logger *zerolog.Logger

func Connect() (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	logger.Info().Msg("Connected to the DB")

	return dbpool, err
}

func init() {

	//Init logger
	l, logFile := lr.NewLogger()
	defer logFile.Close()
	logger = l
	logFile = logFile

	//Init .env
	if err := godotenv.Load(filepath.Join(".", ".env")); err != nil {
		logger.Fatal().Msg("Error loading .env file")
	}
}
