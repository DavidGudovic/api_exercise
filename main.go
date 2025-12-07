package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/DavidGudovic/api_exercise/internal/app"
	"github.com/DavidGudovic/api_exercise/internal/routes"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "http server port")
	flag.Parse()

	application, err := app.NewApplication()

	if err != nil {
		panic(err)
	}

	r := routes.SetupRoutes(application)

	defer func() { _ = application.DB.Close() }()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      r,
	}

	application.Logger.Println("Server listening on port", port)

	application.Logger.Fatal(
		server.ListenAndServe(),
	)
}
