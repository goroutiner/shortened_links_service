package services_test

import (
	"errors"
	"shortened_links_service/internal/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockStorage — мок-хранилище, реализующее StorageInterface
type MockStorage struct {
	mock.Mock
}

// GetShortLink — заглушка для метода получения сокращённой ссылки
func (m *MockStorage) GetShortLink(originalLink string) (string, error) {
	args := m.Called(originalLink)
	return args.String(0), args.Error(1)
}

// GetOriginalLink — заглушка для метода получения оригинального URL
func (m *MockStorage) GetOriginalLink(shortLink string) (string, error) {
	args := m.Called(shortLink)
	return args.String(0), args.Error(1)
}

// SaveLinks — заглушка для метода сохранения ссылки.
func (m *MockStorage) SaveLinks(shortLink, originalLink string) {
	_ = m.Called(shortLink, originalLink)
}

// TestGetShortLink проверяет бизнес-логику генерации сокращенной ссылки с использованием mock-хранилища
func TestGetShortLink(t *testing.T) {
	mockStore := new(MockStorage)
	service := services.NewShortenerService(mockStore)

	testLink := "https://exemple.com"

	// Проверка случая, когда ссылка уже есть в базе
	mockStore.On("GetShortLink", testLink).Return("abcd123456", nil)

	result, err := service.GetShortLink(testLink)
	require.NoError(t, err)
	assert.Equal(t, "abcd123456", result)

	// Проверка создания новой сокращённой ссылки
	testLink = "https://newsite.com"

	mockCallGetShortLink := mockStore.On("GetShortLink", testLink).Return("", errors.New("not found"))
	mockCallSaveLinks := mockStore.On("SaveLinks", mock.Anything, testLink).Return(nil)

	result, err = service.GetShortLink(testLink)
	require.NoError(t, err)
	assert.Len(t, result, 10, "Сокращённая ссылка должна быть длиной 10 символов")

	mockCallGetShortLink.Unset()
	mockCallSaveLinks.Unset()

	// Проверка случая, когда оригинальная ссылка не корректная
	testLink = "httpss:///bad___link"

	_, err = service.GetShortLink(testLink)
	require.Error(t, err)

	// Проверка случая, когда оригинальная ссылка пустая
	testLink = ""

	_, err = service.GetShortLink(testLink)
	require.Error(t, err)

	mockStore.AssertExpectations(t)
}

// TestGetOriginalLink проверяет получение оригинальной ссылки с использованием mock-хранилища
func TestGetOriginalLink(t *testing.T) {
	mockStore := new(MockStorage)
	service := services.NewShortenerService(mockStore)

	testLink := "https://exemple.com"

	// Проверка получения оригинальной URL
	mockStore.On("GetOriginalLink", "abcd123456").Return(testLink, nil)

	originalURL, err := service.GetOriginalLink("abcd123456")
	require.NoError(t, err)
	assert.Equal(t, testLink, originalURL)

	// Проверка обработки ошибки, если ссылки нет
	mockStore.On("GetOriginalLink", "nonexistent").Return("", errors.New("not found"))

	_, err = service.GetOriginalLink("nonexistent")
	assert.Error(t, err)

	mockStore.AssertExpectations(t)
}
