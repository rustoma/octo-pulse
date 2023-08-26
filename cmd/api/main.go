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

var (
	logger = lr.Logger
)

func main() {

	//Init logger
	logFile, err := lr.InitLogger()
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
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		logger.Fatal().Err(err)
	}
	log.Println("Connected to the DB")
	defer dbpool.Close()

}
