package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rustoma/octo-pulse/internal/ai"
	"github.com/rustoma/octo-pulse/internal/db"
	"github.com/rustoma/octo-pulse/internal/services"
	postgresstore "github.com/rustoma/octo-pulse/internal/storage/postgresStore"
	ts "github.com/rustoma/octo-pulse/internal/tasks"

	lr "github.com/rustoma/octo-pulse/internal/logger"
)

var logger *zerolog.Logger
var logFile *os.File

func main() {
	defer logFile.Close()
	dbpool, err := db.Connect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		logger.Fatal().Err(err).Msg("")
	}
	defer dbpool.Close()

	var (
		ai             = ai.NewAI()
		store          = postgresstore.NewPostgresStorage(dbpool)
		articleService = services.NewArticleService(store.Article, ai)
		tasks          = ts.NewTasks(articleService)
	)

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: os.Getenv("REDIS_ADDR"), Password: os.Getenv("REDIS_PASSWORD")},
		asynq.Config{
			Concurrency: 2,
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(ts.TypeArticleGenerateDescription, tasks.Article.HandleGenerateDescription)
	if err := srv.Run(mux); err != nil {
		logger.Fatal().Msgf("could not run server: %v", err)
	}

}

func init() {
	l, lFile := lr.NewLogger()
	defer logFile.Close()
	logger = l
	logFile = lFile

	if err := godotenv.Load(filepath.Join(".", ".env")); err != nil {
		logger.Fatal().Msg("Error loading .env file")
	}
}
