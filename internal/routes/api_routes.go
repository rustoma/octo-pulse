package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rustoma/octo-pulse/internal/api"
	"github.com/rustoma/octo-pulse/internal/controllers"
	m "github.com/rustoma/octo-pulse/internal/middleware"
	"github.com/rustoma/octo-pulse/internal/services"
	"github.com/rustoma/octo-pulse/internal/tasks"
)

type ApiControllers struct {
	Auth      *controllers.AuthController
	Article   *controllers.ArticleController
	Task      *controllers.TaskController
	Domain    *controllers.DomainController
	Category  *controllers.CategoryController
	File      *controllers.FileController
	Image     *controllers.ImageController
	BasicPage *controllers.BasicPageController
	Email     *controllers.EmailController
	Author    *controllers.AuthorController
}

type ApiServices struct {
	Auth services.AuthService
}

func NewApiRoutes(controllers ApiControllers, services ApiServices, tasks *tasks.Tasks) http.Handler {
	middlewares := m.NewMiddleware(services.Auth)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middlewares.EnableCORS)

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middlewares.RequireApiKey)

		r.Get("/articles", api.MakeHTTPHandler(controllers.Article.HandleGetArticles))
		r.Get("/articles/{id}", api.MakeHTTPHandler(controllers.Article.HandleGetArticle))

		r.Get("/domains/{id}", api.MakeHTTPHandler(controllers.Domain.HandleGetDomainPublicData))

		r.Get("/categories", api.MakeHTTPHandler(controllers.Category.HandleGetCategories))
		r.Get("/domain-categories/{id}", api.MakeHTTPHandler(controllers.Category.HandleGetDomainCategories))

		r.Get("/basic-pages", api.MakeHTTPHandler(controllers.BasicPage.HandleGetBasicPages))
		r.Get("/basic-pages/slug/{slug}", api.MakeHTTPHandler(controllers.BasicPage.HandleGetBasicPageBySlug))

		r.Post("/emails", api.MakeHTTPHandler(controllers.Email.HandleSendEmail))
	})

	r.Get("/assets/images/*", api.MakeHTTPHandler(controllers.Image.HandleGetImageByPath))

	r.Route("/api/v1/dashboard/auth", func(r chi.Router) {
		r.Post("/login", api.MakeHTTPHandler(controllers.Auth.HandleLogin))
		r.Post("/logout", api.MakeHTTPHandler(controllers.Auth.HandleLogout))
		r.Post("/refresh", api.MakeHTTPHandler(controllers.Auth.HandleRefreshToken))
	})

	r.Route("/api/v1/dashboard", func(r chi.Router) {
		r.Use(middlewares.RequireAuth())

		r.Get("/domains", api.MakeHTTPHandler(controllers.Domain.HandleGetDomains))
		r.Get("/domains/{id}", api.MakeHTTPHandler(controllers.Domain.HandleGetDomain))

		r.Get("/domain-categories/{id}", api.MakeHTTPHandler(controllers.Category.HandleGetDomainCategories))

		r.Get("/articles", api.MakeHTTPHandler(controllers.Article.HandleGetArticles))
		r.Post("/articles", api.MakeHTTPHandler(controllers.Article.HandleCreateArticle))
		r.Get("/articles/{id}", api.MakeHTTPHandler(controllers.Article.HandleGetArticle))
		r.Put("/articles/{id}", api.MakeHTTPHandler(controllers.Article.HandleUpdateArticle))
		r.Delete("/articles/{id}", api.MakeHTTPHandler(controllers.Article.HandleDeleteArticle))
		r.Post("/articles/{id}/generate-description", api.MakeHTTPHandler(controllers.Article.HandleGenerateDescritption))
		r.Get("/articles/{id}/remove-duplicates", api.MakeHTTPHandler(controllers.Article.HandleRemoveDuplicatesFromArticle))
		r.Post("/articles/generate", api.MakeHTTPHandler(controllers.Article.HandleGenerateArticles))

		r.Get("/categories", api.MakeHTTPHandler(controllers.Category.HandleGetCategories))
		r.Get("/categories/{id}", api.MakeHTTPHandler(controllers.Category.HandleGetCategory))

		r.Get("/authors", api.MakeHTTPHandler(controllers.Author.HandleGetAuthors))

		r.Post("/files/articles", api.MakeHTTPHandler(controllers.File.HandleCreateArticles))

		r.Post("/tasks", api.MakeHTTPHandler(controllers.Task.HandleGetTasksInfo))

		r.Get("/images", api.MakeHTTPHandler(controllers.Image.HandleGetImages))
		r.Post("/images/category-id/{id}", api.MakeHTTPHandler(controllers.Image.HandleUploadImage))
		r.Get("/images/{id}", api.MakeHTTPHandler(controllers.Image.HandleGetImage))
		r.Get("/image-categories", api.MakeHTTPHandler(controllers.Image.HandleGetImageCategories))
		r.Get("/assets/images/*", api.MakeHTTPHandler(controllers.Image.HandleGetImageByPath))
	})

	return r
}
