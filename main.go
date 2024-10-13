package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/Burzich/dvault/internal/config"
	"github.com/Burzich/dvault/internal/dvault"
	"github.com/Burzich/dvault/internal/dvault/handler"
	"github.com/Burzich/dvault/internal/server"
	"github.com/jackc/pgx"
)

func main() {
	cfg, err := config.Default()
	if err != nil {
		log.Fatal(err)
	}

	pgxCfg, err := pgx.ParseConnectionString(cfg.Postgres.Addr)
	if err != nil {
		log.Fatal(err)
		return
	}

	connPool, err := pgx.NewConnPool(pgx.ConnPoolConfig{ConnConfig: pgxCfg})
	if err != nil {
		log.Fatal(err)
		return
	}
	defer connPool.Close()

	_, err = connPool.Exec("SELECT 1")
	if err != nil {
		log.Fatal(err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	vault := dvault.NewDVault(logger)
	vaultHandler := handler.NewHandler(vault)

	srv := server.NewServer(cfg.Server.Addr, vaultHandler)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
