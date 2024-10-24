package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LoggerLevel string `json:"logger_level" validate:"required,oneof=DEBUG INFO WARN ERROR" env:"LOGGER_LEVEL"`
	Server      `json:"server"`
	Dvault      `json:"dvault"`
}

type Dvault struct {
	StorageType   string      `json:"storage_type" validate:"required,oneof=redis postgres fs" env:"STORAGE_TYPE"`
	StorageConfig interface{} `json:"storage_config"`

	EncryptionMethod string `json:"encryption_method" validate:"required" env:"ENCRYPTION_METHOD"`
}

type FSStorageConfig struct {
	MountPath string `json:"mount_path" validate:"required" env:"MOUNT_PATH"`
}

type RedisStorageConfig struct {
	Connection string `json:"mount_path" validate:"required" env:"DB_REDIS"`
	SSLEnabled bool   `json:"ssl_enabled" env:"REDIS_SSL_ENABLED"`
	CertPath   string `json:"cert_path" env:"CERT_REDIS_PATH_DVAULT"`
}

type PostgresqlStorageConfig struct {
	Connection string `json:"mount_path" validate:"required" env:"DB"`
	SSLEnabled bool   `json:"ssl_enabled" env:"DB_SSL_ENABLED"`
	CertPath   string `json:"cert_path" env:"CERT_DB_PATH_DVAULT"`
}

type Server struct {
	Addr       string `json:"addr" validate:"required,hostname_port" env:"PORT"`
	SSLEnabled bool   `json:"ssl_enabled" env:"SERVER_SSL_ENABLED"`
	CertPath   string `json:"cert_path" env:"SERVER_CERT_PATH"`
	KeyPath    string `json:"key_path" env:"SERVER_KEY_PATH"`
}

func Default() (Config, error) {
	return Config{
		LoggerLevel: "DEBUG",
		Server:      Server{Addr: ":8080"},
		Dvault: Dvault{
			StorageType: "postgres",
			StorageConfig: PostgresqlStorageConfig{
				Connection: "postgres://postgres:password@localhost/dvault?sslmode=disable",
				CertPath:   "",
			},
			/*			StorageConfig: RedisStorageConfig{
						Connection: "redis:@localhost:6379/db",
						CertPath:   "",
					},*/
			EncryptionMethod: "chacha20-poly1305",
		},
	}, nil
}

func ReadEnv() (Config, error) {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return Config{}, err
	}

	if cfg.LoggerLevel == "" {
		cfg.LoggerLevel = "INFO"
	}
	if cfg.Dvault.EncryptionMethod == "" {
		cfg.Dvault.EncryptionMethod = "aes"
	}

	switch cfg.StorageType {
	case "postgres":
		var pgcfg PostgresqlStorageConfig

		if err := cleanenv.ReadEnv(&pgcfg); err != nil {
			return Config{}, err
		}

		cfg.StorageConfig = pgcfg
	case "redis":
		var redisCfg RedisStorageConfig

		if err := cleanenv.ReadEnv(&redisCfg); err != nil {
			return Config{}, err
		}

		cfg.StorageConfig = redisCfg
	case "fs":
		var fs FSStorageConfig

		if err := cleanenv.ReadEnv(&fs); err != nil {
			return Config{}, err
		}

		cfg.StorageConfig = fs
	default:
		return Config{}, fmt.Errorf("invalid storage type %s", cfg.StorageConfig)
	}

	if err := validator.New().Struct(cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
