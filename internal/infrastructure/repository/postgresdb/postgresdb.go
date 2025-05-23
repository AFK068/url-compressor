package postgresdb

import (
	"context"
	"errors"

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

func (r *PostgresRepository) SaveURL(ctx context.Context, originalURL string) (shortenedURL string, err error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return "", err
	}

	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback(ctx))
		}
	}()

	existingShortURL, err := r.getExistingShortURL(ctx, tx, originalURL)
	if err == nil {
		err = tx.Commit(ctx)
		if err != nil {
			return "", err
		}

		return existingShortURL, nil
	}

	if err != pgx.ErrNoRows {
		return "", err
	}

	id, err := r.insertURL(ctx, tx, originalURL)
	if err != nil {
		return "", err
	}

	if id >= r.maxSize {
		return "", &apperrors.ErrRepositoryIsFull{Message: "repository is full"}
	}

	shortenedURL, err = r.shortener.Encode(id)
	if err != nil {
		return "", err
	}

	err = r.updateShortURL(ctx, tx, id, shortenedURL)
	if err != nil {
		return "", err
	}

	err = tx.Commit(ctx)
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

func (r *PostgresRepository) getExistingShortURL(ctx context.Context, tx pgx.Tx, originalURL string) (string, error) {
	query, args, err := squirrel.Select("short_url").
		From("urls").
		Where(squirrel.Eq{"url": originalURL}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return "", err
	}

	var shortURL string
	err = tx.QueryRow(ctx, query, args...).Scan(&shortURL)

	return shortURL, err
}

func (r *PostgresRepository) insertURL(ctx context.Context, tx pgx.Tx, originalURL string) (uint64, error) {
	query, args, err := squirrel.Insert("urls").
		Columns("url").
		Values(originalURL).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return 0, err
	}

	var id uint64
	err = tx.QueryRow(ctx, query, args...).Scan(&id)

	return id, err
}

func (r *PostgresRepository) updateShortURL(ctx context.Context, tx pgx.Tx, id uint64, shortURL string) error {
	query, args, err := squirrel.Update("urls").
		Set("short_url", shortURL).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, query, args...)

	return err
}
