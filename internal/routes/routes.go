package routes

import (
	"github.com/DavidGudovic/api_exercise/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(application *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(
			application.Middleware.Authenticate,
			application.Middleware.RequireAuthenticatedUser,
		)

		r.Get("/workouts", application.WorkoutHandler.HandleGetAllWorkouts)
		r.Post("/workouts", application.WorkoutHandler.HandleCreateWorkout)
		r.Get("/workouts/{workoutID}", application.WorkoutHandler.HandleGetWorkoutByID)
		r.Put("/workouts/{workoutID}", application.WorkoutHandler.HandleUpdateWorkout)
		r.Delete("/workouts/{workoutID}", application.WorkoutHandler.HandleDeleteWorkout)
	})

	r.Post("/tokens/authentication", application.TokenHandler.HandleCreateToken)
	r.Post("/users/register", application.UserHandler.HandleRegisterUser)
	r.Get("/health", application.HealthCheck)

	return r
}
