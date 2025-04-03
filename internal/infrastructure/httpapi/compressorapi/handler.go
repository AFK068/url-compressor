package compressorapi

import (
	"errors"

	"github.com/AFK068/compressor/internal/domain"
	"github.com/AFK068/compressor/internal/domain/apperrors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	compressortypes "github.com/AFK068/compressor/internal/api/openapi/compressor/v1"
)

type Handler struct {
	repository domain.Repository
	logger     *zap.Logger
}

func NewHandler(repository domain.Repository, logger *zap.Logger) *Handler {
	return &Handler{
		repository: repository,
		logger:     logger,
	}
}

func (h *Handler) GetUrl(ctx echo.Context, params compressortypes.GetUrlParams) error { //nolint
	h.logger.Info("Get URL request received", zap.String("shortUrl", params.ShortUrl))

	originalURL, err := h.repository.GetURL(ctx.Request().Context(), params.ShortUrl)

	var errURLNotFound *apperrors.ErrURLNotFound
	if errors.As(err, &errURLNotFound) {
		h.logger.Error("URL not found", zap.String("shortUrl", params.ShortUrl))
		return SendNotFoundResponse(ctx, ErrLinkNotFound, ErrDescriptionLinkNotFound)
	}

	if err != nil {
		h.logger.Error("Failed to get URL", zap.Error(err))
		return SendBadRequestResponse(ctx, ErrFailedToGetURL, ErrDescriptionFailedToGetURL)
	}

	h.logger.Info("Successfully retrieved URL", zap.String("url", originalURL))

	return SendSuccessResponse(ctx, compressortypes.UrlResponse{
		Url: aws.String(originalURL),
	})
}

func (h *Handler) PostUrl(ctx echo.Context) error { //nolint
	h.logger.Info("Post URL request received")

	var request compressortypes.AddUrlRequest
	if err := ctx.Bind(&request); err != nil {
		h.logger.Error("Failed to bind request", zap.Error(err))
		return SendBadRequestResponse(ctx, ErrInvalidRequestBody, ErrDescriptionInvalidRequestBody)
	}

	if request.Url == nil || *request.Url == "" {
		h.logger.Error("URL is empty")
		return SendBadRequestResponse(ctx, ErrInvalidRequestBody, ErrDescriptionInvalidRequestBody)
	}

	short, err := h.repository.SaveURL(ctx.Request().Context(), *request.Url)
	if err != nil {
		h.logger.Error("Failed to save URL", zap.Error(err))
		return SendBadRequestResponse(ctx, ErrFailedToPostURL, ErrDescriptionFailedToPostURL)
	}

	h.logger.Info("Successfully saved URL", zap.String("shortUrl", short))

	return SendSuccessResponse(ctx, compressortypes.UrlResponse{
		Url: aws.String(short),
	})
}
