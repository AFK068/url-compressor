package compressorapi_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AFK068/compressor/internal/domain/apperrors"
	"github.com/AFK068/compressor/internal/infrastructure/httpapi/compressorapi"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	compressortypes "github.com/AFK068/compressor/internal/api/openapi/compressor/v1"
	repomock "github.com/AFK068/compressor/internal/domain/mocks"
)

func Test_GetUrl_Success(t *testing.T) {
	repoMock := repomock.NewRepository(t)

	repoMock.On("GetURL", mock.Anything, "shortUrl").Return("http://example.com", nil)
	handler := compressorapi.NewHandler(repoMock, zap.NewNop())

	req := httptest.NewRequest("GET", "/url", http.NoBody)
	rec := httptest.NewRecorder()

	c := echo.New().NewContext(req, rec)

	err := handler.GetUrl(c, compressortypes.GetUrlParams{ShortUrl: "shortUrl"})
	assert.NoError(t, err)

	assert.Equal(t, 200, rec.Code)
	repoMock.AssertExpectations(t)
}

func Test_GetUrl_NotFound_Failure(t *testing.T) {
	repoMock := repomock.NewRepository(t)

	repoMock.On("GetURL", mock.Anything, "shortUrl").Return("", &apperrors.ErrURLNotFound{Message: "url not found"})
	handler := compressorapi.NewHandler(repoMock, zap.NewNop())

	req := httptest.NewRequest("GET", "/url", http.NoBody)
	rec := httptest.NewRecorder()

	c := echo.New().NewContext(req, rec)

	err := handler.GetUrl(c, compressortypes.GetUrlParams{ShortUrl: "shortUrl"})
	assert.NoError(t, err)

	assert.Equal(t, 404, rec.Code)
	repoMock.AssertExpectations(t)
}

func Test_GetUrl_BadRequest_Filure(t *testing.T) {
	repoMock := repomock.NewRepository(t)

	repoMock.On("GetURL", mock.Anything, "shortUrl").Return("", &apperrors.ErrRepositoryIsFull{Message: "repository is full"})
	handler := compressorapi.NewHandler(repoMock, zap.NewNop())

	req := httptest.NewRequest("GET", "/url", http.NoBody)
	rec := httptest.NewRecorder()

	c := echo.New().NewContext(req, rec)

	err := handler.GetUrl(c, compressortypes.GetUrlParams{ShortUrl: "shortUrl"})
	assert.NoError(t, err)

	assert.Equal(t, 400, rec.Code)
	repoMock.AssertExpectations(t)
}

func Test_PostUrl_Success(t *testing.T) {
	repoMock := repomock.NewRepository(t)

	repoMock.On("SaveURL", mock.Anything, "http://example.com").Return("shortUrl", nil)

	handler := compressorapi.NewHandler(repoMock, zap.NewNop())

	body := `{"url": "http://example.com"}`
	req := httptest.NewRequest("POST", "/url", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := handler.PostUrl(c)
	assert.NoError(t, err)

	assert.Equal(t, 200, rec.Code)

	repoMock.AssertExpectations(t)
}

func Test_PostUrl_BadRequest_Failure(t *testing.T) {
	repoMock := repomock.NewRepository(t)

	repoMock.On("SaveURL", mock.Anything, "http://example.com").Return("", &apperrors.ErrRepositoryIsFull{Message: "repository is full"})

	handler := compressorapi.NewHandler(repoMock, zap.NewNop())

	body := `{"url": "http://example.com"}`
	req := httptest.NewRequest("POST", "/url", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := handler.PostUrl(c)
	assert.NoError(t, err)

	assert.Equal(t, 400, rec.Code)

	repoMock.AssertExpectations(t)
}

func Test_PostUrl_InvalidRequestBody_Failure(t *testing.T) {
	repoMock := repomock.NewRepository(t)

	handler := compressorapi.NewHandler(repoMock, zap.NewNop())

	body := `{"asd": "http://example.com"}`
	req := httptest.NewRequest("POST", "/url", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := handler.PostUrl(c)
	assert.NoError(t, err)

	assert.Equal(t, 400, rec.Code)

	repoMock.AssertExpectations(t)
}
