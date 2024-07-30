package main

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"log/slog"
	"messagio_assignment/internal/adapters/kafkaprod"
	"messagio_assignment/internal/adapters/pgstore"
	"messagio_assignment/internal/config"
	"messagio_assignment/internal/graceful"
	"messagio_assignment/internal/logger"
	"messagio_assignment/internal/ports/kafkacons"
	"messagio_assignment/internal/ports/rest"
	"messagio_assignment/internal/usecases"
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

	// Создание логгера
	slogger := logger.FromConfig(cfg.Environment, cfg.LogLevel)

	sarama.Logger = logger.NewSaramaLogger(slogger, slog.LevelInfo)
	sarama.DebugLogger = logger.NewSaramaLogger(slogger, slog.LevelDebug)

	// Перехват сигналов для graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Сетап graceful shutdown
	closer := graceful.NewCloser(slogger)
	defer func() {
		closer.Shutdown(cfg.ShutdownTimeout, 0)
	}()

	saramaCfg := sarama.NewConfig()
	saramaCfg.ClientID = cfg.Kafka.ClientID

	// Создание Kafka Producers
	kafkaProd, err := kafkaprod.New(slogger, saramaCfg, cfg.Kafka)
	if err != nil {
		slogger.Error("kafkaProd.New", logger.Err(err))
		return
	}
	closer.Add(func(ctx context.Context) error {
		err := kafkaProd.Close()
		if err != nil {
			return fmt.Errorf("kafka prod close %w", err)
		}
		slogger.Info("kafka prod is closed")
		return nil
	})

	// Создание Postgres Store
	store, err := pgstore.New(ctx, cfg.Postgres.ConnectionURL, slogger)
	if err != nil {
		slogger.Error("pgstore.New", logger.Err(err))
		return
	}
	closer.Add(func(ctx context.Context) error {
		store.Close()
		slogger.Info("pgstore is closed")
		return nil
	})

	// Применение миграций
	if cfg.Postgres.Migrate {
		err := store.Migrate(ctx)
		if err != nil {
			slogger.Error("migrate", logger.Err(err))
			return
		}
	}

	// Создание usecase
	messageUC := usecases.NewMessageUC(store.Message(), kafkaProd.Messages())

	// Создание и запуск Kafka Consumers
	kafkaCons, err := kafkacons.New(slogger, messageUC, saramaCfg, cfg.Kafka)
	if err != nil {
		slogger.Error("kafkacons.New", logger.Err(err))
		return
	}
	closer.Add(func(ctx context.Context) error {
		err := kafkaCons.Close()
		if err != nil {
			return fmt.Errorf("kafka cons close: %w", err)
		}
		slogger.Info("kafka cons is closed")
		return nil
	})
	go func() {
		kafkaCons.ProcessedMsgs().StartConsume(ctx)
	}()

	// Создание и запуск rest http сервера
	server := rest.NewServer(cfg, messageUC, slogger)
	closer.Add(func(ctx context.Context) error {
		if err := server.Shutdown(ctx); err != nil {
			return fmt.Errorf("rest http server shutdown: %w", err)
		}
		slogger.Info("rest http server is closed")
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
