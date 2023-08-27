package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rustoma/octo-pulse/internal/controllers"
	"github.com/rustoma/octo-pulse/internal/db"
	lr "github.com/rustoma/octo-pulse/internal/logger"
	"github.com/rustoma/octo-pulse/internal/routes"
	"github.com/rustoma/octo-pulse/internal/services"
	"github.com/rustoma/octo-pulse/internal/storage"
	postgresstore "github.com/rustoma/octo-pulse/internal/storage/postgresStore"
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
		//Storage
		userStore = postgresstore.NewUserStore(dbpool)
		store     = &storage.Store{
			User: userStore,
		}
		//Services
		authService = services.NewAuthService(userStore)
		//Controllers
		authController = controllers.NewAuthController(authService)
		apiControllers = routes.ApiControllers{
			Auth: authController,
		}
	)

	_ = store

	//start a web server
	log.Println("Starting application on port", os.Getenv("PORT"))
	err = http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), routes.NewApiRoutes(apiControllers))
	if err != nil {
		log.Fatal(err)
	}
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
