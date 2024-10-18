package config

import (
	"encoding/json"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LoggerLevel string `json:"logger_level" validate:"required,oneof=DEBUG INFO WARN ERROR" env:"LOGGER_LEVEL"`
	Server      `json:"server"`
	Postgres    `json:"postgres"`
	Dvault      `json:"dvault"`
}

type Dvault struct {
	MountPath        string `json:"mount_path" validate:"required" env:"MOUNT_PATH"`
	EncryptionMethod string `json:"encryption_method" validate:"required" env:"ENCRYPTION_METHOD"`
}

type Postgres struct {
	Addr string `json:"addr" validate:"required,url" env:"DB"`
}

type Server struct {
	Addr string `json:"addr" validate:"required,hostname_port" env:"PORT"`
}

func Default() (Config, error) {
	return Config{
		LoggerLevel: "DEBUG",
		Server:      Server{Addr: ":8080"},
		Postgres:    Postgres{Addr: "postgres://postgres:password@localhost:5432/vault"},
		Dvault: Dvault{
			MountPath:        "./data",
			EncryptionMethod: "chacha20-poly1305",
		},
	}, nil
}

func ReadFile(fileName string) (Config, error) {
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{}
	if err = json.Unmarshal(bytes, &cfg); err != nil {
		return Config{}, err
	}

	if cfg.LoggerLevel == "" {
		cfg.LoggerLevel = "INFO"
	}
	if cfg.Dvault.EncryptionMethod == "" {
		cfg.Dvault.EncryptionMethod = "aes"
	}

	if err := validator.New().Struct(cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
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

	if err := validator.New().Struct(cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
