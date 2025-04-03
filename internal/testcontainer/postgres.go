package testcontainer

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/AFK068/compressor/internal/config"
	"github.com/AFK068/compressor/internal/migration"
	"github.com/AFK068/compressor/pkg/logger"
	"github.com/docker/go-connections/nat"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
)

type CleanFunc func() error

type PostgresTestcontainer struct {
	testcontainers.Container
	Config *config.Config
}

const (
	DefaultTestContainerImage = "postgres:13"
)

func NewPostgresTestcontainerContainer(ctx context.Context, cfg *config.Config) (*PostgresTestcontainer, error) {
	portWithTCP := fmt.Sprintf("%s/tcp", cfg.Storage.Port)

	req := testcontainers.ContainerRequest{
		Image:        DefaultTestContainerImage,
		ExposedPorts: []string{portWithTCP},
		Env: map[string]string{
			"POSTGRES_USER":     cfg.Storage.User,
			"POSTGRES_PASSWORD": cfg.Storage.Password,
			"POSTGRES_DB":       cfg.Storage.DatabaseName,
		},
		WaitingFor: wait.ForListeningPort(nat.Port(portWithTCP)),
	}

	postgres, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("error starting Postgres container: %w", err)
	}

	return &PostgresTestcontainer{
		Container: postgres,
		Config:    cfg,
	}, nil
}

func (p *PostgresTestcontainer) SetupTestPostgresContainer(ctx context.Context) (*pgxpool.Pool, CleanFunc, error) {
	contextWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	port, err := p.Start(contextWithTimeout)
	if err != nil {
		return nil, nil, err
	}

	p.SetMappedPort(port)

	log := logger.New()

	if err := p.Migrate(log); err != nil {
		return nil, nil, err
	}

	dbPool, err := pgxpool.New(ctx, p.Config.GetPostgresConnectionString())
	if err != nil {
		return nil, nil, err
	}

	clean := func() error {
		dbPool.Close()

		if err := p.Stop(); err != nil {
			return err
		}

		return nil
	}

	return dbPool, clean, nil
}

func (p *PostgresTestcontainer) Start(ctx context.Context) (int, error) {
	port, err := p.Container.MappedPort(ctx, nat.Port(p.Config.Storage.Port))
	if err != nil {
		return 0, err
	}

	return port.Int(), nil
}

func (p *PostgresTestcontainer) Stop() error {
	contextWithTimeout, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return p.Container.Terminate(contextWithTimeout)
}

func (p *PostgresTestcontainer) Migrate(log *zap.Logger) error {
	return migration.RunMigration(p.Config, log)
}

func (p *PostgresTestcontainer) SetMappedPort(mappedPort int) {
	p.Config.Storage.Port = strconv.Itoa(mappedPort)
}
