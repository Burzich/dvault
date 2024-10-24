package postgres

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
	"github.com/jackc/pgx"
)

type Storage struct {
	conn   *pgx.ConnPool
	logger *slog.Logger
}

func NewPostgresStorage(config config.PostgresqlStorageConfig, logger *slog.Logger) (*Storage, error) {
	pgxCfg, err := pgx.ParseConnectionString(config.Connection)
	if err != nil {
		logger.Error("can't parse postgres connection string", slog.String("error", err.Error()))
		return nil, err
	}

	var tlsConfig *tls.Config
	if config.SSLEnabled {
		clientCertPath := config.CertPath + "/postgresql-client.crt"
		clientKeyPath := config.CertPath + "/postgresql-client.key"
		caCertPath := config.CertPath + "/ca-postgresql.crt"

		caCert, err := os.ReadFile(caCertPath)
		if err != nil {
			log.Fatalf("Can't load ca cert: %v", err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		clientCert, err := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
		if err != nil {
			log.Fatalf("Can't load client cert: %v", err)
		}

		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{clientCert},
			RootCAs:      caCertPool,
		}
	}

	pgxCfg.TLSConfig = tlsConfig
	connPool, err := pgx.NewConnPool(pgx.ConnPoolConfig{ConnConfig: pgxCfg})
	if err != nil {
		logger.Error("can't create connection pool", slog.String("error", err.Error()))
		return nil, err
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS data_table (
        id SERIAL PRIMARY KEY,
        key TEXT UNIQUE NOT NULL,
        data BYTEA
    );`

	_, err = connPool.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	return &Storage{
		conn:   connPool,
		logger: logger,
	}, nil
}

func (f Storage) Put(ctx context.Context, path string, data []byte) error {
	query := `INSERT INTO data_table (key, data) VALUES ($1, $2) ON CONFLICT (key) DO UPDATE SET data = $2;`
	_, err := f.conn.ExecEx(ctx, query, nil, path, data)
	if err != nil {
		return err
	}

	return nil
}

func (f Storage) Get(ctx context.Context, path string) ([]byte, error) {
	var byteData []byte
	query := `SELECT data FROM data_table WHERE key = $1`

	err := f.conn.QueryRowEx(ctx, query, nil, path).Scan(&byteData)
	if err != nil {
		return nil, err
	}

	return byteData, nil
}

func (f Storage) Delete(ctx context.Context, path string) error {
	query := `DELETE FROM data_table WHERE key = $1`
	_, err := f.conn.ExecEx(ctx, query, nil, path)
	if err != nil {
		return err
	}

	return nil
}

func (f Storage) List(ctx context.Context, path string) ([]string, error) {
	var keys []string
	query := `SELECT key FROM data_table WHERE key LIKE $1`
	rows, err := f.conn.QueryEx(ctx, query, nil, path+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var key string
		err := rows.Scan(&key)
		if err != nil {
			return nil, err
		}

		key = strings.TrimPrefix(key, path+string(os.PathSeparator))
		if !strings.Contains(key, string(os.PathSeparator)) {
			continue
		}
		key, _, _ = strings.Cut(key, string(os.PathSeparator))

		keys = append(keys, key)
	}

	slices.Sort(keys)
	keys = slices.Compact(keys)

	return keys, nil
}

func (f Storage) Close() error {
	f.conn.Close()
	return nil
}
