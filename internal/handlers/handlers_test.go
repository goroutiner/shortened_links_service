package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"shortened_links_service/internal/entities"
	"shortened_links_service/internal/handlers"
	"shortened_links_service/internal/services"
	"shortened_links_service/internal/storage"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestShortenLink тестирует обработчик handlers.ShortenLink() на корректность обработки ссылок
func TestShortenLink(t *testing.T) {
	store := storage.NewMemoryStore()
	service := services.NewShortenerService(store)
	handler := handlers.ShortenLink(service)

	testLink := "https://finance.ozon.ru"

	// Создаём тестовый HTTP-запрос
	requestBody, _ := json.Marshal(entities.Link{OriginalLink: testLink})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/shorten", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Создаем маршрутизатор и регистрируем маршрут с параметром
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/shorten", handler)

	// Создаём тестовый HTTP-ответ
	respRec := httptest.NewRecorder()

	mux.ServeHTTP(respRec, req)

	// Проверка код ответа
	assert.Equalf(t, http.StatusOK, respRec.Code, "Ожидался статус 200, но получен %d", respRec.Code)

	// Проверка JSON-ответа
	var response map[string]string
	err := json.Unmarshal(respRec.Body.Bytes(), &response)
	assert.NoErrorf(t, err, "Ошибка парсинга JSON-ответа: %v", err)

	// Проверка получения сокращенной ссылки из тела ответа
	shortLink := response["short_link"]
	require.NotEmpty(t, shortLink, "Сокращённая ссылка должна быть не пустой")
	require.Equal(t, 10, len(shortLink), "Число символов в строке не равно 10")

	// Проверка получения оригинальной ссылки по сокращённой
	originalLink, err := service.GetOriginalLink(shortLink)
	assert.NoError(t, err, "Ошибка при получении оригинального URL")
	require.NotEmpty(t, originalLink, "Оригинальная ссылка должна быть не пустой")
	assert.Equalf(t, testLink, originalLink, `Ожидалась оригинальная URL "%s", но получена "%s"`, testLink, originalLink)
}

// TestRerouteLink тестирует обработчик handlers.RerouteLink() на корректность обработки ссылок
func TestRerouteLink(t *testing.T) {
	store := storage.NewMemoryStore()
	service := services.NewShortenerService(store)
	handler := handlers.RerouteLink(service)

	testLink := "https://finance.ozon.ru"

	shortLink, err := service.GetShortLink(testLink)
	assert.NoError(t, err, "Ошибка добавления тестовой ссылки в хранилище")

	// Создаём тестовый HTTP-запрос
	req := httptest.NewRequest(http.MethodGet, "/api/v1/"+shortLink, nil)

	// Создаем маршрутизатор и регистрируем маршрут с параметром
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/{short_link}", handler)

	// Создаём тестовый HTTP-ответ
	respRec := httptest.NewRecorder()

	// Вызываем обработчик
	mux.ServeHTTP(respRec, req)

	// Проверка кода ответа
	assert.Equalf(t, http.StatusOK, respRec.Code, "Ожидался статус 200, но получен %d", respRec.Code)

	// Проверка JSON-ответа
	var response map[string]string
	err = json.Unmarshal(respRec.Body.Bytes(), &response)
	assert.NoErrorf(t, err, "Ошибка парсинга JSON-ответа: %v", err)

	// Проверка получения оригинальной ссылки из тела ответа
	originalLink := response["original_link"]
	require.NotEmpty(t, originalLink, "Оригинальная ссылка должна быть не пустой")
	assert.Equalf(t, testLink, originalLink, `Ожидался original_link "%s", но получен "%s"`, testLink, originalLink)
}
