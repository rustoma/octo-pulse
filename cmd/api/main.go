package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rustoma/octo-pulse/internal/ai"
	"github.com/rustoma/octo-pulse/internal/controllers"
	"github.com/rustoma/octo-pulse/internal/db"
	lr "github.com/rustoma/octo-pulse/internal/logger"
	"github.com/rustoma/octo-pulse/internal/routes"
	"github.com/rustoma/octo-pulse/internal/services"
	postgresstore "github.com/rustoma/octo-pulse/internal/storage/postgresStore"
	ts "github.com/rustoma/octo-pulse/internal/tasks"
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
	defer dbpool.Close()

	var (
		//AI
		ai = ai.NewAI()
		//Storage
		store = postgresstore.NewPostgresStorage(dbpool)
		//Services
		authService    = services.NewAuthService(store.User)
		articleService = services.NewArticleService(store.Article, ai)
		domainService  = services.NewDomainService(store.Domain)
		//Tasks
		tasks         = ts.NewTasks(articleService)
		taskInspector = ts.NewTaskInspector()
		//Controllers
		authController    = controllers.NewAuthController(authService)
		articleController = controllers.NewArticleController(articleService, tasks.Article)
		taskController    = controllers.NewTaskController(taskInspector)
		domainController  = controllers.NewDomainController(domainService)
		apiControllers    = routes.ApiControllers{
			Auth:    authController,
			Article: articleController,
			Task:    taskController,
			Domain:  domainController,
		}
		apiServices = routes.ApiServices{
			Auth: authService,
		}
	)

	_ = store

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
