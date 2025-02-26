package storage

// StorageInterface определяет интерфейс для работы с различными хранилищами данных в "in-memory" и "postgres" режимах
type StorageInterface interface {
	SaveLinks(shortLink, originalLink string)
	GetShortLink(shortLink string) (string, error)
	GetOriginalLink(originalLink string) (string, error)
}
