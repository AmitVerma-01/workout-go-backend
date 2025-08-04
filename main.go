package main

import (
	"flag"
	"fmt"
	"go_beginner/internals/app"
	"go_beginner/internals/routes"
	"net/http"
	"time"
)
func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "Port to run the server on")
	flag.Parse()
	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
	defer app.DB.Close()
	app.Logger.Printf("Application started successfully on port %d", port )
	r := routes.SetipRoutes(app)
	server := http.Server{
		Addr:  fmt.Sprintf(":%d", port),
		Handler: r,
		IdleTimeout: time.Minute,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatalf("Failed to start server: %v", err)
	} else {
		app.Logger.Println("Server started on port 8080")
	}
}
