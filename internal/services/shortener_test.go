package services_test

import (
	"shortened_links_service/internal/services"
	"shortened_links_service/internal/storage"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestShortenerService_InMemory тестирует работу сервиса "in-memory" режиме
func TestShortenerService_InMemory(t *testing.T) {
	store := storage.NewMemoryStore()
	service := services.NewShortenerService(store)

	testLink := "https://finance.ozon.ru"

	// Проверка создания сокращённой ссылки
	shortLink, err := service.GetShortLink(testLink)
	require.NoError(t, err, "Ошибка при создании сокращённой ссылки")
	assert.NotEmpty(t, shortLink, "Сокращённая ссылка должна быть не пустой")
	require.Equal(t, 10, len(shortLink), "Число символов в строке не равно 10")

	// Проверка получения оригинального URL по сокращённому
	originalLink, err := service.GetOriginalLink(shortLink)
	require.NoError(t, err, "Ошибка при получении оригинального URL")
	assert.Equal(t, testLink, originalLink, "Оригинальный URL должен совпадать")

	// Проверка обработки ошибки, если ссылки нет в базе
	_, err = service.GetOriginalLink("nonexistent")
	assert.Error(t, err, "Должна быть ошибка, если ссылка не найдена")

	// Проверка ошибки, при обработки пустой оригинальной ссылки
	testLink = ""
	_, err = service.GetShortLink(testLink)
	require.Error(t, err, "Должна быть при создании, если оригинальная ссылка пустая")

	// Проверка получения ошибки, при обработки некорректной оригинальной ссылки
	testLink = "http///_somelink_"
	_, err = service.GetShortLink(testLink)
	require.Error(t, err, "Должна быть при создании, если оригинальная ссылка некорректная")
}

// TestShortenerService_Postgres тестирует работу сервиса в "postgres" режиме
func TestShortenerService_Postgres(t *testing.T) {
	psqlUrl := "postgres://user:password@localhost:5432/test_db?sslmode=disable"
	db, err := storage.NewDatabaseStore(psqlUrl)
	require.NoError(t, err, "Ошибка подключения к PostgreSQL")
	defer db.Close()

	store := storage.NewDatabaseConection(db)
	service := services.NewShortenerService(store)

	testLink := "https://finance.ozon.ru"

	// Проверка создания сокращённой ссылки
	shortLink, err := service.GetShortLink(testLink)
	require.NoError(t, err)
	assert.NotEmpty(t, shortLink)
	require.Equal(t, 10, len(shortLink), "Число символов в строке не равно 10")

	// Проверка получения оригинального URL по сокращённому
	originalLink, err := service.GetOriginalLink(shortLink)
	require.NoError(t, err)
	assert.Equal(t, testLink, originalLink)

	// Проверка обработки ошибки, если ссылки нет в базе
	_, err = service.GetOriginalLink("nonexistent")
	assert.Error(t, err)

	// Проверка получения ошибки, при обработки пустой оригинальной ссылки
	testLink = ""
	_, err = service.GetShortLink(testLink)
	require.Error(t, err, "Должна быть при создании, если оригинальная ссылка пустая")

	// Проверка получения ошибки, при обработки некорректной оригинальной ссылки
	testLink = "http///_somelink_"
	_, err = service.GetShortLink(testLink)
	require.Error(t, err, "Должна быть при создании, если оригинальная ссылка не корректная")
}
