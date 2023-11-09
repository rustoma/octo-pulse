package main

import (
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rustoma/octo-pulse/internal/ai"
	"github.com/rustoma/octo-pulse/internal/db"
	"github.com/rustoma/octo-pulse/internal/services"
	"github.com/rustoma/octo-pulse/internal/storage"
	postgresstore "github.com/rustoma/octo-pulse/internal/storage/postgresStore"
	sqlstore "github.com/rustoma/octo-pulse/internal/storage/sqlStore"
	ts "github.com/rustoma/octo-pulse/internal/tasks"
	"github.com/rustoma/octo-pulse/internal/validator"

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

	// Init SQL DB
	db, err := db.SqlConnect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to SQL database: %v\n", err)
		logger.Fatal().Err(err).Msg("")
	}
	logger.Info().Msg("Connected to the SQL DB")

	var (
		validator      = validator.NewValidator()
		ai             = ai.NewAI()
		postgressStore = postgresstore.NewPostgresStorage(dbpool)
		sqlStore       = sqlstore.NewSqlStorage(db)
		store          = storage.Store{
			User:              postgressStore.User,
			Role:              postgressStore.Role,
			Domain:            postgressStore.Domain,
			Category:          postgressStore.Category,
			Author:            postgressStore.Author,
			Article:           postgressStore.Article,
			CategoriesDomains: postgressStore.CategoriesDomains,
			Scrapper:          sqlStore.Scrapper,
		}
		articleService  = services.NewArticleService(store.Article, validator.Article, ai)
		domainService   = services.NewDomainService(store.Domain)
		categoryService = services.NewCategoryService(store.Category, store.CategoriesDomains)
		scrapperService = services.NewScrapperService(store.Scrapper, validator.Scrapper)
		tasks           = ts.NewTasks(articleService, domainService, scrapperService, categoryService, ai)
	)

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: os.Getenv("REDIS_ADDR"), Password: os.Getenv("REDIS_PASSWORD")},
		asynq.Config{
			Concurrency: 1,
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(ts.TypeArticleGenerateDescription, tasks.Article.HandleGenerateDescription)
	mux.HandleFunc(ts.TypeArticleGenerateArticles, tasks.Article.HandleGenerateArticles)
	mux.HandleFunc(ts.TypeScrapperUpdateQuestion, tasks.Scrapper.HandleUpdateQuestionTask)
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
