package services

import (
	"crypto/rand"
	"errors"
	"math/big"
	"net/url"
	"shortened_links_service/internal/config"
	"shortened_links_service/internal/storage"
	"strings"
)

// generateShortLink генерирует сокращенную ссылку, используя криптографически безопасную библиотеку crypto/rand
func generateShortLink() string {
	shortLink := strings.Builder{}
	shortLink.Grow(config.LenShortLink)

	for i := 0; i < config.LenShortLink; i++ {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(config.CharList))))
		shortLink.WriteByte(config.CharList[randomIndex.Int64()])
	}

	return shortLink.String()
}

type ShortenerService struct {
	storage storage.StorageInterface
}

// isValidURL проверка на соответствие ссылки ее минимальным требованиям
func isValidURL(link string) bool {
	parsedURL, err := url.Parse(link)
	return err == nil && parsedURL.Scheme != "" && parsedURL.Host != ""
}

// NewShortenerService возвращает сервис по созданию сокращенных ссылок
func NewShortenerService(s storage.StorageInterface) *ShortenerService {
	return &ShortenerService{storage: s}
}

// GetShortLink возвращает и сохраняет созданную сокращёную ссылку
func (s *ShortenerService) GetShortLink(originalLink string) (string, error) {
	var (
		shortLink string
		err       error
	)

	if !isValidURL(originalLink) {
		return "", errors.New("invalid URL format")
	}

	// проверка на наличие уже существующих ссылок
	shortLink, err = s.storage.GetShortLink(originalLink)
	if err == nil && shortLink != "" {
		return shortLink, nil
	}

	shortLink = generateShortLink()
	err = s.storage.SaveLinks(shortLink, originalLink)
	if err != nil {
		return "", err
	}

	return shortLink, nil

}

// GetOriginalLink возвращает оригинальную ссылку
func (s *ShortenerService) GetOriginalLink(shortLink string) (string, error) {
	originalLink, err := s.storage.GetOriginalLink(shortLink)
	if err != nil {
		return "", err
	}
	return originalLink, nil
}
