package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Burzich/dvault/internal/config"
	"github.com/Burzich/dvault/internal/dvault"
	"github.com/Burzich/dvault/internal/dvault/handler"
	"github.com/Burzich/dvault/internal/server"
)

func main() {
	cfg, err := config.Default()
	if err != nil {
		log.Fatal(err)
	}

	/*	pgxCfg, err := pgx.ParseConnectionString(cfg.Postgres.Addr)
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
		}*/

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	vault, err := dvault.NewDVault(logger, cfg.MountPath)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	vaultHandler := handler.NewHandler(vault)

	srv := server.NewServer(cfg.Server.Addr, vaultHandler)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	ready := make(chan struct{})
	go func() {
		<-ctx.Done()
		stop()
		logger.Info("shutdown starting")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("shutdown server", slog.String("error", err.Error()))
		}
		ready <- struct{}{}
	}()

	logger.Info("starting server")
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logger.Error("http server listen and serve", slog.String("error", err.Error()))
	}
	<-ready
	logger.Info("server shutdown")
}
