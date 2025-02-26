package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"shortened_links_service/internal/entities"
	"shortened_links_service/internal/handlers"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockShortenerService — мок для ShortenerService
type MockShortenerService struct {
	mock.Mock
}

// GetShortLink — заглушка для метода получения сокращённой ссылки
func (m *MockShortenerService) GetShortLink(originalLink string) (string, error) {
	args := m.Called(originalLink)
	return args.String(0), args.Error(1)
}

// GetOriginalLink — заглушка для метода получения оригинального URL
func (m *MockShortenerService) GetOriginalLink(shortLink string) (string, error) {
	args := m.Called(shortLink)
	return args.String(0), args.Error(1)
}

// TestGetShortLinkIntegration проверяет работу метода GetShortLink хендлера с подключением к http серверу
func TestGetShortLinkIntegration(t *testing.T) {
	mockService := new(MockShortenerService)
	handler := handlers.RegisterShortenerHandler(mockService)

	testOriginalLink := "https://example.com"
	expectedShortLink := "abcd123456"

	// Создаём тестовый HTTP-запрос
	requestBody, _ := json.Marshal(entities.Link{OriginalLink: testOriginalLink})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/shorten", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Создаем маршрутизатор и регистрируем маршрут с параметром
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/shorten", handler.GetShortLink())

	// Создаём тестовый HTTP-ответ
	respRec := httptest.NewRecorder()

	mockService.On("GetShortLink", testOriginalLink).Return(expectedShortLink, nil)
	mux.ServeHTTP(respRec, req)

	// Проверка кода ответа
	assert.Equalf(t, http.StatusOK, respRec.Code, "Ожидался статус 200, но получен %d", respRec.Code)

	// Проверка JSON-ответа
	var response map[string]string
	err := json.Unmarshal(respRec.Body.Bytes(), &response)
	assert.NoErrorf(t, err, "Ошибка парсинга JSON-ответа: %v", err)

	// Проверка получения сокращенной ссылки из тела ответа
	shortLink := response["short_link"]
	require.Equal(t, expectedShortLink, shortLink)

	mockService.AssertExpectations(t)
}

// TestGetOriginalLinkIntegration проверяет работу метода GetOriginalLink хендлера с подключением к http серверу
func TestGetOriginalLinkIntegration(t *testing.T) {
	mockService := new(MockShortenerService)
	handler := handlers.RegisterShortenerHandler(mockService)

	expectedOriginalLink := "https://example.com"
	testShortLink := "abcd123456"

	// Создаём тестовый HTTP-запрос
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/%s", testShortLink), nil)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/{short_link}", handler.GetOriginalLink())

	// Создаём тестовый HTTP-ответ
	respRec := httptest.NewRecorder()

	mockService.On("GetOriginalLink", testShortLink).Return(expectedOriginalLink, nil)
	mux.ServeHTTP(respRec, req)

	// Проверка кода ответа
	assert.Equalf(t, http.StatusOK, respRec.Code, "Ожидался статус 200, но получен %d", respRec.Code)

	// Проверка JSON-ответа
	var response map[string]string
	err := json.Unmarshal(respRec.Body.Bytes(), &response)
	assert.NoErrorf(t, err, "Ошибка парсинга JSON-ответа: %v", err)

	// Проверка получения оригинальной ссылки из тела ответа
	originalLink := response["original_link"]
	require.Equal(t, expectedOriginalLink, originalLink)

	mockService.AssertExpectations(t)
}
