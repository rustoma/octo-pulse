package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	lr "github.com/rustoma/octo-pulse/internal/logger"
	"github.com/rustoma/octo-pulse/internal/services"
	postgresstore "github.com/rustoma/octo-pulse/internal/storage/postgresStore"
	"os"
	"path/filepath"
)

func main() {

	dirPath := os.Args[1]

	//Init logger
	logger, logFile := lr.NewLogger()
	defer logFile.Close()

	//Init .env
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logger.Fatal().Msg("Error loading .env file")
	}

	if err := godotenv.Load(filepath.Join(dir, ".env")); err != nil {
		logger.Fatal().Msg("Error loading .env file")
	}

	//Init DB
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		logger.Fatal().Err(err).Send()
	}

	logger.Info().Msg("Connected to the DB")
	defer dbpool.Close()

	var (
		store       = postgresstore.NewPostgresStorage(dbpool)
		fileService = services.NewFileService(store.Article, store.Domain, store.Category, store.Image)
	)

	logger.Info().Msg("Renaming files from the directory: " + dirPath)
	fileService.RenameFilesUsingSlug(dirPath)
	logger.Info().Msg("Files renamed successfully: " + dirPath)

	logger.Info().Msg("Scanning images from the directory: " + dirPath)
	err = fileService.InsertJPGImagesFromDir(dirPath, 2)
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	logger.Info().Msg("Images added successfully")

}
