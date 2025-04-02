package inmemoryrepo_test

import (
	"context"
	"testing"

	"github.com/AFK068/compressor/internal/domain/apperrors"
	"github.com/AFK068/compressor/internal/infrastructure/repository/inmemoryrepo"
	"github.com/stretchr/testify/assert"

	shortenermock "github.com/AFK068/compressor/internal/domain/mocks"
)

func Test_SaveURL_Success(t *testing.T) {
	shortenerMock := shortenermock.NewShortener(t)
	repo := inmemoryrepo.New(shortenerMock, 10)

	shortenerMock.On("Encode", uint64(0)).Return("shortenedURL", nil).Once()
	shortenerMock.On("Encode", uint64(1)).Return("shortenedURL2", nil).Once()

	shortURL, err := repo.SaveURL(context.Background(), "http://example.com")
	assert.NoError(t, err)
	assert.Equal(t, "shortenedURL", shortURL)

	shortURL2, err := repo.SaveURL(context.Background(), "http://example.com/2")
	assert.NoError(t, err)
	assert.Equal(t, "shortenedURL2", shortURL2)

	shortenerMock.AssertExpectations(t)
}

func Test_SaveURL_RepoIsFull_Failure(t *testing.T) {
	shortenerMock := shortenermock.NewShortener(t)
	repo := inmemoryrepo.New(shortenerMock, 1)

	shortenerMock.On("Encode", uint64(0)).Return("shortenedURL", nil).Once()

	_, err := repo.SaveURL(context.Background(), "http://example.com")
	assert.NoError(t, err)

	_, err = repo.SaveURL(context.Background(), "http://example.com/2")
	assert.Error(t, err)
	assert.IsType(t, &apperrors.ErrRepositoryIsFull{}, err)

	shortenerMock.AssertExpectations(t)
}

func Test_SaveURL_EncodeError_Failure(t *testing.T) {
	shortenerMock := shortenermock.NewShortener(t)
	repo := inmemoryrepo.New(shortenerMock, 10)

	shortenerMock.On("Encode", uint64(0)).Return("", assert.AnError).Once()

	_, err := repo.SaveURL(context.Background(), "http://example.com")
	assert.Error(t, err)

	shortenerMock.AssertExpectations(t)
}

func Test_GetURL_Success(t *testing.T) {
	shortenerMock := shortenermock.NewShortener(t)
	repo := inmemoryrepo.New(shortenerMock, 10)

	shortenerMock.On("Decode", "shortenedURL").Return(uint64(0), nil).Once()
	shortenerMock.On("Encode", uint64(0)).Return("shortenedURL", nil).Once()

	short, err := repo.SaveURL(context.Background(), "http://example.com")
	assert.NoError(t, err)

	originalURL, err := repo.GetURL(context.Background(), short)
	assert.NoError(t, err)

	assert.Equal(t, "http://example.com", originalURL)

	shortenerMock.AssertExpectations(t)
}

func Test_GetURL_DecodeError_Failure(t *testing.T) {
	shortenerMock := shortenermock.NewShortener(t)
	repo := inmemoryrepo.New(shortenerMock, 10)

	shortenerMock.On("Decode", "shortenedURL").Return(uint64(0), assert.AnError).Once()

	_, err := repo.GetURL(context.Background(), "shortenedURL")
	assert.Error(t, err)

	shortenerMock.AssertExpectations(t)
}

func Test_GetURL_URLNotFound_Failure(t *testing.T) {
	shortenerMock := shortenermock.NewShortener(t)
	repo := inmemoryrepo.New(shortenerMock, 10)

	shortenerMock.On("Decode", "shortenedURL").Return(uint64(0), nil).Once()

	originalURL, err := repo.GetURL(context.Background(), "shortenedURL")
	assert.Error(t, err)
	assert.IsType(t, &apperrors.ErrURLNotFound{}, err)
	assert.Empty(t, originalURL)

	shortenerMock.AssertExpectations(t)
}
