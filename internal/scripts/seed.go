package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/rustoma/octo-pulse/internal/fixtures"
	lr "github.com/rustoma/octo-pulse/internal/logger"
	"github.com/rustoma/octo-pulse/internal/services"
	postgresstore "github.com/rustoma/octo-pulse/internal/storage/postgresStore"
)

func main() {
	//Init logger
	logger, logFile := lr.NewLogger()
	defer logFile.Close()

	//Init .env
	err := godotenv.Load()
	if err != nil {
		logger.Fatal().Msg("Error loading .env file")
	}

	//Init DB
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("SEED_DATABASE_URL"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		logger.Fatal().Err(err).Send()
	}

	logger.Info().Msg("Connected to the DB")
	defer dbpool.Close()

	var (
		store       = postgresstore.NewPostgresStorage(dbpool)
		authService = services.NewAuthService(store.User)
		fixtures    = fixtures.NewFixtures(authService)
	)

	adminRole := fixtures.CreateRole("Admin")
	editorRole := fixtures.CreateRole("Editor")

	_, err = store.Role.InsertRole(adminRole)

	if err != nil {
		panic(err)
	}

	_, err = store.Role.InsertRole(editorRole)

	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	adminUser := fixtures.CreateUser("admin@admin.com", "admin", 1)
	editorUser := fixtures.CreateUser("editor@editor.com", "editor", 2)

	_, err = store.User.InsertUser(adminUser)

	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	_, err = store.User.InsertUser(editorUser)

	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	for i := 0; i < 3; i++ {
		domain := fixtures.CreateDomain(fmt.Sprintf("exampleDomain%d.com", i+1))
		_, err = store.Domain.InsertDomain(domain)

		if err != nil {
			logger.Fatal().Err(err).Send()
		}
	}

	homeCategory := fixtures.CreateCategory("Home")
	generalCategory := fixtures.CreateCategory("General")
	newsCategory := fixtures.CreateCategory("News")

	_ = homeCategory
	_ = generalCategory
	_ = newsCategory

	john := fixtures.CreateAuthor("John", "Doe", "Lorem ipsum dolor", "https://thispersondoesnotexist.com/")
	jane := fixtures.CreateAuthor("Jane", "Doe", "Lorem ipsum dolor", "https://thispersondoesnotexist.com/")

	_ = john
	_ = jane

	for i := 0; i < 10; i++ {
		title := fmt.Sprintf("Home Article-%d", i+1)
		desc := "Lorem ipsum dolor"
		imageUrl := ""
		isPubished := true
		authorId := 1
		categoryId := 1
		domainId := 1
		article := fixtures.CreateArticle(title, desc, imageUrl, isPubished, authorId, categoryId, domainId)

		_ = article
	}

	for i := 0; i < 20; i++ {
		title := fmt.Sprintf("General Article %d", i+1)
		desc := "Lorem ipsum dolor"
		imageUrl := ""
		isPubished := true
		authorId := 1
		categoryId := 1
		domainId := 1
		article := fixtures.CreateArticle(title, desc, imageUrl, isPubished, authorId, categoryId, domainId)

		_ = article
	}

	for i := 0; i < 15; i++ {
		title := fmt.Sprintf("News Article %d", i+1)
		desc := "Lorem ipsum dolor"
		imageUrl := ""
		isPubished := true
		authorId := 1
		categoryId := 1
		domainId := 1
		article := fixtures.CreateArticle(title, desc, imageUrl, isPubished, authorId, categoryId, domainId)

		_ = article
	}

}
