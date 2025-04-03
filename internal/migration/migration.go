package migration

import (
	"errors"
	"fmt"

	"github.com/AFK068/compressor/internal/config"
	"github.com/golang-migrate/migrate/v4"
	"go.uber.org/zap"

	_ "github.com/golang-migrate/migrate/v4/database/postgres" //nolint
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5"
)

func RunMigration(cfg *config.Config, log *zap.Logger) error {
	log.Info("Running migration")

	migrator, err := migrate.New(
		fmt.Sprintf("file://%s", cfg.Migration.MigrationsPath),
		cfg.GetPostgresConnectionString(),
	)

	if err != nil {
		return fmt.Errorf("creating migrator: %w", err)
	}

	if err := migrator.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info("no migrations to apply")
			return nil
		}

		return fmt.Errorf("applying migrations: %w", err)
	}

	log.Info("migration completed")

	return nil
}
