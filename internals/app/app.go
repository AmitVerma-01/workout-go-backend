package app

import (
	"database/sql"
	"fmt"
	"go_beginner/internals/api"
	"go_beginner/internals/middleware"
	"go_beginner/internals/store"
	"go_beginner/migrations"
	"log"
	"net/http"
	"os"
)

type Application struct {
	Logger *log.Logger
	WorkoutHandler *api.WorkoutHandler
	UserHandler *api.UserHandler
	DB *sql.DB 
	TokenHandler *api.TokenHandler
	Middleware middleware.UserMiddleware
}
 
func NewApplication() (*Application , error) {
	logger := log.New(os.Stdout, "", log.Ldate | log.Ltime )
	pgDB , err := store.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(fmt.Errorf("failed to apply migrations: %w", err))
	}
	workoutStore := store.NewPostgresWorkoutStore(pgDB)
	userStore := store.NewPostgresUserStore(pgDB)
	tokenStore := store.NewPostgresTokenStore(pgDB)
	middlewareHandler := middleware.UserMiddleware{UserStore: userStore}
	app := &Application{
		Logger: logger,
		WorkoutHandler:  api.NewWorkoutHandler(workoutStore, logger),
		UserHandler: api.NewUserHandler(userStore, logger),
		TokenHandler: api.NewTokenHandler(tokenStore, userStore, logger),
		DB: pgDB,
		Middleware: middlewareHandler,
	}
	return app, nil
}
   
func (a *Application) HealthCheck(w http.ResponseWriter , r *http.Request){
	fmt.Fprintf(w, "All good!\n")
}