package compressorapi

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/labstack/echo/v4"

	compressortypes "github.com/AFK068/compressor/internal/api/openapi/compressor/v1"
)

const (
	ErrFailedToGetURL     = "Failed to get URL"
	ErrFailedToPostURL    = "Failed to post URL"
	ErrInvalidRequestBody = "invalid_request_body"
	ErrLinkNotFound       = "link_not_found"

	ErrDescriptionFailedToGetURL     = "Failed to get URL"
	ErrDescriptionFailedToPostURL    = "Failed to post URL"
	ErrDescriptionInvalidRequestBody = "Invalid request body"
	ErrDescriptionLinkNotFound       = "Link not found"
)

func SendSuccessResponse(ctx echo.Context, data any) error {
	return ctx.JSON(http.StatusOK, data)
}

func SendBadRequestResponse(ctx echo.Context, err, description string) error {
	return ctx.JSON(http.StatusBadRequest, compressortypes.ApiErrorResponse{
		Description:      aws.String(description),
		Code:             aws.String("400"),
		ExceptionMessage: aws.String(err),
	})
}

func SendNotFoundResponse(ctx echo.Context, err, description string) error {
	return ctx.JSON(http.StatusNotFound, compressortypes.ApiErrorResponse{
		Description:      aws.String(description),
		Code:             aws.String("404"),
		ExceptionMessage: aws.String(err),
	})
}
