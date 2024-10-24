package redis

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/Burzich/dvault/internal/config"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	logger *slog.Logger
	redis  *redis.Client
}

func NewRedisStorage(config config.RedisStorageConfig, logger *slog.Logger) (*Storage, error) {
	options, err := redis.ParseURL(config.Connection)
	if err != nil {
		return nil, err
	}

	var tlsConfig *tls.Config
	if config.SSLEnabled {
		clientCertPath := config.CertPath + "/redis-client.crt"
		clientKeyPath := config.CertPath + "/redis-client.key"
		caCertPath := config.CertPath + "/ca-redis.crt"

		caCert, err := os.ReadFile(caCertPath)
		if err != nil {
			log.Fatalf("Не удалось загрузить CA сертификат: %v", err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		clientCert, err := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
		if err != nil {
			log.Fatalf("Не удалось загрузить клиентский сертификат: %v", err)
		}

		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{clientCert},
			RootCAs:      caCertPool,
		}
	}

	rdb := redis.NewClient(options)
	options.TLSConfig = tlsConfig

	return &Storage{
		redis:  rdb,
		logger: logger,
	}, nil
}

func (f Storage) Put(ctx context.Context, path string, data []byte) error {
	res := f.redis.Set(ctx, path, data, 0)
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (f Storage) Get(ctx context.Context, path string) ([]byte, error) {
	res := f.redis.Get(ctx, path)
	if res.Err() != nil {
		return nil, res.Err()
	}

	return []byte(res.Val()), nil
}

func (f Storage) Delete(ctx context.Context, path string) error {
	res := f.redis.Del(ctx, path)
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}

func (f Storage) List(ctx context.Context, path string) ([]string, error) {
	var cursor uint64

	var dirs []string
	for {
		var keys []string
		var err error
		keys, cursor, err = f.redis.Scan(ctx, cursor, path+"*", 10).Result()
		if err != nil {
			return nil, err
		}

		for _, key := range keys {
			key = strings.TrimPrefix(key, path+string(os.PathSeparator))

			if !strings.Contains(key, string(os.PathSeparator)) {
				continue
			}
			key, _, _ = strings.Cut(key, string(os.PathSeparator))
			dirs = append(dirs, key)
		}

		if cursor == 0 {
			break
		}
	}

	slices.Sort(dirs)
	dirs = slices.Compact(dirs)

	return dirs, nil
}

func (f Storage) Close() error {
	return nil
}
