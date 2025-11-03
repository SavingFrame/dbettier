package database

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Store the first container instance so we can terminate it in TestMain
var sharedContainer *postgres.PostgresContainer

// TerminateSharedContainer terminates the shared container.
// This should be called from TestMain after all tests complete.
func TerminateSharedContainer() {
	if sharedContainer != nil {
		ctx := context.Background()
		if err := sharedContainer.Terminate(ctx); err != nil {
			// Just log, don't fail since tests are already done
			_ = err
		}
	}
}

// SetupTestDatabase creates a test database connection using either:
// 1. Existing PostgreSQL instance (if TEST_POSTGRES_HOST env var is set)
// 2. Shared testcontainer (reuses a single PostgreSQL container for all tests)
//
// Returns a connected database and a cleanup function that should be deferred.
// Note: The cleanup function only disconnects from the database; it does NOT
// terminate the container (the container is shared across all tests).
func SetupTestDatabase(t *testing.T) (*Database, func()) {
	t.Helper()

	// Check for existing database via environment variables
	if host := os.Getenv("TEST_POSTGRES_HOST"); host != "" {
		t.Log("Using existing PostgreSQL instance from environment variables")

		port := 5432
		if portStr := os.Getenv("TEST_POSTGRES_PORT"); portStr != "" {
			if p, err := strconv.Atoi(portStr); err == nil {
				port = p
			}
		}

		username := os.Getenv("TEST_POSTGRES_USER")
		if username == "" {
			username = "postgres"
		}

		password := os.Getenv("TEST_POSTGRES_PASSWORD")
		if password == "" {
			password = "postgres"
		}

		database := os.Getenv("TEST_POSTGRES_DB")
		if database == "" {
			database = "postgres"
		}

		db := NewDatabase(host, username, password, port, database)

		if err := db.Connect(); err != nil {
			t.Fatalf("Failed to connect to existing database: %v", err)
		}

		cleanup := func() {
			db.Disconnect()
		}

		return db, cleanup
	}

	t.Log("Using shared PostgreSQL testcontainer")
	ctx := context.Background()

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2),
		),
		testcontainers.WithReuseByName("dbettier-test-postgres"),
	)
	if err != nil {
		t.Fatalf("Failed to start/reuse postgres container: %v", err)
	}

	sharedContainer = postgresContainer

	host, err := postgresContainer.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get container host: %v", err)
	}

	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("Failed to get container port: %v", err)
	}

	db := NewDatabase(host, "testuser", "testpass", port.Int(), "testdb")

	if err := db.Connect(); err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	cleanup := func() {
		db.Disconnect()
	}

	return db, cleanup
}

// CreateSchema creates a new schema in the database.
// If the schema already exists, it will be dropped and recreated.
func CreateSchema(t *testing.T, db *Database, schemaName string) {
	t.Helper()
	ctx := context.Background()

	_, err := db.Connection.Exec(ctx, "DROP SCHEMA IF EXISTS "+schemaName+" CASCADE;")
	if err != nil {
		t.Fatalf("Failed to drop existing schema %s: %v", schemaName, err)
	}

	_, err = db.Connection.Exec(ctx, "CREATE SCHEMA "+schemaName+";")
	if err != nil {
		t.Fatalf("Failed to create schema %s: %v", schemaName, err)
	}
}

// DropSchemas drops one or more schemas from the database.
// Logs warnings if schemas cannot be dropped but doesn't fail the test.
func DropSchemas(t *testing.T, db *Database, schemaNames ...string) {
	t.Helper()
	ctx := context.Background()

	for _, schemaName := range schemaNames {
		_, err := db.Connection.Exec(ctx, "DROP SCHEMA IF EXISTS "+schemaName+" CASCADE;")
		if err != nil {
			t.Logf("Warning: Failed to drop schema %s: %v", schemaName, err)
		}
	}
}

// ExecQueries executes multiple SQL queries in sequence.
// Fails the test if any query fails.
func ExecQueries(t *testing.T, db *Database, queries ...string) {
	t.Helper()
	ctx := context.Background()

	for i, query := range queries {
		_, err := db.Connection.Exec(ctx, query)
		if err != nil {
			t.Fatalf("Failed to execute query #%d: %v\nQuery: %s", i+1, err, query)
		}
	}
}
