package database_test

import (
	"log"
	"os"
	"shortened_links_service/internal/storage/database"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	psqlUrl = "postgres://user:password@localhost:5432/test_db?sslmode=disable"
	testDb  *sqlx.DB
	store   *database.Database
)

// TestMain производит соединения с БД и создает таблицу для хранения ссылок
func TestMain(m *testing.M) {
	var err error
	testDb, err = sqlx.Open("pgx", psqlUrl)
	if err != nil {
		log.Fatalf("Ошибка подключения к тестовой БД: %v", err)
	}

	testDb.MustExec(`
		CREATE TABLE IF NOT EXISTS links (
            id SERIAL PRIMARY KEY,
            short_link VARCHAR(10) UNIQUE NOT NULL,
            original_link TEXT UNIQUE NOT NULL
        );
		`)

	store = database.NewDatabaseConection(testDb)

	code := m.Run()

	deleteTable()

	testDb.Close()

	os.Exit(code)
}

// TestSaveLinksInDb проверяет сохранение ссылок в БД
func TestSaveLinksInDb(t *testing.T) {
	originalLink := "https://example.com"
	shortLink := "abcd123456"

	store.SaveLinks(shortLink, originalLink)

	var count int
	err := testDb.Get(&count, "SELECT COUNT(*) FROM links WHERE short_link=$1 AND original_link=$2", shortLink, originalLink)
	require.NoError(t, err)
	assert.Equal(t, 1, count, "Ссылка не сохранились в БД")
}

// TestGetOriginalLinkInDb проверяет получение оригинальной ссылки из БД
func TestGetOriginalLinkInDb(t *testing.T) {
	originalLink := "https://another.com"
	shortLink := "xyz987abc"

	_, err := testDb.Exec("INSERT INTO links (short_link, original_link) VALUES ($1, $2)", shortLink, originalLink)
	require.NoError(t, err)

	result, err := store.GetOriginalLink(shortLink)
	require.NoError(t, err)
	assert.Equal(t, originalLink, result)
}

// TestGetShortLinkInDb проверяет получение сокращённой ссылки из БД
func TestGetShortLinkInDb(t *testing.T) {
	originalLink := "https://test.com"
	shortLink := "lmn456pqr"

	_, err := testDb.Exec("INSERT INTO links (short_link, original_link) VALUES ($1, $2)", shortLink, originalLink)
	require.NoError(t, err)

	result, err := store.GetShortLink(originalLink)
	require.NoError(t, err)
	assert.Equal(t, shortLink, result)
}

// deleteTable удаляет таблицу после каждого теста
func deleteTable() {
	testDb.MustExec("DROP TABLE IF EXISTS links;")
}
