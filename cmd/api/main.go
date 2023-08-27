package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/rustoma/octo-pulse/internal/app"
	"github.com/rustoma/octo-pulse/internal/controllers"
	lr "github.com/rustoma/octo-pulse/internal/logger"
	"github.com/rustoma/octo-pulse/internal/routes"
	"github.com/rustoma/octo-pulse/internal/services"
	"github.com/rustoma/octo-pulse/internal/storage"
	postgresstore "github.com/rustoma/octo-pulse/internal/storage/postgresStore"
)

func main() {
	//Init logger
	logger, logFile, err := lr.NewLogger()
	defer logFile.Close()

	c := app.NewAppCtx(logger)

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

	var (
		//Storage
		userStore = postgresstore.NewUserStore(c, dbpool)
		store     = &storage.Store{
			User: userStore,
		}
		//Services
		authService = services.NewAuthService(c, userStore)
		//Controllers
		authController = controllers.NewAuthController(c, authService)
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
