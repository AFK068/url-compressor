package config

import (
	"fmt"

	"github.com/AFK068/compressor/internal/domain"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Storage   Storage   `yaml:"storage" env-required:"true"`
	Migration Migration `yaml:"migrations" env-required:"true"`
	Shortener Shortener `yaml:"shortener" env-required:"true"`
}

type Storage struct {
	Type         domain.RepositoryType `yaml:"type" env:"STORAGE_TYPE" env-default:"postgres"`
	MaxSize      uint64                `yaml:"max_size" env:"MAX_SIZE" env-default:"100"`
	Host         string                `yaml:"host" env:"POSTGRES_HOST" env-required:"true"`
	Port         string                `yaml:"port" env:"POSTGRES_PORT" env-required:"true"`
	DatabaseName string                `yaml:"database_name" env:"POSTGRES_DATABASE_NAME" env-required:"true"`
	User         string                `yaml:"user" env:"POSTGRES_USER" env-required:"true"`
	Password     string                `yaml:"password" env:"POSTGRES_PASSWORD" env-required:"true"`
}

type Migration struct {
	MigrationsPath string `yaml:"migrations_path" env:"MIGRATIONS_PATH" env-required:"true"`
}

type Shortener struct {
	Port     string `yaml:"port" env:"SHORTENER_PORT" env-default:"8080"`
	Alphabet string `yaml:"alphabet" env:"ALPHABET" env-required:"true"`
	Length   uint64 `yaml:"length" env:"LENGTH" env-required:"true"`
}

func NewConfig(filePath string) (*Config, error) {
	config := &Config{}

	if err := cleanenv.ReadConfig(filePath, config); err != nil {
		return nil, err
	}

	switch config.Storage.Type {
	case domain.InMemoryRepository, domain.PostgresRepository:
	default:
		config.Storage.Type = domain.PostgresRepository
	}

	return config, nil
}

func (cfg *Config) GetPostgresConnectionString() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Storage.User,
		cfg.Storage.Password,
		cfg.Storage.Host,
		cfg.Storage.Port,
		cfg.Storage.DatabaseName,
	)
}
