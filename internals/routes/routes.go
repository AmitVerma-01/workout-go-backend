package routes

import (
	"go_beginner/internals/app"

	"github.com/go-chi/chi/v5"
)

func SetipRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/health", app.HealthCheck)
	r.Get("/workout/{id}", app.WorkoutHandler.HandleGetWorkoutByID)
	r.Post("/workout", app.WorkoutHandler.HandleCreateWorkout)
	r.Patch("/workout/{id}", app.WorkoutHandler.HandleUpdateWorkout)
	r.Delete("/workout/{id}", app.WorkoutHandler.HandleDeleteWorkout)
	r.Get("/workouts", app.WorkoutHandler.HandleGetAllWorkouts)

	r.Post("/user", app.UserHandler.HandleCreateUser)
	r.Get("/user/{id}", app.UserHandler.HandleGetUserByID)
	r.Patch("/user/{id}", app.UserHandler.HandleUpdateUser)
	r.Delete("/user/{id}", app.UserHandler.HandleDeleteUser)
	// r.Get("/users", app.UserHandler.÷)
	
	r.Post("/login", app.TokenHandler.HandleCreateToken)
	// r.Post("/register", app.UserHandler.HandleCreateUser)

	return r 
}