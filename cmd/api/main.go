package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rustoma/octo-pulse/internal/ai"
	"github.com/rustoma/octo-pulse/internal/controllers"
	"github.com/rustoma/octo-pulse/internal/db"
	lr "github.com/rustoma/octo-pulse/internal/logger"
	"github.com/rustoma/octo-pulse/internal/routes"
	"github.com/rustoma/octo-pulse/internal/services"
	"github.com/rustoma/octo-pulse/internal/storage"
	postgresstore "github.com/rustoma/octo-pulse/internal/storage/postgresStore"
	sqlstore "github.com/rustoma/octo-pulse/internal/storage/sqlStore"
	ts "github.com/rustoma/octo-pulse/internal/tasks"
	"github.com/rustoma/octo-pulse/internal/validator"
)

var logger *zerolog.Logger
var logFile *os.File

func main() {
	defer logFile.Close()
	//Init DB
	dbpool, err := db.Connect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		logger.Fatal().Err(err).Msg("")
	}
	logger.Info().Msg("Connected to the Postgress DB")
	defer dbpool.Close()

	//Init SQL DB
	db, err := db.SqlConnect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to SQL database: %v\n", err)
		logger.Fatal().Err(err).Msg("")
	}
	logger.Info().Msg("Connected to the SQL DB")
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	var (
		//AI
		ai = ai.NewAI()
		//Storage
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
		//Validator
		validator = validator.NewValidator()
		//Services
		authService     = services.NewAuthService(store.User)
		articleService  = services.NewArticleService(store.Article, validator.Article, ai)
		domainService   = services.NewDomainService(store.Domain)
		categoryService = services.NewCategoryService(store.Category)
		scrapperService = services.NewScrapperService(store.Scrapper)
		//Tasks
		tasks         = ts.NewTasks(articleService)
		taskInspector = ts.NewTaskInspector()
		//Controllers
		authController     = controllers.NewAuthController(authService)
		articleController  = controllers.NewArticleController(articleService, tasks.Article)
		taskController     = controllers.NewTaskController(taskInspector)
		domainController   = controllers.NewDomainController(domainService)
		categoryController = controllers.NewCategoryController(categoryService)
		apiControllers     = routes.ApiControllers{
			Auth:     authController,
			Article:  articleController,
			Task:     taskController,
			Domain:   domainController,
			Category: categoryController,
		}
		apiServices = routes.ApiServices{
			Auth: authService,
		}
	)

	question, err := scrapperService.GetQuestion(354)

	if err != nil {
		logger.Err(err).Send()
	}

	logger.Info().Interface("question: ", question).Send()

	//start a web server
	log.Println("Starting application on port", os.Getenv("PORT"))
	err = http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), routes.NewApiRoutes(apiControllers, apiServices, tasks))
	if err != nil {
		log.Fatal(err)
	}
}

func init() {

	//Init logger
	l, lFile := lr.NewLogger()
	defer logFile.Close()
	logger = l
	logFile = lFile

	//Init .env
	if err := godotenv.Load(filepath.Join(".", ".env")); err != nil {
		logger.Fatal().Msg("Error loading .env file")
	}
}
