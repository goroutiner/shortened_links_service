package database

import (
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Database struct {
	db *sqlx.DB
}

// NewDatabaseConection возвращает структуру соединения с БД
func NewDatabaseConection(db *sqlx.DB) *Database {
	return &Database{db: db}
}

// NewDatabaseStore создает таблицу в БД для хранения посылок
func NewDatabaseStore(psqlUrl string) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", psqlUrl)
	if err != nil {
		return nil, err
	}

	db.MustExec(`
		CREATE TABLE IF NOT EXISTS links (
            id SERIAL PRIMARY KEY,
            short_link VARCHAR(10) UNIQUE NOT NULL,
            original_link TEXT UNIQUE NOT NULL
        );
		
		CREATE INDEX IF NOT EXISTS link_indx ON links (short_link, original_link);
		`)

	log.Println("Creating a database ...")

	return db, err
}

// SaveLinks сохраняет сокращенную и оригинальную ссылку в таблице links
func (s *Database) SaveLinks(shortLink, originalLink string) {
	_, _ = s.db.Exec("INSERT INTO links (short_link, original_link) VALUES ($1, $2) ON CONFLICT (original_link) DO NOTHING", shortLink, originalLink)
}

// GetOriginalLink возвращает сокращенную ссылку, полученную по ее оригиналу записанному в таблице links
func (s *Database) GetOriginalLink(shortLink string) (string, error) {
	var originalLink string
	err := s.db.Get(&originalLink, "SELECT original_link FROM links WHERE short_link = $1", shortLink)
	return originalLink, err
}

// GetShortLink возвращает оригинальную ссылку, полученную по ее сокращенной форме записанной в таблице links
func (s *Database) GetShortLink(originalLink string) (string, error) {
	var shortLink string
	err := s.db.Get(&shortLink, "SELECT short_link FROM links WHERE original_link = $1", originalLink)
	return shortLink, err
}
