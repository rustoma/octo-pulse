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
	Auth    *controllers.AuthController
	Article *controllers.ArticleController
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
		//AUTH
		r.Post("/login", api.MakeHTTPHandler(controllers.Auth.HandleLogin))
		//r.Post("/logout", api.MakeHTTPHandler(controllers.Auth.HandleLogout))
		r.Post("/refresh", api.MakeHTTPHandler(controllers.Auth.HandleRefreshToken))
	})
	//r.Get("/assets/images/*", api.MakeHTTPHandler(controllers.HandleGetImage))

	r.Route("/api/v1/dashboard", func(r chi.Router) {
		r.Use(middlewares.RequireAuth())
		r.Use(middlewares.RequireApiKey)
		r.Get("/article/{id}/generate-description", api.MakeHTTPHandler(controllers.Article.HandleGenerateDescritption))
	})

	return r
}
