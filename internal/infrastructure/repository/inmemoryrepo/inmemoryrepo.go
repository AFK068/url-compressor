package inmemoryrepo

import (
	"context"
	"sync"

	"github.com/AFK068/compressor/internal/domain"
	"github.com/AFK068/compressor/internal/domain/apperrors"
	rbt "github.com/emirpasic/gods/trees/redblacktree"
)

type InMemoryRepository struct {
	urls      []string
	urlTree   *rbt.Tree
	shortener domain.Shortener
	counter   uint64
	mu        sync.Mutex
	maxSize   uint64
}

func New(shortener domain.Shortener, maxSize uint64) *InMemoryRepository {
	return &InMemoryRepository{
		urls:      make([]string, maxSize),
		shortener: shortener,
		urlTree:   rbt.NewWithStringComparator(),
		maxSize:   maxSize,
	}
}

func (r *InMemoryRepository) SaveURL(_ context.Context, originalURL string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if val, ok := r.urlTree.Get(originalURL); ok {
		return val.(string), nil
	}

	if r.counter >= r.maxSize {
		return "", &apperrors.ErrRepositoryIsFull{Message: "repository is full"}
	}

	shortenedURL, err := r.shortener.Encode(r.counter)
	if err != nil {
		return "", err
	}

	r.urls[r.counter] = originalURL
	r.urlTree.Put(originalURL, shortenedURL)

	r.counter++

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
