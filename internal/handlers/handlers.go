package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"shortened_links_service/internal/entities"
	"shortened_links_service/internal/services"
)

type ShortenerHandler struct {
	service services.ShortenerServiceInterface
}

func RegisterShortenerHandler(service services.ShortenerServiceInterface) *ShortenerHandler {
	return &ShortenerHandler{service: service}
}

// GetShortLink обрабатывает POST запрос с оригинальной ссылкой и возвращет сокращенную ссылку в формате json
func (s *ShortenerHandler) GetShortLink() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var link entities.Link

		if err := json.NewDecoder(r.Body).Decode(&link); err != nil {
			log.Println(err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		shortLink, err := s.service.GetShortLink(link.OriginalLink)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to generate short link", http.StatusInternalServerError)
			return
		}

		resp := map[string]string{"short_link": shortLink}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

// GetOriginalLink обрабатывает GET запрос с сокращенной ссылкой и возвращает оригинальную ссылку в формате json
func (s *ShortenerHandler) GetOriginalLink() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		shortURL := r.PathValue("short_link")

		originalURL, err := s.service.GetOriginalLink(shortURL)
		if err != nil {
			log.Println(err)
			http.NotFound(w, r)
			return
		}

		resp := map[string]string{"original_link": originalURL}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
