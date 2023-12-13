package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/rustoma/octo-pulse/internal/fixtures"
	lr "github.com/rustoma/octo-pulse/internal/logger"
	"github.com/rustoma/octo-pulse/internal/services"
	postgresstore "github.com/rustoma/octo-pulse/internal/storage/postgresStore"
	"math/rand"
	"os"
	"time"
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

	homeDesignDomain := fixtures.CreateDomain("homedesign.com")
	homeDesignDomainId, err := store.Domain.InsertDomain(homeDesignDomain)

	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	newsDomain := fixtures.CreateDomain("hotnews.com")
	newsDomainId, err := store.Domain.InsertDomain(newsDomain)

	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	homeCategory := fixtures.CreateCategory("Home")

	homeCategoryId, err := store.Category.InsertCategory(homeCategory)

	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	generalCategory := fixtures.CreateCategory("General")

	generalCategoryId, err := store.Category.InsertCategory(generalCategory)

	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	newsCategory := fixtures.CreateCategory("News")

	newsCategoryId, err := store.Category.InsertCategory(newsCategory)

	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	john := fixtures.CreateAuthor("John", "Doe", "Lorem ipsum dolor", "https://thispersondoesnotexist.com/")

	johnId, err := store.Author.InsertAuthor(john)

	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	jane := fixtures.CreateAuthor("Jane", "Doe", "Lorem ipsum dolor", "https://thispersondoesnotexist.com/")

	janeId, err := store.Author.InsertAuthor(jane)

	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	err = store.CategoriesDomains.AsignCategoryToDomain(homeCategoryId, homeDesignDomainId)

	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	err = store.CategoriesDomains.AsignCategoryToDomain(newsCategoryId, homeDesignDomainId)

	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	err = store.CategoriesDomains.AsignCategoryToDomain(newsCategoryId, newsDomainId)

	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	constructionCategory := fixtures.CreateImageCategory("Construction")
	_, err = store.ImageCategory.InsertCategory(constructionCategory)

	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	imageFirstName := "mezczyzna-siedzacy-na-zewnatrz-i-glaszczacy-swojego-kota"
	imageSecondName := "odkryty-patio-z-krzeslem-i-stolem"
	imageThirdName := "puste-krzeslo-drewniane-w-salonie"
	imageFirst := fixtures.CreateImage(
		imageFirstName,
		"/assets/image/mezczyzna-siedzacy-na-zewnatrz-i-glaszczacy-swojego-kota.jpg",
		184551,
		".jpg",
		1500,
		998,
		"",
		1,
	)
	imageSecond := fixtures.CreateImage(
		imageSecondName,
		"/assets/image/odkryty-patio-z-krzeslem-i-stolem.jpg",
		371302,
		".jpg",
		1500,
		1000,
		"",
		1,
	)
	imageThird := fixtures.CreateImage(
		imageThirdName,
		"/assets/image/puste-krzeslo-drewniane-w-salonie.jpg",
		296332,
		".jpg",
		1500,
		1203,
		"",
		1,
	)

	_, err = store.Image.InsertImage(imageFirst)
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	_, err = store.Image.InsertImage(imageSecond)
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	_, err = store.Image.InsertImage(imageThird)
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	for i := 0; i < 10; i++ {
		rand.NewSource(time.Now().UnixNano())
		n := 1 + rand.Intn(3-1+1)

		title := fmt.Sprintf("Home Article %d", i+1)
		desc := "Lorem ipsum dolor"
		thumbnail := n
		isPubished := true
		authorId := janeId
		categoryId := homeCategoryId
		domainId := homeDesignDomainId
		featured := false
		article := fixtures.CreateArticle(title, desc, thumbnail, isPubished, authorId, categoryId, domainId, featured)

		_, err = store.Article.InsertArticle(article)

		if err != nil {
			logger.Fatal().Err(err).Send()
		}
	}

	for i := 0; i < 20; i++ {
		rand.NewSource(time.Now().UnixNano())
		n := 1 + rand.Intn(3-1+1)

		title := fmt.Sprintf("General Article %d", i+1)
		desc := "Lorem ipsum dolor"
		thumbnail := n
		isPubished := true
		authorId := johnId
		categoryId := generalCategoryId
		domainId := homeDesignDomainId
		featured := false
		article := fixtures.CreateArticle(title, desc, thumbnail, isPubished, authorId, categoryId, domainId, featured)

		_, err = store.Article.InsertArticle(article)

		if err != nil {
			logger.Fatal().Err(err).Send()
		}
	}

	for i := 0; i < 15; i++ {
		rand.NewSource(time.Now().UnixNano())
		n := 1 + rand.Intn(3-1+1)

		title := fmt.Sprintf("Clean Home Article %d", i+1)
		desc := "Lorem ipsum dolor"
		thumbnail := n
		isPubished := true
		authorId := johnId
		categoryId := homeCategoryId
		domainId := homeDesignDomainId
		featured := false
		article := fixtures.CreateArticle(title, desc, thumbnail, isPubished, authorId, categoryId, domainId, featured)

		_, err = store.Article.InsertArticle(article)

		if err != nil {
			logger.Fatal().Err(err).Send()
		}
	}

	for i := 0; i < 15; i++ {
		rand.NewSource(time.Now().UnixNano())
		n := 1 + rand.Intn(3-1+1)

		title := fmt.Sprintf("News Article %d", i+1)
		desc := "Lorem ipsum dolor"
		thumbnail := n
		isPubished := true
		authorId := johnId
		categoryId := newsCategoryId
		domainId := newsDomainId
		featured := false
		article := fixtures.CreateArticle(title, desc, thumbnail, isPubished, authorId, categoryId, domainId, featured)

		_, err = store.Article.InsertArticle(article)

		if err != nil {
			logger.Fatal().Err(err).Send()
		}
	}

}
