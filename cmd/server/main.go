package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"messagio_assignment/internal/adapters/pgstore"
	"messagio_assignment/internal/config"
	"messagio_assignment/internal/graceful"
	"messagio_assignment/internal/logger"
	"messagio_assignment/internal/ports/rest"
	"os"
	"os/signal"
	"syscall"
)

const DefaultConfigPath = "configs/development.yaml"

func main() {
	configPath := os.Getenv("APP_CONFIG_PATH")
	if configPath == "" {
		configPath = DefaultConfigPath
	}

	cfg, err := config.ReadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	slogger := logger.FromConfig(cfg.Environment, cfg.LogLevel)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	closer := graceful.NewCloser(slogger)
	defer func() {
		closer.Shutdown(cfg.ShutdownTimeout, 0)
	}()

	store, err := pgstore.New(ctx, cfg.Postgres.ConnectionURL, slogger)
	if err != nil {
		log.Println(err)
		return
	}
	closer.Add(func(ctx context.Context) error {
		store.Close()
		return nil
	})

	if cfg.Postgres.Migrate {
		err := store.Migrate(ctx)
		if err != nil {
			log.Printf("migrate: %v\n", err)
			return
		}
	}

	server := rest.NewServer(cfg, store.Message(), slogger)
	closer.Add(func(ctx context.Context) error {
		if err := server.Shutdown(ctx); err != nil {
			return fmt.Errorf("server shutdown: %w", err)
		}
		return nil
	})

	go func() {
		slogger.Info("listening...", slog.String("addr", server.Addr))

		err = server.ListenAndServe()
		if err != nil {
			slogger.Error("listen and serve", logger.Err(err))
		}
	}()

	<-ctx.Done()
}
