package postgresdb_test

import (
	"context"
	"testing"

	"github.com/AFK068/compressor/internal/config"
	"github.com/AFK068/compressor/internal/domain/apperrors"
	"github.com/AFK068/compressor/internal/infrastructure/repository/postgresdb"
	"github.com/AFK068/compressor/internal/testcontainer"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"

	shortenermock "github.com/AFK068/compressor/internal/domain/mocks"
)

const (
	TestConfigPath = "../../../../config/test.yaml"
)

func setupDB(t *testing.T) (*pgxpool.Pool, context.Context) {
	ctx := context.Background()

	config, err := config.NewConfig(TestConfigPath)
	assert.NoError(t, err)

	testContainer, err := testcontainer.NewPostgresTestcontainerContainer(ctx, config)
	assert.NoError(t, err)

	dbPool, cleanup, err := testContainer.SetupTestPostgresContainer(ctx)
	assert.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, cleanup())
	})

	return dbPool, ctx
}

func Test_SaveURL_Success(t *testing.T) {
	dbPool, ctx := setupDB(t)

	shortenerMock := shortenermock.NewShortener(t)
	shortenerMock.On("Encode", uint64(0)).Return("shortURL", nil).Once()

	repo := postgresdb.New(dbPool, shortenerMock, 10)

	short, err := repo.SaveURL(ctx, "originURL")
	assert.NoError(t, err)
	assert.Equal(t, "shortURL", short)

	query := `SELECT url FROM urls WHERE id = 0`

	var originalURL string
	err = dbPool.QueryRow(ctx, query).Scan(&originalURL)

	assert.NoError(t, err)
	assert.Equal(t, "originURL", originalURL)
	shortenerMock.AssertExpectations(t)
}

func Test_URLIsAlreadyExists_Success(t *testing.T) {
	dbPool, ctx := setupDB(t)

	shortenerMock := shortenermock.NewShortener(t)
	shortenerMock.On("Encode", uint64(0)).Return("shortURL", nil).Once()

	repo := postgresdb.New(dbPool, shortenerMock, 10)
	short, err := repo.SaveURL(ctx, "originURL")
	assert.NoError(t, err)
	assert.Equal(t, "shortURL", short)

	short2, err := repo.SaveURL(ctx, "originURL")
	assert.NoError(t, err)
	assert.Equal(t, "shortURL", short2)

	query := `SELECT id FROM urls ORDER BY id DESC LIMIT 1`

	var lastID int
	err = dbPool.QueryRow(ctx, query).Scan(&lastID)

	assert.NoError(t, err)
	assert.Equal(t, 0, lastID)
}

func Test_SaveURL_RepoIsFull_Failure(t *testing.T) {
	dbPool, ctx := setupDB(t)

	shortenerMock := shortenermock.NewShortener(t)
	shortenerMock.On("Encode", uint64(0)).Return("123", nil).Once()
	shortenerMock.On("Encode", uint64(1)).Return("321", nil).Once()

	repo := postgresdb.New(dbPool, shortenerMock, 2)

	_, err := repo.SaveURL(ctx, "originURL1")
	assert.NoError(t, err)

	_, err = repo.SaveURL(ctx, "originURL2")
	assert.NoError(t, err)

	_, err = repo.SaveURL(ctx, "originURL3")
	assert.Error(t, err)

	shortenerMock.AssertExpectations(t)
}

func Test_GetURL_Success(t *testing.T) {
	dbPool, ctx := setupDB(t)

	shortenerMock := shortenermock.NewShortener(t)
	shortenerMock.On("Decode", "shortURL").Return(uint64(0), nil).Once()

	repo := postgresdb.New(dbPool, shortenerMock, 10)

	shortenerMock.On("Encode", uint64(0)).Return("shortURL", nil).Once()

	_, err := repo.SaveURL(ctx, "originURL")
	assert.NoError(t, err)

	originalURL, err := repo.GetURL(ctx, "shortURL")
	assert.NoError(t, err)
	assert.Equal(t, "originURL", originalURL)
}

func Test_GetURl_LinkNotFound_Failure(t *testing.T) {
	dbPool, ctx := setupDB(t)

	shortenerMock := shortenermock.NewShortener(t)
	shortenerMock.On("Decode", "shortURL").Return(uint64(0), nil).Once()

	repo := postgresdb.New(dbPool, shortenerMock, 10)

	_, err := repo.GetURL(ctx, "shortURL")
	assert.Error(t, err)
	assert.IsType(t, &apperrors.ErrURLNotFound{}, err)

	shortenerMock.AssertExpectations(t)
}
