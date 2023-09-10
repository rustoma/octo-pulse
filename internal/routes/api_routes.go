package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rustoma/octo-pulse/internal/api"
	"github.com/rustoma/octo-pulse/internal/controllers"
)

type ApiControllers struct {
	Auth *controllers.AuthController
}

func NewApiRoutes(controllers ApiControllers) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	//r.Use(app.enableCORS)

	r.Route("/api/v1", func(r chi.Router) {
		//AUTH
		r.Post("/login", api.MakeHTTPHandler(controllers.Auth.HandleLogin))
		//r.Post("/logout", api.MakeHTTPHandler(controllers.Auth.HandleLogout))
		r.Post("/refresh", api.MakeHTTPHandler(controllers.Auth.HandleRefreshToken))
	})
	//r.Get("/assets/images/*", api.MakeHTTPHandler(controllers.HandleGetImage))

	return r
}
