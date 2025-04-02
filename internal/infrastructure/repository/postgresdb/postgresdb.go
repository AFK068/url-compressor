package postgresdb

import (
	"context"

	"github.com/AFK068/compressor/internal/domain"
	"github.com/AFK068/compressor/internal/domain/apperrors"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool      *pgxpool.Pool
	shortener domain.Shortener
	maxSize   uint64
}

func New(pool *pgxpool.Pool, shortener domain.Shortener, maxSize uint64) *PostgresRepository {
	return &PostgresRepository{
		pool:      pool,
		shortener: shortener,
		maxSize:   maxSize,
	}
}

func (r *PostgresRepository) SaveURL(ctx context.Context, originalURL string) (string, error) {
	query, args, err := squirrel.Insert("urls").
		Columns("url").
		Values(originalURL).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return "", err
	}

	var id uint64

	err = r.pool.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", &apperrors.ErrRepositoryIsFull{Message: "repository is full"}
		}

		return "", err
	}

	if id >= r.maxSize {
		return "", &apperrors.ErrRepositoryIsFull{Message: "repository is full"}
	}

	shortenedURL, err := r.shortener.Encode(id)
	if err != nil {
		return "", err
	}

	return shortenedURL, nil
}

func (r *PostgresRepository) GetURL(ctx context.Context, shortenedURL string) (string, error) {
	id, err := r.shortener.Decode(shortenedURL)
	if err != nil {
		return "", err
	}

	if id >= r.maxSize {
		return "", &apperrors.ErrURLNotFound{Message: "url not found"}
	}

	query, args, err := squirrel.Select("url").
		From("urls").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return "", err
	}

	var originalURL string

	err = r.pool.QueryRow(ctx, query, args...).Scan(&originalURL)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", &apperrors.ErrURLNotFound{Message: "url not found"}
		}

		return "", err
	}

	if originalURL == "" {
		return "", &apperrors.ErrURLNotFound{Message: "url not found"}
	}

	return originalURL, nil
}
