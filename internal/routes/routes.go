package routes

import (
	"github.com/DavidGudovic/api_exercise/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(application *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", application.HealthCheck)

	r.Get("/workouts/{workoutID}", application.WorkoutHandler.HandleGetWorkoutByID)
	r.Post("/workouts/{workoutID}", application.WorkoutHandler.HandleCreateWorkout)
	r.Put("/workouts/{workoutID}", application.WorkoutHandler.HandleUpdateWorkout)
	r.Delete("/workouts/{workoutID}", application.WorkoutHandler.HandleDeleteWorkout)

	return r
}
