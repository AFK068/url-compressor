package server

import (
	"context"
	"time"

	"github.com/AFK068/compressor/internal/config"
	"github.com/AFK068/compressor/internal/infrastructure/httpapi/compressorapi"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"

	compressortypes "github.com/AFK068/compressor/internal/api/openapi/compressor/v1"
)

type Compressor struct {
	Config  *config.Config
	Handler *compressorapi.Handler
	Echo    *echo.Echo
	logger  *zap.Logger
}

func NewCompressor(cfg *config.Config, handler *compressorapi.Handler, logger *zap.Logger) *Compressor {
	return &Compressor{
		Config:  cfg,
		Handler: handler,
		Echo:    echo.New(),
		logger:  logger,
	}
}

func (c *Compressor) Start() error {
	compressortypes.RegisterHandlers(c.Echo, c.Handler)

	return c.Echo.Start(":" + c.Config.Shortener.Port)
}

func (c *Compressor) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return c.Echo.Shutdown(ctx)
}

func (c *Compressor) RegisterHooks(lc fx.Lifecycle, log *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			log.Info("Starting compressor server")

			go func() {
				if err := c.Start(); err != nil {
					log.Error("Failed to start compressor server", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(context.Context) error {
			log.Info("Stopping compressor server")

			if err := c.Stop(); err != nil {
				log.Error("Failed to stop compressor server", zap.Error(err))
			}

			return nil
		},
	})
}
