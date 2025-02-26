package memory

import (
	"errors"
	"sync"
)

type Memory struct {
	shortToOriginal map[string]string // shortToOriginal словарь для связи short_link -> original_link
	originalToShort map[string]string // originalToShort словарь для связи original_link -> short_link
	mu              sync.RWMutex
}

// NewMemoryStore возвращает структуру для хранения ссылок во внутренней памяти
func NewMemoryStore() *Memory {
	return &Memory{
		shortToOriginal: make(map[string]string),
		originalToShort: make(map[string]string),
	}
}

// SaveLinks сохраняет сокращенную и оригинальную ссылку в Memory
func (m *Memory) SaveLinks(shortLink, originalLink string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.shortToOriginal[shortLink] = originalLink
	m.originalToShort[originalLink] = shortLink
}

// GetShortLink возвращает сокращенную ссылку, полученную по ее оригиналу записанному в Memory
func (m *Memory) GetShortLink(originalLink string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	shortLink, has := m.originalToShort[originalLink]
	if !has {
		return "", errors.New("short link is not found")
	}
	return shortLink, nil
}

// GetOriginalLink возвращает оригинальную ссылку, полученную по ее сокращенной форме записанной в Memory
func (m *Memory) GetOriginalLink(shortLink string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	originalLink, has := m.shortToOriginal[shortLink]
	if !has {
		return "", errors.New("original link is not found")
	}
	return originalLink, nil
}
