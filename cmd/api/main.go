package main

import (
	"flag"
	"log"
	"log/slog"
	"net/http"

	"sequencesender/internal/handlers"
	"sequencesender/internal/services"
	"sequencesender/pkg/dbclient"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

var developmentFlag = flag.Bool("development", false, "Load environment variables from .env file for development")

func init() {
	flag.Parse()
}

func main() {
	// load environment variables only when development flag is passed
	if *developmentFlag {
		slog.Info("Running in development mode")

		if err := godotenv.Load(); err != nil {
			slog.Error("failed to load .env file", "error", err)
			return
		} else {
			slog.Info("environment vars loaded from .env file")
		}
	}

	dbClient, err := dbclient.NewSQLXConnection(dbclient.GetDBConnectionString())
	if err != nil {
		slog.Error(" initialize database connection failed ", slog.String("error", err.Error()))
		return
	}

	if err := dbClient.Ping(); err != nil {
		slog.Error("failed to ping database", slog.String("error", err.Error()))
		return
	}
	slog.Info("database connection established successfully")

	sequenceService := services.NewSequenceService(dbClient)

	sequenceHandler := handlers.NewSequenceHandler(sequenceService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(rr chi.Router) {
		sequenceHandler.RegisterRoutes(rr)
	})

	server := &http.Server{
		Addr:    ":3000",
		Handler: r,
	}

	slog.Info("server listening on :3000")
	log.Fatal(server.ListenAndServe())
}
