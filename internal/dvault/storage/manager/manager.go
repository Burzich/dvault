package manager

import (
	"errors"
	"fmt"
	"log/slog"

	config2 "github.com/Burzich/dvault/internal/config"
	"github.com/Burzich/dvault/internal/dvault/storage"
	fs "github.com/Burzich/dvault/internal/dvault/storage/disc"
	"github.com/Burzich/dvault/internal/dvault/storage/postgres"
	"github.com/Burzich/dvault/internal/dvault/storage/redis"
)

func CreateStorage(storageType string, cfg interface{}, logger *slog.Logger) (storage.Storage, error) {
	switch storageType {
	case "postgres":
		postgresCfg, ok := cfg.(config2.PostgresqlStorageConfig)
		if !ok {
			return nil, errors.New("invalid fs config")
		}
		return postgres.NewPostgresStorage(postgresCfg, logger)
	case "fs":
		fsCfg, ok := cfg.(config2.FSStorageConfig)
		if !ok {
			return nil, errors.New("invalid fs config")
		}
		return fs.NewFSStorage(fsCfg), nil
	case "redis":
		redisCfg, ok := cfg.(config2.RedisStorageConfig)
		if !ok {
			return nil, errors.New("invalid fs config")
		}

		return redis.NewRedisStorage(redisCfg, logger)
	default:
		return nil, fmt.Errorf("invalid storage type %s", storageType)
	}
}
