package domain

import "context"

type Repository interface {
	SaveURL(ctx context.Context, originalURL string) (string, error)
	GetURL(ctx context.Context, shortenedURL string) (string, error)
}
