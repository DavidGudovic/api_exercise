package routes

import (
	"github.com/DavidGudovic/api_exercise/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(application *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", application.HealthCheck)

	return r
}
