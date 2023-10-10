package db

import (
	"context"
	"database/sql"
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
	return dbpool, err
}

func SqlConnect() (*sql.DB, error) {
	db, err := sql.Open("mysql", os.Getenv("BOT_DATABASE_URL"))
	return db, err
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
