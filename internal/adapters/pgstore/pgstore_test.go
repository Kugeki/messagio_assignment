package pgstore

import (
	"context"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestPGStore(t *testing.T) {
	suite.Run(t, new(PGStoreTestSuite))
}

type postgresContainer struct {
	*postgres.PostgresContainer
	ConnectionString string
}

func createPostgresContainer(ctx context.Context) (*postgresContainer, error) {
	pgContainer, err := postgres.Run(ctx,
		"postgres:16.3-alpine",
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.WithSQLDriver("pgx"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(20*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	log.Println(connStr)

	return &postgresContainer{
		PostgresContainer: pgContainer,
		ConnectionString:  connStr,
	}, nil
}

type PGStoreTestSuite struct {
	suite.Suite
	pgContainer *postgresContainer
	ctx         context.Context

	store *Store
}

var EmptyDBSnapshot = "empty-snapshot"

func (su *PGStoreTestSuite) SetupSuite() {
	su.ctx = context.Background()
	pgContainer, err := createPostgresContainer(su.ctx)
	su.Require().NoError(err)

	su.pgContainer = pgContainer

	su.store, err = New(context.Background(), su.pgContainer.ConnectionString, nil)
	su.Require().NoError(err)
	su.Require().NotNil(su.store)

	err = su.store.Migrate(context.Background())
	su.Require().NoError(err)

	su.store.Close()

	err = su.pgContainer.Snapshot(su.ctx, postgres.WithSnapshotName(EmptyDBSnapshot))
	su.Require().NoError(err)
}

func (su *PGStoreTestSuite) TearDownSuite() {
	su.store.Close()

	if err := su.pgContainer.Terminate(su.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
}

func (su *PGStoreTestSuite) SetupTest() {
	su.RestoreDB()
	su.ReconnectStore()
}

func (su *PGStoreTestSuite) TearDownTest() {
	su.store.Close()
}

func (su *PGStoreTestSuite) SetupSubTest() {
	su.RestoreDB()
	su.ReconnectStore()
}

func (su *PGStoreTestSuite) TearDownSubTest() {
	su.store.Close()
}

func (su *PGStoreTestSuite) ReconnectStore() {
	var err error
	su.store, err = New(context.Background(), su.pgContainer.ConnectionString, nil)
	su.Require().NoError(err)
	su.Require().NotNil(su.store)
}

func (su *PGStoreTestSuite) RestoreDB() {
	err := su.pgContainer.Restore(su.ctx)
	su.Require().NoError(err)
}
