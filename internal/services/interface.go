package services

// ShortenerServiceInterface - интерфейс для сервиса сокращения ссылок
type ShortenerServiceInterface interface {
	GetShortLink(originalLink string) (string, error)
	GetOriginalLink(shortLink string) (string, error)
}
