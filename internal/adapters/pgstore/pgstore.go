package pgstore

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
	"log/slog"
	"messagioassignment/internal/logger"
	ms "messagioassignment/migrations"

	_ "github.com/jackc/pgx/v5/stdlib" // for store.Migrate
)

type Store struct {
	db  *pgxpool.Pool
	log *slog.Logger

	messageRepo *MessageRepoPG
}

// New create new Store and connects to a database. Need call Close after this before goroutine shutdown.
func New(ctx context.Context, dbURL string, log *slog.Logger) (*Store, error) {
	s := &Store{log: log}

	if s.log == nil {
		s.log = logger.NewEraseLogger()
	}
	s.log = s.log.With(slog.String("component", "adapters/pgstore"))

	err := s.open(ctx, dbURL)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Store) open(ctx context.Context, dbURL string) error {
	conf, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return fmt.Errorf("postgres parse config: %w", err)
	}

	pgxLogger := logger.NewPgxLogger(s.log)
	conf.ConnConfig.Tracer = pgxLogger.TraceLog(ctx)

	pool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return fmt.Errorf("postgres pool create: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		return fmt.Errorf("postgres ping: %w", err)
	}

	s.db = pool

	return nil
}

func (s *Store) Migrate(ctx context.Context) error {
	db, err := sql.Open("pgx", s.db.Config().ConnString())
	if err != nil {
		return err
	}
	defer db.Close()

	migrateProvider, err := goose.NewProvider("postgres",
		db, ms.MigrateFS,
	)
	if err != nil {
		return err
	}

	migrateRes, err := migrateProvider.Up(ctx)
	if err != nil {
		return err
	}

	s.log.Info("applied migrations", slog.Int("count", len(migrateRes)))

	return nil
}

func (s *Store) Close() {
	if s.db == nil {
		s.log.Error("trying to close nil db")
		return
	}
	s.db.Close()

	s.log.Info("db is closed")
}

func (s *Store) Message() *MessageRepoPG {
	if s.messageRepo == nil {
		s.messageRepo = NewMessageRepoPG(s.db)
	}

	return s.messageRepo
}
