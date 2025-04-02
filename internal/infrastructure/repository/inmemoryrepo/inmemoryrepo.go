package inmemoryrepo

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/AFK068/compressor/internal/domain"
	"github.com/AFK068/compressor/internal/domain/apperrors"
)

type InMemoryRepository struct {
	urls      []string
	shortener domain.Shortener
	counter   atomic.Uint64
	mu        sync.Mutex
}

func New(shortener domain.Shortener, maxSize uint64) *InMemoryRepository {
	return &InMemoryRepository{
		urls:      make([]string, maxSize),
		shortener: shortener,
	}
}

func (r *InMemoryRepository) SaveURL(_ context.Context, originalURL string) (string, error) {
	id := r.counter.Add(1) - 1

	if id >= uint64(len(r.urls)) {
		return "", &apperrors.ErrRepositoryIsFull{Message: "repository is full"}
	}

	shortenedURL, err := r.shortener.Encode(id)
	if err != nil {
		return "", err
	}

	r.mu.Lock()
	r.urls[id] = originalURL
	r.mu.Unlock()

	return shortenedURL, nil
}

func (r *InMemoryRepository) GetURL(_ context.Context, shortenedURL string) (string, error) {
	id, err := r.shortener.Decode(shortenedURL)
	if err != nil {
		return "", err
	}

	if id >= uint64(len(r.urls)) {
		return "", &apperrors.ErrURLNotFound{Message: "url not found"}
	}

	r.mu.Lock()
	originalURL := r.urls[id]
	r.mu.Unlock()

	if originalURL == "" {
		return "", &apperrors.ErrURLNotFound{Message: "url not found"}
	}

	return originalURL, nil
}
