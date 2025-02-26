package memory_test

import (
	"os"
	"shortened_links_service/internal/storage/memory"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	store            *memory.Memory
	testOriginalLink = "https://example.com"
	testShortLink    = "abcd123456"
)

// TestMain создает хранилище перед запуском тестов
func TestMain(m *testing.M) {
	store = memory.NewMemoryStore()

	store.SaveLinks(testShortLink, testOriginalLink)

	code := m.Run()
	os.Exit(code)
}

// TestGetShortLinkInMemory проверяет извлечение сокращенной ссылки из памяти
func TestGetShortLinkInMemory(t *testing.T) {
	shortLink, err := store.GetShortLink(testOriginalLink)
	require.NoError(t, err)
	require.Equal(t, testShortLink, shortLink)

	_, err = store.GetShortLink("https://nonexistent.com")
	require.Error(t, err)
}

// TestGetOriginalLinkInMemory проверяет извлечение оригинальной ссылки памяти
func TestGetOriginalLinkInMemory(t *testing.T) {
	original_link, err := store.GetOriginalLink(testShortLink)
	assert.NoError(t, err)
	require.Equal(t, testOriginalLink, original_link)

	_, err = store.GetOriginalLink("nonexistent")
	require.Error(t, err)
}
