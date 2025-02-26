package main

import (
	"log"
	"net/http"
	"shortened_links_service/internal/config"
	"shortened_links_service/internal/entities"
	"shortened_links_service/internal/handlers"
	"shortened_links_service/internal/services"
	"shortened_links_service/internal/storage"
	"shortened_links_service/internal/storage/database"
	"shortened_links_service/internal/storage/memory"
	"time"
)

func main() {
	var (
		store storage.StorageInterface
		err   error
	)

	go handlers.СleanupVisitors()

	switch "in-memory" {
	case config.Mode:
		store = memory.NewMemoryStore()
		log.Println("Using in-memory storage")
	case "postgres":
		entities.Db, err = database.NewDatabaseStore(config.PsqlUrl)
		if err != nil {
			log.Fatal(err.Error())
		}
		defer entities.Db.Close()

		store = database.NewDatabaseConection(entities.Db)
		log.Println("Using PostgreSQL store")
	default:
		log.Fatalf("config.Mode is empty in /internal/config/setting.go")
	}

	service := services.NewShortenerService(store)
	hadler := handlers.RegisterShortenerHandler(service)
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/shorten", hadler.GetShortLink())
	mux.HandleFunc("GET /api/v1/{short_link}", hadler.GetOriginalLink())

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
