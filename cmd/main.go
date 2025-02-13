package main

import (
	"log"
	"net/http"
	"shortened_links_service/internal/config"
	"shortened_links_service/internal/entities"
	"shortened_links_service/internal/handlers"
	"shortened_links_service/internal/services"
	"shortened_links_service/internal/storage"
	"time"
)

func main() {
	var (
		store storage.StorageInterface
		err   error
	)

	go handlers.СleanupVisitors()

	switch config.Mode {
	case "in-memory":
		store = storage.NewMemoryStore()
		log.Println("Using in-memory storage")
	case "postgres":
		entities.Db, err = storage.NewDatabaseStore(config.PsqlUrl)
		if err != nil {
			log.Fatal(err.Error())
		}
		defer entities.Db.Close()

		store = storage.NewDatabaseConection(entities.Db)
		log.Println("Using PostgreSQL store")
	default:
		log.Fatalf("config.Mode is empty in /internal/config/setting.go")
	}

	service := services.NewShortenerService(store)
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/shorten", handlers.ShortenLink(service))
	mux.HandleFunc("GET /api/v1/{short_link}", handlers.RerouteLink(service))

	serv := &http.Server{
		Addr:         config.Port,
		Handler:      handlers.LimiterMiddleware(mux),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	log.Println("Service is running ...")
	if err := serv.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}
}
