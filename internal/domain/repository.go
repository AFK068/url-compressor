package domain

import "context"

type RepositoryType string

const (
	PostgresRepository RepositoryType = "postgres"
	InMemoryRepository RepositoryType = "inmemory"
)

type Repository interface {
	SaveURL(ctx context.Context, originalURL string) (string, error)
	GetURL(ctx context.Context, shortenedURL string) (string, error)
}
