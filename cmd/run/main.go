package main

import (
	"context"

	"github.com/AFK068/compressor/internal/config"
	"github.com/AFK068/compressor/internal/domain"
	"github.com/AFK068/compressor/internal/infrastructure/httpapi/compressorapi"
	"github.com/AFK068/compressor/internal/infrastructure/repository/inmemoryrepo"
	"github.com/AFK068/compressor/internal/infrastructure/repository/postgresdb"
	"github.com/AFK068/compressor/internal/migration"
	"github.com/AFK068/compressor/internal/server"
	"github.com/AFK068/compressor/pkg/logger"
	"github.com/AFK068/compressor/pkg/shortener"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	DevConfigPath = "config/dev.yaml"
)

func NewPostgreDB(cfg *config.Config, shortener domain.Shortener, log *zap.Logger, lc fx.Lifecycle) (domain.Repository, error) {
	if cfg.Storage.Type == domain.InMemoryRepository {
		return inmemoryrepo.New(shortener, cfg.Storage.MaxSize), nil
	}

	err := migration.RunMigration(cfg, log)
	if err != nil {
		log.Fatal("failed to run migration", zap.Error(err))
	}

	dbPool, err := pgxpool.New(context.Background(), cfg.GetPostgresConnectionString())
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			dbPool.Close()
			return nil
		},
	})

	return postgresdb.New(dbPool, shortener, cfg.Storage.MaxSize), nil
}

func main() {
	fx.New(
		fx.Provide(
			// Logger.
			logger.New,

			// Config.
			func() (*config.Config, error) {
				return config.NewConfig(DevConfigPath)
			},

			// Shortener.
			func(cfg *config.Config) (domain.Shortener, error) {
				return shortener.NewShortener(cfg.Shortener.Alphabet, cfg.Shortener.Length)
			},

			// Repository.
			NewPostgreDB,

			// Handler.
			compressorapi.NewHandler,

			// Server.
			server.NewCompressor,
		),
		fx.Invoke(
			func(s *server.Compressor, lc fx.Lifecycle, log *zap.Logger) {
				s.RegisterHooks(lc, log)
			},
		),
	).Run()
}
